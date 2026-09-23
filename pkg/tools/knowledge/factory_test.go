package knowledge

import (
	"context"
	"testing"

	"github.com/nabob-ai/nomad/pkg/knowledge/core"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/vector"
)

func TestFactory_AddAndSearch(t *testing.T) {
	pipe, err := core.NewPipeline(core.PipelineConfig{
		Store:    vector.NewMemoryStore(),
		Embedder: vector.NewMockEmbedder(16),
	})
	if err != nil {
		t.Fatalf("new pipeline: %v", err)
	}

	f := NewFactory(pipe)
	add, _ := f.KnowledgeAddTool()
	search, _ := f.KnowledgeSearchTool()

	_, err = add.Execute(context.Background(), map[string]any{
		"text":      "Go has goroutines for concurrency.",
		"namespace": "ns1",
	}, &tools.ToolContext{})
	if err != nil {
		t.Fatalf("add exec error: %v", err)
	}

	resp, err := search.Execute(context.Background(), map[string]any{
		"query":     "goroutines",
		"namespace": "ns1",
		"top_k":     3,
	}, &tools.ToolContext{})
	if err != nil {
		t.Fatalf("search exec error: %v", err)
	}

	resMap, ok := resp.(map[string]any)
	if !ok {
		t.Fatalf("unexpected resp type: %T", resp)
	}
	results, ok := resMap["results"].([]map[string]any)
	if !ok || len(results) == 0 {
		t.Fatalf("expected results")
	}
}
