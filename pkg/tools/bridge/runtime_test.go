package bridge

import (
	"context"
	"testing"
	"time"
)

func TestPythonRuntime_Execute(t *testing.T) {
	runtime := NewPythonRuntime(nil)

	if !runtime.IsAvailable() {
		t.Skip("Python not available")
	}

	ctx := context.Background()
	code := `
result = _input['a'] + _input['b']
print(result)
`
	input := map[string]any{
		"a": 10,
		"b": 20,
	}

	result, err := runtime.Execute(ctx, code, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	// 输出应该是 30
	if result.Output != "30" {
		t.Errorf("expected output '30', got %v", result.Output)
	}
}

func TestPythonRuntime_JSONOutput(t *testing.T) {
	runtime := NewPythonRuntime(nil)

	if !runtime.IsAvailable() {
		t.Skip("Python not available")
	}

	ctx := context.Background()
	code := `
import json
result = {"sum": _input['a'] + _input['b'], "product": _input['a'] * _input['b']}
print(json.dumps(result))
`
	input := map[string]any{
		"a": float64(5),
		"b": float64(3),
	}

	result, err := runtime.Execute(ctx, code, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	// 输出应该被解析为 JSON
	output, ok := result.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", result.Output)
	}

	if output["sum"] != float64(8) {
		t.Errorf("expected sum=8, got %v", output["sum"])
	}
}

func TestNodeJSRuntime_Execute(t *testing.T) {
	runtime := NewNodeJSRuntime(nil)

	if !runtime.IsAvailable() {
		t.Skip("Node.js not available")
	}

	ctx := context.Background()
	code := `
const result = _input.a + _input.b;
console.log(result);
`
	input := map[string]any{
		"a": 10,
		"b": 20,
	}

	result, err := runtime.Execute(ctx, code, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	if result.Output != "30" {
		t.Errorf("expected output '30', got %v", result.Output)
	}
}

func TestBashRuntime_Execute(t *testing.T) {
	runtime := NewBashRuntime(nil)

	if !runtime.IsAvailable() {
		t.Skip("Bash not available")
	}

	ctx := context.Background()
	code := `echo "Hello, World!"`
	input := map[string]any{}

	result, err := runtime.Execute(ctx, code, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	if result.Output != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %v", result.Output)
	}
}

func TestRuntimeManager_Execute(t *testing.T) {
	manager := NewRuntimeManager(nil)

	// 检查可用语言
	langs := manager.AvailableLanguages()
	if len(langs) == 0 {
		t.Skip("No runtimes available")
	}

	ctx := context.Background()

	// 测试每个可用的运行时
	for _, lang := range langs {
		t.Run(string(lang), func(t *testing.T) {
			var code string
			switch lang {
			case LangPython:
				code = "print('test')"
			case LangNodeJS:
				code = "console.log('test')"
			case LangBash:
				code = "echo test"
			}

			result, err := manager.Execute(ctx, lang, code, nil)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			if !result.Success {
				t.Errorf("expected success, got error: %s", result.Error)
			}
		})
	}
}

func TestRuntime_Timeout(t *testing.T) {
	config := &RuntimeConfig{
		Timeout:   1 * time.Second,
		MaxOutput: 1024,
	}

	runtime := NewPythonRuntime(config)

	if !runtime.IsAvailable() {
		t.Skip("Python not available")
	}

	ctx := context.Background()
	code := `
import time
time.sleep(10)
print("done")
`

	result, err := runtime.Execute(ctx, code, nil)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	// 应该超时失败
	if result.Success {
		t.Error("expected timeout failure")
	}

	if result.Error != "execution timeout" {
		t.Errorf("expected timeout error, got: %s", result.Error)
	}
}

func TestRuntime_SyntaxError(t *testing.T) {
	runtime := NewPythonRuntime(nil)

	if !runtime.IsAvailable() {
		t.Skip("Python not available")
	}

	ctx := context.Background()
	code := `
this is not valid python
`

	result, err := runtime.Execute(ctx, code, nil)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	// 应该失败
	if result.Success {
		t.Error("expected syntax error failure")
	}

	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filename string
		expected Language
	}{
		{"script.py", LangPython},
		{"app.js", LangNodeJS},
		{"app.mjs", LangNodeJS},
		{"script.sh", LangBash},
		{"script.bash", LangBash},
		{"file.txt", ""},
		{"file.go", ""},
	}

	for _, tt := range tests {
		got := DetectLanguage(tt.filename)
		if got != tt.expected {
			t.Errorf("DetectLanguage(%q) = %q, want %q", tt.filename, got, tt.expected)
		}
	}
}
