/**
 * useChat Composable (重构版)
 *
 * 使用 Pinia stores 管理状态，简化逻辑
 */

import { computed, onMounted, reactive } from "vue";
import type { ChatConfig, TextMessage } from "@/types";
import { useNomadClient } from "./useNomadClient";
import { useWebSocket } from "./useWebSocket";
import { generateId } from "@/utils/format";

// 导入 Pinia stores
import { useChatStore, useThinkingStore, useToolsStore, useTodosStore, useApprovalStore } from "@/stores";

export function useChat(config: ChatConfig) {
  // ==================
  // Stores
  // ==================

  const chatStore = useChatStore();
  const thinkingStore = useThinkingStore();
  const toolsStore = useToolsStore();
  const todosStore = useTodosStore();
  const approvalStore = useApprovalStore();

  // ==================
  // 配置
  // ==================

  const isDemoMode = config.demoMode ?? true;
  const apiUrl = config.apiUrl || import.meta.env.VITE_API_URL || "http://localhost:8080";
  const wsUrlOverride = config.wsUrl || import.meta.env.VITE_WS_URL;

  // ==================
  // 客户端初始化
  // ==================

  const { client } = useNomadClient({
    baseUrl: apiUrl,
    apiKey: config.apiKey,
    wsUrl: wsUrlOverride,
  });

  const { connect, getInstance, isConnected: wsConnected } = useWebSocket();

  // 连接状态
  const connectionState = computed(() => (isDemoMode ? true : wsConnected.value));

  // ==================
  // Demo 模式相关
  // ==================

  let demoCursor = 0;

  const fallbackResponses = [
    "我已经为你生成了一个新的多 Agent 工作流，包含大纲、评价器和部署策略。",
    "Nomad 的沙箱已准备好，所有写入都被限制在 /workspace 目录，你可以放心执行指令。",
    "我为这个会话自动挂载了上下文记忆，后续可以直接引用历史工单。",
    "Streaming 模式已打开，等待后端返回 token，平均延迟 220ms。",
  ];

  const pickDemoResponse = (content: string) => {
    const list = config.demoResponses?.length ? config.demoResponses : fallbackResponses;
    const index = demoCursor % list.length;
    demoCursor += 1;
    const template = list[index] ?? "";
    return template.includes("{question}") ? template.split("{question}").join(content) : template;
  };

  // ==================
  // 事件处理器
  // ==================

  /**
   * 处理 Agent 事件（统一分发）
   */
  const handleAgentEvent = (type: string, ev: any, messageId?: string) => {
    const currentMessageId = messageId || chatStore.lastAssistantMessage?.id || "";

    // === 思维事件 ===
    if (type === "think_chunk_start") {
      thinkingStore.startThinking(currentMessageId);
      return;
    }

    if (type === "think_chunk") {
      thinkingStore.handleThinkChunk(ev.delta || "");
      return;
    }

    if (type === "think_chunk_end") {
      thinkingStore.endThinking();
      return;
    }

    // === 工具事件 ===
    if (type === "tool:start") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolStart(call);
      // 添加工具调用步骤到思维块
      if (currentMessageId) {
        thinkingStore.addToolCallStep(currentMessageId, call.name, call.arguments);
      }
      return;
    }

    if (type === "tool:progress") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolProgress(call.id || call.ID, ev.progress ?? call.progress ?? 0, ev.message, ev.metadata);
      return;
    }

    if (type === "tool:intermediate") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolIntermediate(call.id || call.ID, ev.label || "", ev.data);
      return;
    }

    if (type === "tool:end") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolEnd(call);
      // 添加工具结果步骤到思维块
      if (currentMessageId) {
        thinkingStore.addToolResultStep(currentMessageId, call.result);
      }
      return;
    }

    if (type === "tool:error") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolError(call.id || call.ID, ev.error || call.error);
      return;
    }

    if (type === "tool:cancelled") {
      const call = ev.Call || ev.call || {};
      toolsStore.handleToolCancelled(call.id || call.ID, ev.reason);
      return;
    }

    // === 审批事件 ===
    if (type === "permission_required") {
      const call = ev.call || {};
      approvalStore.addApprovalRequest({
        id: ev.request_id || generateId("approval"),
        messageId: currentMessageId,
        toolName: call.name || "",
        args: call.arguments || {},
        reason: ev.reason || "",
        timestamp: Date.now(),
      });

      // 添加审批步骤到思维块
      if (currentMessageId) {
        thinkingStore.addApprovalStep(currentMessageId, call.name, call.arguments);
      }
      return;
    }

    // === Todo 事件 ===
    if (type === "todo_update") {
      todosStore.updateTodos(ev.todos || []);
      return;
    }

    // === AskUser 事件 ===
    if (type === "ask_user") {
      // 添加 AskUser 消息到消息列表
      const askUserMsg = {
        id: generateId("msg"),
        type: "ask-user" as const,
        role: "assistant" as const,
        content: {
          request_id: ev.request_id,
          questions: ev.questions || [],
          answered: false,
        },
        createdAt: Date.now(),
      };
      chatStore.addMessage(askUserMsg);
      return;
    }

    // === 状态变更事件 ===
    if (type === "state_changed") {
      const state = ev.state;
      if (state === "working" || state === "running") {
        chatStore.updateAgentStatus("thinking");
      } else if (state === "idle" || state === "ready" || state === "completed") {
        chatStore.updateAgentStatus("idle");
      } else if (state === "failed") {
        chatStore.updateAgentStatus("error");
      }
      return;
    }

    // === Token 使用统计 ===
    if (type === "token_usage") {
      console.log("Token usage:", ev);
      return;
    }

    // === 错误事件 ===
    if (type === "error") {
      console.error("Agent error:", ev.message, ev.detail);
      return;
    }
  };

  // ==================
  // 发送消息
  // ==================

  const sendMessage = async (content: string) => {
    if (!content.trim()) return;

    // 创建用户消息
    const userMessage: TextMessage = {
      id: generateId("msg"),
      type: "text",
      role: "user",
      content: { text: content },
      createdAt: Date.now(),
      status: "pending",
    };
    chatStore.addMessage(userMessage);

    // 创建助手消息占位符（使用 reactive 确保响应式）
    const assistantMessage: TextMessage = reactive({
      id: generateId("msg"),
      type: "text",
      role: "assistant",
      content: { text: "" },
      createdAt: Date.now(),
    }) as TextMessage;
    chatStore.addMessage(assistantMessage);

    chatStore.setTyping(true);
    chatStore.updateAgentStatus("thinking");
    userMessage.status = "sent";
    chatStore.setCurrentInput("");

    try {
      if (isDemoMode) {
        // Demo 模式
        await new Promise((resolve) => setTimeout(resolve, config.demoDelay ?? 800));
        (assistantMessage.content as { text: string }).text = pickDemoResponse(content);
        assistantMessage.status = "sent";
        chatStore.setTyping(false);
        chatStore.updateAgentStatus("idle");
      } else {
        // WebSocket 流式模式
        const ws = getInstance();

        if (ws && wsConnected.value) {
          // 监听 WebSocket 消息
          const unsubscribe = ws.onMessage((message: any) => {
            if (message.type === "text_delta" && message.payload?.text) {
              // 使用批量更新优化性能
              chatStore.handleTextChunk(assistantMessage.id, message.payload.text);
            } else if (message.type === "chat_complete") {
              assistantMessage.status = "sent";
              chatStore.setTyping(false);
              chatStore.updateAgentStatus("idle");
              unsubscribe();

              // 触发回调
              if (config.onReceive) {
                config.onReceive(assistantMessage);
              }
            } else if (message.type === "error") {
              assistantMessage.content.text = `❌ ${message.payload?.message || "发送失败"}`;
              userMessage.status = "error";
              chatStore.setTyping(false);
              chatStore.updateAgentStatus("idle");
              unsubscribe();
              if (config.onError) {
                config.onError(new Error(message.payload?.message));
              }
            } else if (message.type === "agent_event") {
              const ev = message.payload?.event;
              const evType = message.payload?.type || ev?.type || ev?.EventType;
              if (ev && evType) {
                handleAgentEvent(evType, ev, assistantMessage.id);
              }
            }
          });

          // 发送聊天消息
          ws.send({
            type: "chat",
            payload: {
              template_id: config.agentId || "chat",
              input: content,
              model_config: config.modelConfig,
            },
          });
        } else {
          // HTTP Fallback
          const response = await client.agents.chatDirect(content, config.agentId || "chat");
          assistantMessage.content.text = response.text || response.data?.text || "无响应";
          assistantMessage.status = "sent";
          chatStore.setTyping(false);
          chatStore.updateAgentStatus("idle");
        }
      }
    } catch (error: any) {
      console.error("Send message error:", error);

      assistantMessage.content.text = `❌ 发送失败: ${error.message || "未知错误"}`;
      userMessage.status = "error";
      chatStore.setTyping(false);
      chatStore.updateAgentStatus("idle");

      if (config.onError) {
        config.onError(error);
      }
    }

    // 触发回调
    if (config.onSend) {
      config.onSend(userMessage);
    }
  };

  // ==================
  // 其他操作
  // ==================

  const answerQuestion = async (requestId: string, answers: Record<string, any>) => {
    const ws = getInstance();
    if (!ws || !wsConnected.value) return;

    ws.send({
      type: "user_answer",
      payload: {
        request_id: requestId,
        answers,
      },
    });

    // 更新消息状态
    const msgIndex = chatStore.messages.findIndex((m: any) => m.type === "ask-user" && m.content.request_id === requestId);
    if (msgIndex !== -1) {
      const msg = chatStore.messages[msgIndex] as any;
      msg.content.answered = true;
      msg.content.answers = answers;
    }
  };

  const controlTool = async (toolCallId: string, action: "cancel" | "pause" | "resume") => {
    const ws = getInstance();
    if (!ws || !wsConnected.value) return;
    ws.send({
      type: "tool:control",
      payload: {
        tool_call_id: toolCallId,
        action,
      },
    });
  };

  // ==================
  // 初始化
  // ==================

  onMounted(async () => {
    // 设置 Agent 信息
    if (config.agentId || config.agentProfile) {
      chatStore.setAgent({
        id: config.agentId || "demo-agent",
        name: config.agentProfile?.name || "Nomad Copilot",
        description: config.agentProfile?.description || "多模态执行、自动规划、符合企业安全的 Agent",
        avatar: config.agentProfile?.avatar,
        metadata: {
          model: "nomad:builder",
        },
      });
    }

    // 添加欢迎消息
    if (config.welcomeMessage && chatStore.messages.length === 0) {
      const welcomeText = typeof config.welcomeMessage === "string" ? config.welcomeMessage : config.welcomeMessage.type === "text" ? config.welcomeMessage.content.text : "👋 你好，我是 Nomad Copilot。";

      chatStore.addTextMessage("assistant", welcomeText);
    }

    // 初始化 WebSocket 连接
    if (!isDemoMode) {
      const wsUrl = wsUrlOverride || apiUrl.replace(/^http/, "ws") + "/v1/ws";
      console.log("🚀 Initializing WebSocket connection to:", wsUrl);
      try {
        await connect(wsUrl);
        console.log("✅ WebSocket initialized in useChat");
      } catch (error) {
        console.error("❌ Failed to initialize WebSocket:", error);
      }
    }
  });

  // ==================
  // Return (保持向后兼容的 API)
  // ==================

  return {
    // 状态（从 stores 导出）
    messages: computed(() => chatStore.messages),
    isTyping: computed(() => chatStore.isTyping),
    isConnected: wsConnected,
    connectionState,
    currentInput: computed(() => chatStore.currentInput),
    agent: computed(() => chatStore.agent),
    isThinking: computed(() => chatStore.isTyping),
    thinkingContent: computed(() => thinkingStore.currentThought),
    currentStep: computed(() => 0), // TODO: 与 workflow store 集成
    todos: computed(() => todosStore.todos),
    pendingAskUser: computed(() => null), // TODO: 需要实现

    // 方法
    sendMessage,
    sendImage: async (file: File) => {
      // TODO: 实现图片上传
      console.log("Send image:", file.name);
    },
    retryMessage: async (message: any) => {
      if (message.type === "text" && message.role === "user") {
        await sendMessage(message.content.text);
      }
    },
    deleteMessage: (messageId: string) => {
      chatStore.deleteMessage(messageId);
    },
    clearMessages: () => {
      chatStore.clearMessages();
      thinkingStore.clearAllSteps();
      toolsStore.clearAllTools();
      todosStore.clearAllTodos();
      approvalStore.clearPendingApprovals();
    },
    approveAction: (requestId: string) => {
      approvalStore.approve(requestId);
    },
    rejectAction: (requestId: string, reason?: string) => {
      approvalStore.reject(requestId, reason);
    },
    toolRunsList: computed(() => toolsStore.toolRunsList),
    controlTool,
    answerQuestion,
  };
}
