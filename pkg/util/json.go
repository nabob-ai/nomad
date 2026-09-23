package util

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

// DeterministicJSON 提供确定性的 JSON 序列化
// 主要用于 KV-Cache 优化：确保相同的数据结构总是生成相同的 JSON 输出
//
// Go 的 map 迭代顺序是随机的，这会导致：
// - 相同内容的 JSON 输出不同
// - LLM Provider 的 KV-Cache 失效
// - 成本增加 10x（根据 Manus 团队数据）
//
// 此包通过递归排序 map 的 key 来解决这个问题

// MarshalDeterministic 确定性地序列化为 JSON
// 对于 map 类型，会按 key 字典序排序
func MarshalDeterministic(v any) ([]byte, error) {
	normalized := normalizeValue(reflect.ValueOf(v))
	return json.Marshal(normalized)
}

// MarshalDeterministicIndent 确定性地序列化为格式化的 JSON
func MarshalDeterministicIndent(v any, prefix, indent string) ([]byte, error) {
	normalized := normalizeValue(reflect.ValueOf(v))
	return json.MarshalIndent(normalized, prefix, indent)
}

// MarshalDeterministicToBuffer 确定性地序列化到 buffer（减少内存分配）
func MarshalDeterministicToBuffer(v any, buf *bytes.Buffer) error {
	normalized := normalizeValue(reflect.ValueOf(v))
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(normalized)
}

// normalizeValue 递归处理值，确保 map 的 key 有序
func normalizeValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	// 处理指针和接口
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		return normalizeMap(v)
	case reflect.Slice:
		return normalizeSlice(v)
	case reflect.Array:
		return normalizeArray(v)
	case reflect.Struct:
		return normalizeStruct(v)
	default:
		// 基本类型直接返回
		if v.CanInterface() {
			return v.Interface()
		}
		return nil
	}
}

// normalizeMap 按 key 排序处理 map
func normalizeMap(v reflect.Value) any {
	if v.IsNil() {
		return nil
	}

	// 获取所有 key 并排序
	keys := v.MapKeys()
	sortedKeys := make([]string, 0, len(keys))
	keyMap := make(map[string]reflect.Value, len(keys))

	for _, k := range keys {
		// 将 key 转换为字符串用于排序
		var keyStr string
		switch k.Kind() {
		case reflect.String:
			keyStr = k.String()
		default:
			// 对于非字符串 key，序列化为 JSON 字符串
			keyBytes, _ := json.Marshal(k.Interface())
			keyStr = string(keyBytes)
		}
		sortedKeys = append(sortedKeys, keyStr)
		keyMap[keyStr] = k
	}

	sort.Strings(sortedKeys)

	// 构建有序的 map（使用 orderedMap 保持顺序）
	result := make(orderedMap, 0, len(sortedKeys))
	for _, keyStr := range sortedKeys {
		k := keyMap[keyStr]
		val := v.MapIndex(k)
		result = append(result, orderedMapEntry{
			Key:   keyStr,
			Value: normalizeValue(val),
		})
	}

	return result
}

// orderedMapEntry 有序 map 的条目
type orderedMapEntry struct {
	Key   string
	Value any
}

// orderedMap 有序 map，实现 json.Marshaler 接口
type orderedMap []orderedMapEntry

// MarshalJSON 按顺序输出 JSON
func (m orderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i, entry := range m {
		if i > 0 {
			buf.WriteByte(',')
		}

		// 写入 key
		keyBytes, err := json.Marshal(entry.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		// 写入 value
		valBytes, err := json.Marshal(entry.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// normalizeSlice 处理 slice
func normalizeSlice(v reflect.Value) any {
	if v.IsNil() {
		return nil
	}

	result := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = normalizeValue(v.Index(i))
	}
	return result
}

// normalizeArray 处理数组
func normalizeArray(v reflect.Value) any {
	result := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = normalizeValue(v.Index(i))
	}
	return result
}

// normalizeStruct 处理结构体
// 注意：Go 的 json.Marshal 对结构体已经是确定性的（按字段定义顺序）
// 但我们仍需递归处理嵌套的 map
func normalizeStruct(v reflect.Value) any {
	t := v.Type()
	result := make(orderedMap, 0)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// 跳过未导出字段
		if field.PkgPath != "" {
			continue
		}

		// 获取 JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// 解析 tag
		name := field.Name
		omitempty := false
		if jsonTag != "" {
			parts := splitTag(jsonTag)
			if parts[0] != "" {
				name = parts[0]
			}
			for _, opt := range parts[1:] {
				if opt == "omitempty" {
					omitempty = true
				}
			}
		}

		// 处理 omitempty
		if omitempty && isEmptyValue(fieldVal) {
			continue
		}

		result = append(result, orderedMapEntry{
			Key:   name,
			Value: normalizeValue(fieldVal),
		})
	}

	return result
}

// splitTag 分割 JSON tag
func splitTag(tag string) []string {
	var result []string
	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] != ',' {
			i++
		}
		result = append(result, tag[:i])
		if i < len(tag) {
			tag = tag[i+1:]
		} else {
			break
		}
	}
	return result
}

// isEmptyValue 检查是否是空值
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
