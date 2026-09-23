package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestToolHandlers 测试 Tool 相关的处理器
func TestToolHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	var toolID string

	t.Run("CreateTool", func(t *testing.T) {
		body := `{
			"name": "test_tool",
			"description": "A test tool",
			"type": "custom",
			"schema": {
				"type": "object",
				"properties": {
					"input": {"type": "string"}
				}
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/v1/tools", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data := resp["data"].(map[string]any)
		toolID = data["id"].(string)
		assert.NotEmpty(t, toolID)
	})

	t.Run("ListTools", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/tools", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("GetTool", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/tools/"+toolID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), toolID)
	})

	t.Run("UpdateTool", func(t *testing.T) {
		body := `{
			"description": "Updated test tool"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/v1/tools/"+toolID, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteTool", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/v1/tools/"+toolID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("GetDeletedTool", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/tools/"+toolID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestMiddlewareHandlers 测试 Middleware 相关的处理器
func TestMiddlewareHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	var middlewareID string

	t.Run("CreateMiddleware", func(t *testing.T) {
		body := `{
			"name": "test_middleware",
			"type": "custom",
			"description": "A test middleware",
			"priority": 10
		}`

		req := httptest.NewRequest(http.MethodPost, "/v1/middlewares", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		data := resp["data"].(map[string]any)
		middlewareID = data["id"].(string)
		assert.NotEmpty(t, middlewareID)
	})

	t.Run("ListMiddlewares", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/middlewares", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMiddleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/middlewares/"+middlewareID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListMiddlewareRegistry", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/middlewares/registry", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Registry may return different structure
	})

	t.Run("DeleteMiddleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/v1/middlewares/"+middlewareID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

// TestTelemetryHandlers 测试 Telemetry 相关的处理器
func TestTelemetryHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("RecordMetric", func(t *testing.T) {
		body := `{
			"name": "test_metric",
			"type": "counter",
			"value": 42.0,
			"tags": {
				"env": "test"
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/v1/telemetry/metrics", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("ListMetrics", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/telemetry/metrics", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("RecordTrace", func(t *testing.T) {
		t.Skip("Trace recording needs proper implementation")
	})

	t.Run("RecordLog", func(t *testing.T) {
		body := `{
			"level": "info",
			"message": "Test log message",
			"timestamp": "2024-01-01T00:00:00Z"
		}`

		req := httptest.NewRequest(http.MethodPost, "/v1/telemetry/logs", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

// TestSystemHandlers 测试 System 相关的处理器
func TestSystemHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("GetSystemInfo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/system/info", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "version")
	})

	t.Run("GetSystemHealth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/system/health", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetSystemStats", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/system/stats", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListConfig", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/system/config", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateConfig", func(t *testing.T) {
		body := `{
			"value": "test_value"
		}`

		req := httptest.NewRequest(http.MethodPut, "/v1/system/config/test_key", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetConfig", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/system/config/test_key", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteConfig", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/v1/system/config/test_key", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

// TestMCPHandlers 测试 MCP 相关的处理器
func TestMCPHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	var serverID string

	t.Run("CreateMCPServer", func(t *testing.T) {
		body := `{
			"name": "test_mcp_server",
			"url": "http://localhost:8080",
			"type": "http"
		}`

		req := httptest.NewRequest(http.MethodPost, "/v1/mcp/servers", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		data := resp["data"].(map[string]any)
		serverID = data["id"].(string)
		assert.NotEmpty(t, serverID)
	})

	t.Run("ListMCPServers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/mcp/servers", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMCPServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/mcp/servers/"+serverID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateMCPServer", func(t *testing.T) {
		body := `{
			"name": "updated_mcp_server"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/v1/mcp/servers/"+serverID, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteMCPServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/v1/mcp/servers/"+serverID, nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

// TestEvalHandlers 测试 Eval 相关的处理器
func TestEvalHandlers(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("RunTextEval", func(t *testing.T) {
		t.Skip("Eval implementation needs proper setup")
	})

	t.Run("ListEvals", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/eval/evals", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateBenchmark", func(t *testing.T) {
		t.Skip("Benchmark creation needs proper setup")
	})

	t.Run("ListBenchmarks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/eval/benchmarks", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestValidationErrors 测试各种验证错误
func TestValidationErrors(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "EmptyAgentName",
			method:         http.MethodPost,
			path:           "/v1/agents",
			body:           `{"template_id": ""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidToolSchema",
			method:         http.MethodPost,
			path:           "/v1/tools",
			body:           `{"name": "test", "schema": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidRoomName",
			method:         http.MethodPost,
			path:           "/v1/rooms",
			body:           `{"name": ""}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			srv.Router().ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestContentTypeValidation 测试 Content-Type 验证
func TestContentTypeValidation(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("MissingContentType", func(t *testing.T) {
		body := `{"template_id": "chat"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/agents", strings.NewReader(body))
		// 不设置 Content-Type
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		// Gin 默认会尝试解析，但可能失败
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("WrongContentType", func(t *testing.T) {
		t.Skip("Content-Type validation varies by implementation")
	})
}

// TestQueryParameters 测试查询参数
func TestQueryParameters(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("ListAgentsWithPrefix", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/pool/agents?prefix=test", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListWithPagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/agents?page=1&limit=10", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestResponseFormat 测试响应格式
func TestResponseFormat(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("SuccessResponseFormat", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/rooms", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Contains(t, resp, "success")
		assert.True(t, resp["success"].(bool))
		assert.Contains(t, resp, "data")
	})

	t.Run("ErrorResponseFormat", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/agents/nonexistent", nil)
		w := httptest.NewRecorder()

		srv.Router().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Contains(t, resp, "success")
		assert.False(t, resp["success"].(bool))
		assert.Contains(t, resp, "error")

		errorObj := resp["error"].(map[string]any)
		assert.Contains(t, errorObj, "code")
		assert.Contains(t, errorObj, "message")
	})
}
