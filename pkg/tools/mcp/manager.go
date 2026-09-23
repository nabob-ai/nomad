package mcp

import (
	"context"
	"fmt"
	"sync"

	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/tools/search"
)

// MCPManager MCP Server 管理器
// 管理多个 MCP Server 连接和工具注册
type MCPManager struct {
	mu       sync.RWMutex
	servers  map[string]*MCPServer
	registry *tools.Registry
}

// NewMCPManager 创建 MCP Manager
func NewMCPManager(registry *tools.Registry) *MCPManager {
	return &MCPManager{
		servers:  make(map[string]*MCPServer),
		registry: registry,
	}
}

// AddServer 添加 MCP Server
func (m *MCPManager) AddServer(config *MCPServerConfig) (*MCPServer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.servers[config.ServerID]; exists {
		return nil, fmt.Errorf("server already exists: %s", config.ServerID)
	}

	// 创建 Server
	server, err := NewMCPServer(config, m.registry)
	if err != nil {
		return nil, fmt.Errorf("create mcp server: %w", err)
	}

	m.servers[config.ServerID] = server
	return server, nil
}

// ConnectServer 连接指定的 MCP Server 并注册工具
func (m *MCPManager) ConnectServer(ctx context.Context, serverID string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	// 连接并发现工具
	if err := server.Connect(ctx); err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	// 注册工具到 Registry
	if err := server.RegisterTools(); err != nil {
		return fmt.Errorf("register tools: %w", err)
	}

	return nil
}

// ConnectAll 连接所有已添加的 MCP Server
func (m *MCPManager) ConnectAll(ctx context.Context) error {
	m.mu.RLock()
	serverIDs := make([]string, 0, len(m.servers))
	for id := range m.servers {
		serverIDs = append(serverIDs, id)
	}
	m.mu.RUnlock()

	// 连接所有 Server
	for _, serverID := range serverIDs {
		if err := m.ConnectServer(ctx, serverID); err != nil {
			return fmt.Errorf("connect server %s: %w", serverID, err)
		}
	}

	return nil
}

// GetServer 获取指定的 MCP Server
func (m *MCPManager) GetServer(serverID string) (*MCPServer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, exists := m.servers[serverID]
	return server, exists
}

// ListServers 列出所有 Server ID
func (m *MCPManager) ListServers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.servers))
	for id := range m.servers {
		ids = append(ids, id)
	}
	return ids
}

// GetServerCount 获取 Server 数量
func (m *MCPManager) GetServerCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.servers)
}

// GetTotalToolCount 获取所有 Server 提供的工具总数
func (m *MCPManager) GetTotalToolCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, server := range m.servers {
		count += server.GetToolCount()
	}
	return count
}

// RemoveServer 移除 MCP Server
func (m *MCPManager) RemoveServer(serverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[serverID]; !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	delete(m.servers, serverID)
	return nil
}

// ConnectServerDeferred 连接 MCP Server 但使用延迟加载模式
// 只发现工具并添加到索引，不立即注册到 Registry
func (m *MCPManager) ConnectServerDeferred(ctx context.Context, serverID string, index *search.ToolIndex) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	// 连接并发现工具
	if err := server.Connect(ctx); err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	// 将工具添加到索引（延迟加载模式）
	if err := server.IndexToolsToIndex(index); err != nil {
		return fmt.Errorf("index tools: %w", err)
	}

	return nil
}

// ConnectAllDeferred 连接所有 MCP Server 使用延迟加载模式
func (m *MCPManager) ConnectAllDeferred(ctx context.Context, index *search.ToolIndex) error {
	m.mu.RLock()
	serverIDs := make([]string, 0, len(m.servers))
	for id := range m.servers {
		serverIDs = append(serverIDs, id)
	}
	m.mu.RUnlock()

	// 连接所有 Server（延迟模式）
	for _, serverID := range serverIDs {
		if err := m.ConnectServerDeferred(ctx, serverID, index); err != nil {
			return fmt.Errorf("connect server %s: %w", serverID, err)
		}
	}

	return nil
}

// ActivateTool 激活延迟加载的 MCP 工具
func (m *MCPManager) ActivateTool(toolName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 遍历所有服务器查找工具
	for _, server := range m.servers {
		if err := server.RegisterToolDeferred(toolName); err == nil {
			return nil // 成功注册
		}
	}

	return fmt.Errorf("tool not found in any server: %s", toolName)
}

// GetAllToolIndexEntries 获取所有服务器的工具索引条目
func (m *MCPManager) GetAllToolIndexEntries() []search.ToolIndexEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]search.ToolIndexEntry, 0)
	for _, server := range m.servers {
		entries = append(entries, server.GetToolIndexEntries()...)
	}
	return entries
}
