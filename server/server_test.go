package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/router"
	"github.com/nabob-ai/nomad/pkg/sandbox"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*Server, func()) {
	// Create test store
	st, err := store.NewJSONStore(t.TempDir())
	require.NoError(t, err)

	// Create dependencies
	toolRegistry := tools.NewRegistry()
	sandboxFactory := sandbox.NewFactory()
	providerFactory := NewMockProviderFactory()
	templateRegistry := agent.NewTemplateRegistry()

	// Register test templates
	chatTemplate := &types.AgentTemplateDefinition{
		ID:           "chat",
		Version:      "1.0.0",
		SystemPrompt: "You are a helpful assistant.",
		Model:        "test-model",
		Tools:        []string{},
	}
	templateRegistry.Register(chatTemplate)

	defaultModel := &types.ModelConfig{
		Provider: "mock",
		Model:    "test-model",
	}
	routes := []router.StaticRouteEntry{
		{Task: "chat", Priority: router.PriorityQuality, Model: defaultModel},
	}
	rt := router.NewStaticRouter(defaultModel, routes)

	agentDeps := &agent.Dependencies{
		Store:            st,
		ToolRegistry:     toolRegistry,
		SandboxFactory:   sandboxFactory,
		ProviderFactory:  providerFactory,
		TemplateRegistry: templateRegistry,
		Router:           rt,
	}

	deps := &Dependencies{
		Store:     st,
		AgentDeps: agentDeps,
	}

	// Create server with test config
	config := DefaultConfig()
	config.Auth.APIKey.Enabled = false // Disable auth for tests

	srv, err := New(config, deps)
	require.NoError(t, err)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Stop(ctx)
	}

	return srv, cleanup
}

func TestServerHealth(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestServerMetrics(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServerSystemInfo(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/system/info", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "version")
	assert.Contains(t, w.Body.String(), "go_version")
}

func TestServerPoolStats(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/pool/stats", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "total_agents")
}

func TestServerRoomsList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/rooms", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Agent Handler Tests

func TestAgentCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"template_id": "chat",
		"name": "Test Agent",
		"model_config": {
			"provider": "mock",
			"model": "test-model"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/agents", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "id")
}

func TestAgentList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/agents", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestAgentGetNotFound(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/agents/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not_found")
}

// Pool Handler Tests

func TestPoolListAgents(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/pool/agents", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "agents")
}

func TestPoolGetAgentNotFound(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/pool/agents/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not_found")
}

// Room Handler Tests

func TestRoomCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"name": "Test Room",
		"metadata": {
			"description": "A test room"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "Test Room")
}

func TestRoomGetNotFound(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/rooms/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not_found")
}

// Memory Handler Tests

func TestMemoryCreateWorking(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"key": "test_key",
		"value": "test_value",
		"type": "string",
		"ttl": 3600
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/memory/working", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "test_key")
}

func TestMemoryListWorking(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/memory/working", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Session Handler Tests

func TestSessionCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"agent_id": "test-agent",
		"context": {
			"user_id": "user-123"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/sessions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "test-agent")
}

func TestSessionList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/sessions", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Workflow Handler Tests

func TestWorkflowCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"name": "Test Workflow",
		"description": "A test workflow",
		"version": "1.0.0",
		"steps": [
			{
				"id": "step1",
				"name": "First Step",
				"type": "agent"
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/workflows", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "Test Workflow")
}

func TestWorkflowList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/workflows", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Tool Handler Tests

func TestToolCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

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
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "test_tool")
}

func TestToolList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/tools", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Middleware Handler Tests

func TestMiddlewareCreate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

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
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "test_middleware")
}

func TestMiddlewareList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/middlewares", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Telemetry Handler Tests

func TestTelemetryRecordMetric(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

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
	assert.Contains(t, w.Body.String(), "success")
}

func TestTelemetryListMetrics(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/telemetry/metrics", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// System Handler Tests

func TestSystemConfigList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/system/config", nil)
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestSystemConfigUpdate(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{
		"value": "test_value"
	}`

	req := httptest.NewRequest(http.MethodPut, "/v1/system/config/test_key", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
