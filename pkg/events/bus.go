package events

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/nabob-ai/nomad/pkg/types"
)

// EventHandler 事件处理器函数
type EventHandler func(event any)

// EventBusConfig 事件总线配置
type EventBusConfig struct {
	MaxTimelineSize int           // 最大事件数 (默认 10000)
	MaxTimelineAge  time.Duration // 最大事件年龄 (默认 1小时)
	CleanupInterval time.Duration // 清理间隔 (默认 5分钟)
	EnableArchive   bool          // 是否启用归档 (预留)
}

// DefaultEventBusConfig 默认配置
func DefaultEventBusConfig() *EventBusConfig {
	return &EventBusConfig{
		MaxTimelineSize: 10000,
		MaxTimelineAge:  1 * time.Hour,
		CleanupInterval: 5 * time.Minute,
		EnableArchive:   false,
	}
}

// EventBus 三通道事件总线
type EventBus struct {
	mu sync.RWMutex

	// 配置
	config *EventBusConfig

	// 事件序列
	cursor    int64
	timeline  []types.AgentEventEnvelope
	bookmarks map[int64]types.Bookmark

	// 订阅者管理
	progressSubs map[string]chan types.AgentEventEnvelope
	controlSubs  map[string]chan types.AgentEventEnvelope
	monitorSubs  map[string]chan types.AgentEventEnvelope

	// 回调处理器
	controlHandlers map[string][]EventHandler
	monitorHandlers map[string][]EventHandler

	// 清理 Worker
	cleanupTicker *time.Ticker
	cleanupDone   chan struct{}
	cleanupWg     sync.WaitGroup // 等待清理 goroutine 退出
}

// NewEventBus 创建新的事件总线（使用默认配置）
func NewEventBus() *EventBus {
	return NewEventBusWithConfig(DefaultEventBusConfig())
}

// NewEventBusWithConfig 创建带配置的事件总线
func NewEventBusWithConfig(config *EventBusConfig) *EventBus {
	if config == nil {
		config = DefaultEventBusConfig()
	}

	eb := &EventBus{
		config:          config,
		timeline:        make([]types.AgentEventEnvelope, 0, 1000),
		bookmarks:       make(map[int64]types.Bookmark),
		progressSubs:    make(map[string]chan types.AgentEventEnvelope),
		controlSubs:     make(map[string]chan types.AgentEventEnvelope),
		monitorSubs:     make(map[string]chan types.AgentEventEnvelope),
		controlHandlers: make(map[string][]EventHandler),
		monitorHandlers: make(map[string][]EventHandler),
		cleanupDone:     make(chan struct{}),
	}

	// 启动清理 worker
	eb.startCleanupWorker()

	return eb
}

// startCleanupWorker 启动后台清理协程
func (eb *EventBus) startCleanupWorker() {
	if eb.config == nil || eb.config.CleanupInterval <= 0 {
		return
	}

	eb.cleanupTicker = time.NewTicker(eb.config.CleanupInterval)
	eb.cleanupWg.Add(1)

	go func() {
		defer eb.cleanupWg.Done()
		for {
			select {
			case <-eb.cleanupTicker.C:
				eb.cleanup()
			case <-eb.cleanupDone:
				eb.cleanupTicker.Stop()
				return
			}
		}
	}()
}

// cleanup 清理过期和超量的事件
func (eb *EventBus) cleanup() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(eb.timeline) == 0 {
		return
	}

	now := time.Now().Unix()
	cutoffIdx := 0

	// 按时间清理：移除超过 MaxTimelineAge 的事件
	if eb.config.MaxTimelineAge > 0 {
		maxAge := int64(eb.config.MaxTimelineAge.Seconds())
		for i, env := range eb.timeline {
			if now-env.Bookmark.Timestamp <= maxAge {
				cutoffIdx = i
				break
			}
			// 删除对应的 bookmark
			delete(eb.bookmarks, env.Cursor)
		}
	}

	// 按数量清理：保留最新的 MaxTimelineSize 个事件
	if eb.config.MaxTimelineSize > 0 && len(eb.timeline)-cutoffIdx > eb.config.MaxTimelineSize {
		newCutoff := len(eb.timeline) - eb.config.MaxTimelineSize
		if newCutoff > cutoffIdx {
			// 删除多余事件的 bookmarks
			for i := cutoffIdx; i < newCutoff; i++ {
				delete(eb.bookmarks, eb.timeline[i].Cursor)
			}
			cutoffIdx = newCutoff
		}
	}

	// 执行切片截断
	if cutoffIdx > 0 {
		eb.timeline = eb.timeline[cutoffIdx:]
	}
}

// Close 关闭事件总线，释放资源
func (eb *EventBus) Close() {
	// 先停止清理 worker（不持锁，避免死锁）
	if eb.cleanupDone != nil {
		close(eb.cleanupDone)
		eb.cleanupWg.Wait() // 等待清理 goroutine 完全退出
		eb.cleanupDone = nil
	}

	// 然后清理所有资源（持锁）
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 关闭所有订阅 channel
	for id, ch := range eb.progressSubs {
		close(ch)
		delete(eb.progressSubs, id)
	}
	for id, ch := range eb.controlSubs {
		close(ch)
		delete(eb.controlSubs, id)
	}
	for id, ch := range eb.monitorSubs {
		close(ch)
		delete(eb.monitorSubs, id)
	}

	// 清空数据
	eb.timeline = nil
	eb.bookmarks = nil
}

// emit 发送事件到总线(内部方法)
func (eb *EventBus) emit(channel types.AgentChannel, event any) types.AgentEventEnvelope {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 增加cursor
	eb.cursor++

	// 创建Bookmark
	bookmark := types.Bookmark{
		Cursor:    eb.cursor,
		Timestamp: time.Now().Unix(),
	}

	// 封装事件
	envelope := types.AgentEventEnvelope{
		Cursor:   eb.cursor,
		Bookmark: bookmark,
		Event:    event,
	}

	// 保存到时间线
	eb.timeline = append(eb.timeline, envelope)
	eb.bookmarks[eb.cursor] = bookmark

	// 检查是否是重要事件（done事件必须送达）
	_, isDoneEvent := event.(*types.ProgressDoneEvent)

	// 分发到对应通道的订阅者
	switch channel {
	case types.ChannelProgress:
		for _, ch := range eb.progressSubs {
			if isDoneEvent {
				// done 事件使用带超时的发送，确保送达
				select {
				case ch <- envelope:
				case <-time.After(5 * time.Second):
					// 超时，记录日志但继续
				}
			} else {
				select {
				case ch <- envelope:
				default:
					// 非阻塞发送,如果channel满了则跳过
				}
			}
		}
	case types.ChannelControl:
		for _, ch := range eb.controlSubs {
			select {
			case ch <- envelope:
			default:
			}
		}
		// 调用Control回调处理器
		eb.invokeHandlers(eb.controlHandlers, event)
	case types.ChannelMonitor:
		for _, ch := range eb.monitorSubs {
			select {
			case ch <- envelope:
			default:
			}
		}
		// 调用Monitor回调处理器
		eb.invokeHandlers(eb.monitorHandlers, event)
	}

	return envelope
}

// invokeHandlers 调用事件处理器
func (eb *EventBus) invokeHandlers(handlers map[string][]EventHandler, event any) {
	// 获取事件类型
	eventType := ""
	if e, ok := event.(types.EventType); ok {
		eventType = e.EventType()
	}

	// 调用特定类型的处理器
	if hs, ok := handlers[eventType]; ok {
		for _, h := range hs {
			go h(event) // 异步调用
		}
	}

	// 调用通配符处理器
	if hs, ok := handlers["*"]; ok {
		for _, h := range hs {
			go h(event)
		}
	}
}

// EmitProgress 发送Progress事件
func (eb *EventBus) EmitProgress(event any) types.AgentEventEnvelope {
	return eb.emit(types.ChannelProgress, event)
}

// EmitControl 发送Control事件
func (eb *EventBus) EmitControl(event any) types.AgentEventEnvelope {
	return eb.emit(types.ChannelControl, event)
}

// EmitMonitor 发送Monitor事件
func (eb *EventBus) EmitMonitor(event any) types.AgentEventEnvelope {
	return eb.emit(types.ChannelMonitor, event)
}

// Subscribe 订阅指定通道的事件(返回channel)
func (eb *EventBus) Subscribe(channels []types.AgentChannel, opts *types.SubscribeOptions) <-chan types.AgentEventEnvelope {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 创建缓冲channel(避免阻塞)
	ch := make(chan types.AgentEventEnvelope, 100)

	// 生成唯一订阅ID
	subID := generateSubID()

	// 注册到对应通道
	if len(channels) == 0 {
		channels = []types.AgentChannel{types.ChannelProgress, types.ChannelControl, types.ChannelMonitor}
	}

	for _, channel := range channels {
		switch channel {
		case types.ChannelProgress:
			eb.progressSubs[subID] = ch
		case types.ChannelControl:
			eb.controlSubs[subID] = ch
		case types.ChannelMonitor:
			eb.monitorSubs[subID] = ch
		}
	}

	// 如果指定了since,回放历史事件
	if opts != nil && opts.Since != nil {
		go eb.replay(ch, opts.Since, opts.Kinds, channels)
	}

	return ch
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(ch <-chan types.AgentEventEnvelope) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var writeCh chan types.AgentEventEnvelope
	found := false

	// 从所有订阅 map 中查找并移除（需要检查所有 map，因为同一个 channel 可能订阅了多个通道）
	// 在第一次找到时保存双向 channel 用于关闭
	for id, subCh := range eb.progressSubs {
		if subCh == ch {
			delete(eb.progressSubs, id)
			if !found {
				writeCh = subCh
				found = true
			}
		}
	}
	for id, subCh := range eb.controlSubs {
		if subCh == ch {
			delete(eb.controlSubs, id)
			if !found {
				writeCh = subCh
				found = true
			}
		}
	}
	for id, subCh := range eb.monitorSubs {
		if subCh == ch {
			delete(eb.monitorSubs, id)
			if !found {
				writeCh = subCh
				found = true
			}
		}
	}

	// 只关闭一次 channel
	if found && writeCh != nil {
		close(writeCh)
	}
}

// OnControl 注册Control事件处理器
func (eb *EventBus) OnControl(eventType string, handler EventHandler) func() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.controlHandlers[eventType] = append(eb.controlHandlers[eventType], handler)

	// 返回取消函数
	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		// 从处理器列表中移除
		handlers := eb.controlHandlers[eventType]
		for i, h := range handlers {
			// Go中函数比较困难,这里简化处理
			if &h == &handler {
				eb.controlHandlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// OnMonitor 注册Monitor事件处理器
func (eb *EventBus) OnMonitor(eventType string, handler EventHandler) func() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.monitorHandlers[eventType] = append(eb.monitorHandlers[eventType], handler)

	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		handlers := eb.monitorHandlers[eventType]
		for i, h := range handlers {
			if &h == &handler {
				eb.monitorHandlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// replay 回放历史事件
func (eb *EventBus) replay(ch chan types.AgentEventEnvelope, since *types.Bookmark, kinds []string, channels []types.AgentChannel) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	// 创建事件类型过滤器
	kindFilter := make(map[string]bool)
	if len(kinds) > 0 {
		for _, k := range kinds {
			kindFilter[k] = true
		}
	}

	// 创建通道过滤器
	channelFilter := make(map[types.AgentChannel]bool)
	for _, c := range channels {
		channelFilter[c] = true
	}

	// 遍历时间线,发送符合条件的事件
	for _, envelope := range eb.timeline {
		// 跳过since之前的事件
		if since != nil && envelope.Bookmark.Cursor <= since.Cursor {
			continue
		}

		// 检查通道过滤
		if e, ok := envelope.Event.(types.EventType); ok {
			if len(channelFilter) > 0 && !channelFilter[e.Channel()] {
				continue
			}

			// 检查类型过滤
			if len(kindFilter) > 0 && !kindFilter[e.EventType()] {
				continue
			}
		}

		// 发送事件
		select {
		case ch <- envelope:
		default:
			return // channel已关闭或满
		}
	}
}

// GetCursor 获取当前cursor
func (eb *EventBus) GetCursor() int64 {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return eb.cursor
}

// GetLastBookmark 获取最后一个bookmark
func (eb *EventBus) GetLastBookmark() *types.Bookmark {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.cursor == 0 {
		return nil
	}

	if bm, ok := eb.bookmarks[eb.cursor]; ok {
		return &bm
	}
	return nil
}

// GetTimeline 获取完整时间线（不推荐用于大量事件，请使用 GetTimelineRange）
func (eb *EventBus) GetTimeline() []types.AgentEventEnvelope {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	// 返回副本
	timeline := make([]types.AgentEventEnvelope, len(eb.timeline))
	copy(timeline, eb.timeline)
	return timeline
}

// GetTimelineRange 获取指定范围的时间线（基于索引）
func (eb *EventBus) GetTimelineRange(start, limit int) []types.AgentEventEnvelope {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if start < 0 {
		start = 0
	}
	if start >= len(eb.timeline) {
		return []types.AgentEventEnvelope{}
	}

	end := start + limit
	if limit <= 0 || end > len(eb.timeline) {
		end = len(eb.timeline)
	}

	result := make([]types.AgentEventEnvelope, end-start)
	copy(result, eb.timeline[start:end])
	return result
}

// GetTimelineSince 获取指定 cursor 之后的所有事件
func (eb *EventBus) GetTimelineSince(cursor int64) []types.AgentEventEnvelope {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	result := make([]types.AgentEventEnvelope, 0)
	for _, envelope := range eb.timeline {
		if envelope.Cursor > cursor {
			result = append(result, envelope)
		}
	}
	return result
}

// GetTimelineFiltered 获取过滤后的时间线
func (eb *EventBus) GetTimelineFiltered(filter func(types.AgentEventEnvelope) bool) []types.AgentEventEnvelope {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	result := make([]types.AgentEventEnvelope, 0)
	for _, envelope := range eb.timeline {
		if filter(envelope) {
			result = append(result, envelope)
		}
	}
	return result
}

// GetTimelineCount 获取当前时间线事件数量
func (eb *EventBus) GetTimelineCount() int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.timeline)
}

// Clear 清空事件总线(用于测试)
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.cursor = 0
	eb.timeline = make([]types.AgentEventEnvelope, 0, 1000)
	eb.bookmarks = make(map[int64]types.Bookmark)
}

// generateSubID 生成订阅ID
func generateSubID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串（使用 crypto/rand）
func randomString(n int) string {
	b := make([]byte, (n+1)/2)
	if _, err := rand.Read(b); err != nil {
		// fallback: 使用时间戳
		return time.Now().Format("150405.000000000")[:n]
	}
	return hex.EncodeToString(b)[:n]
}
