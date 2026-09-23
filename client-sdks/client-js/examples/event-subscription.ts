/**
 * 事件订阅使用示例
 * 展示如何使用 nomad 的三通道事件系统
 */

import { WebSocketClient, SubscriptionManager, isProgressEvent, isControlEvent, isMonitorEvent, isEventType } from "@nomad/client-js";

async function main() {
  // 1. 创建 WebSocket 客户端
  const ws = new WebSocketClient({
    maxReconnectAttempts: 5,
    reconnectDelay: 1000,
    heartbeatInterval: 30000,
    heartbeatTimeout: 10000,
  });

  // 2. 连接到 nomad 服务器
  try {
    await ws.connect("ws://localhost:8080/ws");
    console.log("✅ Connected to nomad");
  } catch (error) {
    console.error("❌ Connection failed:", error);
    return;
  }

  // 3. 创建订阅管理器
  const subscriptionManager = new SubscriptionManager(ws);

  // 4. 订阅所有三个通道
  const subscription = subscriptionManager.subscribe(["progress", "control", "monitor"], {
    agentId: "agent-123",
    eventTypes: ["thinking", "text_chunk", "tool_start", "token_usage"],
  });

  // 5. 处理事件
  try {
    for await (const envelope of subscription) {
      const event = envelope.event;

      // 按通道分类处理
      if (isProgressEvent(event)) {
        handleProgressEvent(event);
      } else if (isControlEvent(event)) {
        handleControlEvent(event);
      } else if (isMonitorEvent(event)) {
        handleMonitorEvent(event);
      }
    }
  } catch (error) {
    console.error("❌ Event subscription error:", error);
  }

  // 6. 清理
  subscription.unsubscribe();
  ws.disconnect();
}

/**
 * 处理 Progress Channel 事件
 */
function handleProgressEvent(event: any) {
  if (isEventType(event, "thinking")) {
    console.log("🤔 AI 正在思考:", event.data.content);
  } else if (isEventType(event, "text_chunk")) {
    process.stdout.write(event.data.delta);
  } else if (isEventType(event, "tool_start")) {
    console.log("🔧 调用工具:", event.data.toolName);
  } else if (isEventType(event, "tool_end")) {
    console.log("✅ 工具完成:", event.data.toolName, "结果:", event.data.result);
  } else if (isEventType(event, "done")) {
    console.log("\n\n✅ 任务完成:", event.data.text);
  } else if (isEventType(event, "error")) {
    console.error("❌ 错误:", event.data.error);
  }
}

/**
 * 处理 Control Channel 事件
 */
function handleControlEvent(event: any) {
  if (isEventType(event, "tool_approval_request")) {
    console.log("⚠️  需要审批工具:", event.data.toolName);
    console.log("   审批 ID:", event.data.approvalId);
    console.log("   参数:", event.data.params);
    // 这里可以调用 API 批准或拒绝
  } else if (isEventType(event, "pause")) {
    console.log("⏸️  执行暂停:", event.data.reason);
  } else if (isEventType(event, "resume")) {
    console.log("▶️  执行恢复:", event.data.timestamp);
  }
}

/**
 * 处理 Monitor Channel 事件
 */
function handleMonitorEvent(event: any) {
  if (isEventType(event, "token_usage")) {
    console.log("📊 Token 使用:", {
      prompt: event.data.promptTokens,
      completion: event.data.completionTokens,
      total: event.data.totalTokens,
    });
  } else if (isEventType(event, "latency")) {
    console.log("⏱️  延迟:", event.data.latencyMs, "ms", "操作:", event.data.operation);
  } else if (isEventType(event, "cost")) {
    console.log("💰 成本:", event.data.cost, event.data.currency);
  } else if (isEventType(event, "compliance")) {
    const status = event.data.passed ? "✅" : "❌";
    console.log(`${status} 合规检查:`, event.data.details);
  }
}

// 运行示例
main().catch(console.error);
