package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/types"
)

// ToolRuntimeHandler 提供工具调用运行态查询
type ToolRuntimeHandler struct {
	store *store.Store
	reg   *RuntimeAgentRegistry
}

// NewToolRuntimeHandler 创建新的运行态处理器
func NewToolRuntimeHandler(st store.Store, reg *RuntimeAgentRegistry) *ToolRuntimeHandler {
	return &ToolRuntimeHandler{store: &st, reg: reg}
}

// GetStatus 获取指定调用的状态
func (h *ToolRuntimeHandler) GetStatus(c *gin.Context) {
	agentID := c.Query("agent_id")
	callID := c.Param("id")
	if agentID == "" || callID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id and call id are required"})
		return
	}

	if snap := h.runtimeSnapshot(agentID, callID); snap != nil {
		c.JSON(http.StatusOK, gin.H{"status": snap})
		return
	}

	records, err := (*h.store).LoadToolCallRecords(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, rec := range records {
		if rec.ID == callID {
			c.JSON(http.StatusOK, gin.H{
				"status": mapRecordToSnapshot(&rec),
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "call not found"})
}

// ListRunning 返回运行中的工具调用
func (h *ToolRuntimeHandler) ListRunning(c *gin.Context) {
	agentID := c.Query("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id is required"})
		return
	}

	// 优先返回实时运行表
	if h.reg != nil {
		if ag := h.reg.Get(agentID); ag != nil {
			c.JSON(http.StatusOK, gin.H{"running": ag.ListRunningToolSnapshots()})
			return
		}
	}

	records, err := (*h.store).LoadToolCallRecords(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	running := make([]types.ToolCallSnapshot, 0)
	for _, rec := range records {
		if rec.State == types.ToolCallStateExecuting || rec.State == types.ToolCallStateQueued || rec.State == types.ToolCallStatePending || rec.State == types.ToolCallStateCancelling {
			snapshot := mapRecordToSnapshot(&rec)
			running = append(running, snapshot)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"running": running,
	})
}

// GetResult 导出工具执行结果
func (h *ToolRuntimeHandler) GetResult(c *gin.Context) {
	agentID := c.Query("agent_id")
	callID := c.Param("id")
	format := c.DefaultQuery("format", "json")
	if agentID == "" || callID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id and call id are required"})
		return
	}

	if snap := h.runtimeSnapshot(agentID, callID); snap != nil {
		switch format {
		case "json":
			c.JSON(http.StatusOK, gin.H{"result": snap.Result, "error": snap.Error})
		default:
			c.String(http.StatusOK, "%v", snap.Result)
		}
		return
	}

	records, err := (*h.store).LoadToolCallRecords(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, rec := range records {
		if rec.ID != callID {
			continue
		}
		switch format {
		case "json":
			c.JSON(http.StatusOK, gin.H{"result": rec.Result, "error": rec.Error})
		default:
			c.String(http.StatusOK, "%v", rec.Result)
		}
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "call not found"})
}

func (h *ToolRuntimeHandler) runtimeSnapshot(agentID, callID string) *types.ToolCallSnapshot {
	if h.reg == nil {
		return nil
	}
	ag := h.reg.Get(agentID)
	if ag == nil {
		return nil
	}
	snap := ag.GetToolSnapshot(callID)
	if snap.ID == "" {
		return nil
	}
	return &snap
}

func mapRecordToSnapshot(rec *types.ToolCallRecord) types.ToolCallSnapshot {
	return types.ToolCallSnapshot{
		ID:           rec.ID,
		Name:         rec.Name,
		State:        rec.State,
		Arguments:    rec.Input,
		Result:       rec.Result,
		Error:        rec.Error,
		Progress:     rec.Progress,
		Intermediate: rec.Intermediate,
		StartedAt:    rec.StartTime,
		UpdatedAt:    rec.UpdatedAt,
	}
}
