// Package stream 提供泛型流处理能力
// 支持类型安全的流操作，包括复制、合并和转换
package stream

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
)

// ErrSkip 用于在 Transform 过程中跳过某个值
// 示例:
//
//	outStream := Transform(inStream, func(s string) (string, error) {
//	    if len(s) == 0 {
//	        return "", ErrSkip
//	    }
//	    return s, nil
//	})
var ErrSkip = errors.New("skip this value")

// ErrRecvAfterClosed 表示在流关闭后调用了 Recv
var ErrRecvAfterClosed = errors.New("recv after stream closed")

// SourceEOF 表示合并流中某个源流的 EOF
type SourceEOF struct {
	Name string
}

func (e *SourceEOF) Error() string {
	return "EOF from source: " + e.Name
}

// GetSourceName 从 SourceEOF 错误中提取源流名称
func GetSourceName(err error) (string, bool) {
	var sErr *SourceEOF
	if errors.As(err, &sErr) {
		return sErr.Name, true
	}
	return "", false
}

// item 表示流中的单个值
type item[T any] struct {
	value T
	err   error
}

// stream 是基于通道的核心实现
type stream[T any] struct {
	items  chan item[T]
	closed chan struct{}

	automaticClose bool
	closedFlag     *uint32
}

func newStream[T any](cap int) *stream[T] {
	return &stream[T]{
		items:  make(chan item[T], cap),
		closed: make(chan struct{}),
	}
}

func (s *stream[T]) send(value T, err error) bool {
	select {
	case <-s.closed:
		return true
	default:
	}

	select {
	case <-s.closed:
		return true
	case s.items <- item[T]{value, err}:
		return false
	}
}

func (s *stream[T]) recv() (T, error) {
	it, ok := <-s.items
	if !ok {
		it.err = io.EOF
	}
	return it.value, it.err
}

func (s *stream[T]) closeSend() {
	close(s.items)
}

func (s *stream[T]) closeRecv() {
	if s.automaticClose {
		if atomic.CompareAndSwapUint32(s.closedFlag, 0, 1) {
			close(s.closed)
		}
		return
	}
	close(s.closed)
}

// Reader 是流的接收端
type Reader[T any] struct {
	typ readerType
	st  *stream[T]

	// 数组读取器
	ar *arrayReader[T]

	// 多流合并读取器
	msr *multiStreamReader[T]

	// 转换读取器
	cvt *convertReader[T]

	// 子读取器（用于 Copy）
	child *childReader[T]
}

type readerType int

const (
	readerTypeStream readerType = iota
	readerTypeArray
	readerTypeMulti
	readerTypeConvert
	readerTypeChild
)

// Recv 从流中接收下一个值
// 当流耗尽时返回 io.EOF
func (r *Reader[T]) Recv() (T, error) {
	switch r.typ {
	case readerTypeStream:
		return r.st.recv()
	case readerTypeArray:
		return r.ar.recv()
	case readerTypeMulti:
		return r.msr.recv()
	case readerTypeConvert:
		return r.cvt.recv()
	case readerTypeChild:
		return r.child.recv()
	default:
		var zero T
		return zero, errors.New("invalid reader type")
	}
}

// Close 关闭读取器并释放资源
func (r *Reader[T]) Close() {
	switch r.typ {
	case readerTypeStream:
		r.st.closeRecv()
	case readerTypeArray:
		// 无需清理
	case readerTypeMulti:
		r.msr.close()
	case readerTypeConvert:
		r.cvt.close()
	case readerTypeChild:
		r.child.close()
	}
}

// SetAutomaticClose 启用 GC 时自动清理
func (r *Reader[T]) SetAutomaticClose() {
	switch r.typ {
	case readerTypeStream:
		if !r.st.automaticClose {
			r.st.automaticClose = true
			var flag uint32
			r.st.closedFlag = &flag
			runtime.SetFinalizer(r, func(reader *Reader[T]) {
				reader.Close()
			})
		}
	case readerTypeMulti:
		for _, s := range r.msr.streams {
			if !s.automaticClose {
				s.automaticClose = true
				var flag uint32
				s.closedFlag = &flag
				runtime.SetFinalizer(s, func(st *stream[T]) {
					st.closeRecv()
				})
			}
		}
	case readerTypeChild:
		r.child.parent.reader.SetAutomaticClose()
	case readerTypeConvert:
		r.cvt.source.SetAutomaticClose()
	}
}

// Copy 创建 n 个独立的读取器
// 调用 Copy 后原读取器不可用
// 每个副本可以独立读取，互不影响
func (r *Reader[T]) Copy(n int) []*Reader[T] {
	if n < 2 {
		return []*Reader[T]{r}
	}

	if r.typ == readerTypeArray {
		copies := make([]*Reader[T], n)
		for i, ar := range r.ar.copy(n) {
			copies[i] = &Reader[T]{typ: readerTypeArray, ar: ar}
		}
		return copies
	}

	return copyReaders(r, n)
}

// Collect 将流中所有值读取到切片中
func (r *Reader[T]) Collect() ([]T, error) {
	var result []T
	for {
		v, err := r.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			return result, err
		}
		result = append(result, v)
	}
}

// Writer 是流的发送端
type Writer[T any] struct {
	st *stream[T]
}

// Send 向流发送一个值
// 如果流已被接收端关闭，返回 true
func (w *Writer[T]) Send(value T, err error) bool {
	return w.st.send(value, err)
}

// Close 关闭写入器，向接收端发送 EOF 信号
func (w *Writer[T]) Close() {
	w.st.closeSend()
}

// Pipe 创建具有指定缓冲容量的新流
// 返回用于接收的 Reader 和用于发送的 Writer
//
// 示例:
//
//	reader, writer := Pipe[string](10)
//	go func() {
//	    defer writer.Close()
//	    writer.Send("hello", nil)
//	    writer.Send("world", nil)
//	}()
//
//	for {
//	    v, err := reader.Recv()
//	    if errors.Is(err, io.EOF) {
//	        break
//	    }
//	    fmt.Println(v)
//	}
func Pipe[T any](cap int) (*Reader[T], *Writer[T]) {
	s := newStream[T](cap)
	return &Reader[T]{typ: readerTypeStream, st: s}, &Writer[T]{st: s}
}

// FromSlice 从切片创建 Reader
func FromSlice[T any](slice []T) *Reader[T] {
	return &Reader[T]{
		typ: readerTypeArray,
		ar:  &arrayReader[T]{arr: slice},
	}
}

// arrayReader 将切片包装为流
type arrayReader[T any] struct {
	arr   []T
	index int
}

func (ar *arrayReader[T]) recv() (T, error) {
	if ar.index < len(ar.arr) {
		v := ar.arr[ar.index]
		ar.index++
		return v, nil
	}
	var zero T
	return zero, io.EOF
}

func (ar *arrayReader[T]) copy(n int) []*arrayReader[T] {
	copies := make([]*arrayReader[T], n)
	for i := range copies {
		copies[i] = &arrayReader[T]{
			arr:   ar.arr,
			index: ar.index,
		}
	}
	return copies
}

func (ar *arrayReader[T]) toStream() *stream[T] {
	s := newStream[T](len(ar.arr) - ar.index)
	for i := ar.index; i < len(ar.arr); i++ {
		s.send(ar.arr[i], nil)
	}
	s.closeSend()
	return s
}

// multiStreamReader 合并多个流
type multiStreamReader[T any] struct {
	streams     []*stream[T]
	nonClosed   []int
	sourceNames []string
}

func newMultiStreamReader[T any](streams []*stream[T], names []string) *multiStreamReader[T] {
	nonClosed := make([]int, len(streams))
	for i := range streams {
		nonClosed[i] = i
	}
	return &multiStreamReader[T]{
		streams:     streams,
		nonClosed:   nonClosed,
		sourceNames: names,
	}
}

func (msr *multiStreamReader[T]) recv() (T, error) {
	var zero T

	for len(msr.nonClosed) > 0 {
		// 简单的轮询 + select
		for i, idx := range msr.nonClosed {
			select {
			case it, ok := <-msr.streams[idx].items:
				if !ok {
					// 从 nonClosed 中移除
					msr.nonClosed = append(msr.nonClosed[:i], msr.nonClosed[i+1:]...)
					if len(msr.sourceNames) > 0 {
						return zero, &SourceEOF{Name: msr.sourceNames[idx]}
					}
					continue
				}
				return it.value, it.err
			default:
				continue
			}
		}

		// 如果没有立即可用的数据，阻塞等待第一个
		if len(msr.nonClosed) > 0 {
			idx := msr.nonClosed[0]
			it, ok := <-msr.streams[idx].items
			if !ok {
				msr.nonClosed = msr.nonClosed[1:]
				if len(msr.sourceNames) > 0 {
					return zero, &SourceEOF{Name: msr.sourceNames[idx]}
				}
				continue
			}
			return it.value, it.err
		}
	}

	return zero, io.EOF
}

func (msr *multiStreamReader[T]) close() {
	for _, s := range msr.streams {
		s.closeRecv()
	}
}

// convertReader 对流值应用转换
type convertReader[T any] struct {
	source  *Reader[any]
	convert func(any) (T, error)
}

func (cr *convertReader[T]) recv() (T, error) {
	for {
		v, err := cr.source.Recv()
		if err != nil {
			var zero T
			return zero, err
		}

		result, err := cr.convert(v)
		if err == nil {
			return result, nil
		}

		if !errors.Is(err, ErrSkip) {
			return result, err
		}
		// ErrSkip: 继续下一个值
	}
}

func (cr *convertReader[T]) close() {
	cr.source.Close()
}

// 流复制实现 - 使用链表 + sync.Once 实现零拷贝

type copyElement[T any] struct {
	once sync.Once
	next *copyElement[T]
	item item[T]
}

type parentReader[T any] struct {
	reader    *Reader[T]
	children  []*copyElement[T]
	closedNum uint32
}

func (p *parentReader[T]) peek(idx int) (T, error) {
	elem := p.children[idx]
	if elem == nil {
		var zero T
		return zero, ErrRecvAfterClosed
	}

	elem.once.Do(func() {
		v, err := p.reader.Recv()
		elem.item = item[T]{value: v, err: err}
		if !errors.Is(err, io.EOF) {
			elem.next = &copyElement[T]{}
			p.children[idx] = elem.next
		}
	})

	if !errors.Is(elem.item.err, io.EOF) {
		p.children[idx] = elem.next
	}

	return elem.item.value, elem.item.err
}

func (p *parentReader[T]) closeChild(idx int) {
	if p.children[idx] == nil {
		return
	}
	p.children[idx] = nil

	if int(atomic.AddUint32(&p.closedNum, 1)) == len(p.children) {
		p.reader.Close()
	}
}

type childReader[T any] struct {
	parent *parentReader[T]
	index  int
}

func (cr *childReader[T]) recv() (T, error) {
	return cr.parent.peek(cr.index)
}

func (cr *childReader[T]) close() {
	cr.parent.closeChild(cr.index)
}

func copyReaders[T any](r *Reader[T], n int) []*Reader[T] {
	parent := &parentReader[T]{
		reader:   r,
		children: make([]*copyElement[T], n),
	}

	// 使用共享的空元素初始化
	elem := &copyElement[T]{}
	for i := range parent.children {
		parent.children[i] = elem
	}

	copies := make([]*Reader[T], n)
	for i := range copies {
		copies[i] = &Reader[T]{
			typ: readerTypeChild,
			child: &childReader[T]{
				parent: parent,
				index:  i,
			},
		}
	}

	return copies
}

// Merge 将多个读取器合并为一个
// 从任何有数据的读取器接收值
func Merge[T any](readers ...*Reader[T]) *Reader[T] {
	if len(readers) == 0 {
		return nil
	}
	if len(readers) == 1 {
		return readers[0]
	}

	streams := make([]*stream[T], len(readers))
	for i, r := range readers {
		streams[i] = readerToStream(r)
	}

	return &Reader[T]{
		typ: readerTypeMulti,
		msr: newMultiStreamReader(streams, nil),
	}
}

// MergeNamed 将多个命名读取器合并为一个
// 当源流结束时，返回带有流名称的 SourceEOF
func MergeNamed[T any](readers map[string]*Reader[T]) *Reader[T] {
	if len(readers) == 0 {
		return nil
	}

	streams := make([]*stream[T], 0, len(readers))
	names := make([]string, 0, len(readers))

	for name, r := range readers {
		streams = append(streams, readerToStream(r))
		names = append(names, name)
	}

	return &Reader[T]{
		typ: readerTypeMulti,
		msr: newMultiStreamReader(streams, names),
	}
}

func readerToStream[T any](r *Reader[T]) *stream[T] {
	switch r.typ {
	case readerTypeStream:
		return r.st
	case readerTypeArray:
		return r.ar.toStream()
	default:
		// 通过读取和发送转换其他类型
		s := newStream[T](5)
		go func() {
			defer s.closeSend()
			for {
				v, err := r.Recv()
				if errors.Is(err, io.EOF) {
					return
				}
				if s.send(v, err) {
					return
				}
			}
		}()
		return s
	}
}

// Transform 将类型 T 的流转换为类型 U
// 返回 ErrSkip 可跳过某个值
func Transform[T, U any](r *Reader[T], fn func(T) (U, error)) *Reader[U] {
	return TransformSimple(r, fn)
}

// TransformSimple 创建新的 goroutine 进行转换
func TransformSimple[T, U any](r *Reader[T], fn func(T) (U, error)) *Reader[U] {
	out, writer := Pipe[U](5)

	go func() {
		defer writer.Close()
		for {
			v, err := r.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				writer.Send(*new(U), err)
				return
			}

			result, err := fn(v)
			if errors.Is(err, ErrSkip) {
				continue
			}
			if err != nil {
				writer.Send(*new(U), err)
				return
			}

			if writer.Send(result, nil) {
				return
			}
		}
	}()

	return out
}

// Filter 创建只包含匹配谓词的值的新读取器
func Filter[T any](r *Reader[T], predicate func(T) bool) *Reader[T] {
	return TransformSimple(r, func(v T) (T, error) {
		if predicate(v) {
			return v, nil
		}
		return v, ErrSkip
	})
}

// Map 对流中的每个值应用函数
func Map[T, U any](r *Reader[T], fn func(T) U) *Reader[U] {
	return TransformSimple(r, func(v T) (U, error) {
		return fn(v), nil
	})
}

// Take 返回只产出前 n 个值的读取器
func Take[T any](r *Reader[T], n int) *Reader[T] {
	out, writer := Pipe[T](n)
	count := 0

	go func() {
		defer writer.Close()
		for count < n {
			v, err := r.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				writer.Send(*new(T), err)
				return
			}
			if writer.Send(v, nil) {
				return
			}
			count++
		}
	}()

	return out
}

// ForEach 消费读取器中的所有值，对每个值调用 fn
func ForEach[T any](r *Reader[T], fn func(T) error) error {
	defer r.Close()
	for {
		v, err := r.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := fn(v); err != nil {
			return err
		}
	}
}
