package logic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	t.Run("record save operations", func(t *testing.T) {
		m := NewMetrics()

		// 记录成功的保存
		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 15*time.Millisecond, nil)
		m.RecordSave("user:456", "behavior", ScopeSession, 20*time.Millisecond, nil)

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(3), snapshot.MemorySaveTotal)
		assert.Equal(t, int64(0), snapshot.MemorySaveErrors)
		assert.Equal(t, int64(2), snapshot.MemoryByNamespace["user:123"])
		assert.Equal(t, int64(1), snapshot.MemoryByNamespace["user:456"])
	})

	t.Run("record save errors", func(t *testing.T) {
		m := NewMetrics()

		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, assert.AnError)

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(2), snapshot.MemorySaveTotal)
		assert.Equal(t, int64(1), snapshot.MemorySaveErrors)
	})

	t.Run("record get operations", func(t *testing.T) {
		m := NewMetrics()

		m.RecordGet(5*time.Millisecond, nil)
		m.RecordGet(10*time.Millisecond, nil)
		m.RecordGet(15*time.Millisecond, assert.AnError)

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(3), snapshot.MemoryGetTotal)
		assert.Equal(t, int64(1), snapshot.MemoryGetErrors)
	})

	t.Run("record delete operations", func(t *testing.T) {
		m := NewMetrics()

		// 先记录一些保存
		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)

		// 然后删除
		m.RecordDelete("user:123")

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(1), snapshot.MemoryDeleteTotal)
		assert.Equal(t, int64(1), snapshot.MemoryByNamespace["user:123"]) // 2 - 1 = 1
	})

	t.Run("record event processing", func(t *testing.T) {
		m := NewMetrics()

		m.RecordEventProcess("user_revision", 50*time.Millisecond, nil)
		m.RecordEventProcess("user_feedback", 30*time.Millisecond, nil)
		m.RecordEventProcess("user_message", 40*time.Millisecond, assert.AnError)

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(3), snapshot.EventProcessTotal)
		assert.Equal(t, int64(1), snapshot.EventProcessErrors)
	})

	t.Run("memory by type and scope", func(t *testing.T) {
		m := NewMetrics()

		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "behavior", ScopeSession, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeGlobal, 10*time.Millisecond, nil)

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(3), snapshot.MemoryByType["preference"])
		assert.Equal(t, int64(1), snapshot.MemoryByType["behavior"])
		assert.Equal(t, int64(2), snapshot.MemoryByScope[ScopeUser])
		assert.Equal(t, int64(1), snapshot.MemoryByScope[ScopeSession])
		assert.Equal(t, int64(1), snapshot.MemoryByScope[ScopeGlobal])
	})

	t.Run("average duration calculation", func(t *testing.T) {
		m := NewMetrics()

		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 20*time.Millisecond, nil)
		m.RecordSave("user:123", "preference", ScopeUser, 30*time.Millisecond, nil)

		snapshot := m.GetSnapshot()
		// 平均应该是 20ms
		assert.InDelta(t, 20*time.Millisecond, snapshot.AvgSaveDuration, float64(5*time.Millisecond))
	})

	t.Run("reset metrics", func(t *testing.T) {
		m := NewMetrics()

		m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
		m.RecordGet(5*time.Millisecond, nil)

		m.Reset()

		snapshot := m.GetSnapshot()
		assert.Equal(t, int64(0), snapshot.MemorySaveTotal)
		assert.Equal(t, int64(0), snapshot.MemoryGetTotal)
		assert.Empty(t, snapshot.MemoryByNamespace)
	})
}

func TestMetricsSnapshot(t *testing.T) {
	t.Run("total memories", func(t *testing.T) {
		snapshot := &MetricsSnapshot{
			MemoryByNamespace: map[string]int64{
				"user:123": 10,
				"user:456": 5,
				"team:789": 3,
			},
		}

		assert.Equal(t, int64(18), snapshot.TotalMemories())
	})

	t.Run("save error rate", func(t *testing.T) {
		snapshot := &MetricsSnapshot{
			MemorySaveTotal:  100,
			MemorySaveErrors: 5,
		}

		assert.InDelta(t, 0.05, snapshot.SaveErrorRate(), 0.001)
	})

	t.Run("get error rate", func(t *testing.T) {
		snapshot := &MetricsSnapshot{
			MemoryGetTotal:  200,
			MemoryGetErrors: 10,
		}

		assert.InDelta(t, 0.05, snapshot.GetErrorRate(), 0.001)
	})

	t.Run("event process error rate", func(t *testing.T) {
		snapshot := &MetricsSnapshot{
			EventProcessTotal:  50,
			EventProcessErrors: 5,
		}

		assert.InDelta(t, 0.1, snapshot.EventProcessErrorRate(), 0.001)
	})

	t.Run("zero division protection", func(t *testing.T) {
		snapshot := &MetricsSnapshot{}

		assert.Equal(t, 0.0, snapshot.SaveErrorRate())
		assert.Equal(t, 0.0, snapshot.GetErrorRate())
		assert.Equal(t, 0.0, snapshot.EventProcessErrorRate())
	})
}

func TestPrometheusExport(t *testing.T) {
	m := NewMetrics()

	// 记录一些数据
	m.RecordSave("user:123", "preference", ScopeUser, 10*time.Millisecond, nil)
	m.RecordGet(5*time.Millisecond, nil)
	m.RecordEventProcess("user_revision", 20*time.Millisecond, nil)

	// 创建一个模拟的 exporter
	exporter := &mockPrometheusExporter{
		gauges:     make(map[string]float64),
		counters:   make(map[string]float64),
		histograms: make(map[string]float64),
	}

	// 导出
	m.ExportToPrometheus(exporter)

	// 验证
	assert.Equal(t, float64(1), exporter.counters["logic_memory_save_total"])
	assert.Equal(t, float64(1), exporter.counters["logic_memory_get_total"])
	assert.Equal(t, float64(1), exporter.counters["logic_memory_event_process_total"])
}

// mockPrometheusExporter 模拟的 Prometheus 导出器
type mockPrometheusExporter struct {
	gauges     map[string]float64
	counters   map[string]float64
	histograms map[string]float64
}

func (m *mockPrometheusExporter) ExportGauge(name string, value float64, labels map[string]string) {
	m.gauges[name] = value
}

func (m *mockPrometheusExporter) ExportCounter(name string, value float64, labels map[string]string) {
	m.counters[name] = value
}

func (m *mockPrometheusExporter) ExportHistogram(name string, value float64, labels map[string]string) {
	m.histograms[name] = value
}
