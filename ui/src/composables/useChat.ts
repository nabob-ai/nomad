/**
 * useChat Composable
 * 管理 Chat 对话逻辑
 */

import { ref, onMounted, reactive, computed } from "vue";
import type { Message, ChatConfig, TextMessage, Agent, AskUserMessage, ToolCallSnapshot, TodoItemData, Question } from "@/types";
import { useNomadClient } from "./useNomadClient";
import { useWebSocket } from "./useWebSocket";
import { generateId } from "@/utils/format";
import { useChatStore } from "@/stores/chat";
import { useThinkingStore } from "@/stores/thinking";
import { useToolsStore } from "@/stores/tools";
import { useTodosStore } from "@/stores/todos";
import { useApprovalStore } from "@/stores/approval";
import { useWorkflowStore } from "@/stores/workflow";

export function useChat(config: ChatConfig) {
  // 初始化 Pinia Stores
  const chatStore = useChatStore();
  const thinkingStore = useThinkingStore();
  const toolsStore = useToolsStore();
  const todosStore = useTodosStore();
  const approvalStore = useApprovalStore();
  const workflowStore = useWorkflowStore();

  // 保留部分本地状态 (不适合放在 store 中的)
  const currentInput = ref("");
  const demoConnection = ref(true);
  const isDemoMode = config.demoMode ?? true;
  const pendingAskUser = ref<{ requestId: string; questions: Question[] } | null>(null);
  const agent = ref<Agent>({
    id: config.agentId || "demo-agent",
    name: config.agentProfile?.name || "Nomad Copilot",
    description: config.agentProfile?.description || "多模态执行、自动规划、符合企业安全的 Agent",
    avatar: config.agentProfile?.avatar,
    status: "idle",
    metadata: {
      model: "nomad:builder",
    },
  });
  const demoCursor = ref(0);

  const apiUrl = config.apiUrl || import.meta.env.VITE_API_URL || "http://localhost:8080";
  const wsUrlOverride = config.wsUrl || import.meta.env.VITE_WS_URL;

  const { client } = useNomadClient({
    baseUrl: apiUrl,
    apiKey: config.apiKey,
    wsUrl: wsUrlOverride,
  });

  const { connect, getInstance, isConnected: wsConnected } = useWebSocket();
  // connectionState 用于组件中判断连接状态
  const connectionState = computed(() => (isDemoMode ? demoConnection.value : wsConnected.value));

  // 初始化 WebSocket 连接
  onMounted(async () => {
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

  const fallbackResponses = [
    "我已经为你生成了一个新的多 Agent 工作流，包含大纲、评价器和部署策略。",
    "Nomad 的沙箱已准备好，所有写入都被限制在 /workspace 目录，你可以放心执行指令。",
    "我为这个会话自动挂载了上下文记忆，后续可以直接引用历史工单。",
    "Streaming 模式已打开，等待后端返回 token，平均延迟 220ms。",
  ];

  const pickDemoResponse = (content: string): string => {
    const list = config.demoResponses?.length ? config.demoResponses : fallbackResponses;
    const index = demoCursor.value % list.length;
    demoCursor.value += 1;
    const template = list[index] ?? "";
    return template.includes("{question}") ? template.split("{question}").join(content) : template;
  };

  // 发送消息
  const sendMessage = async (content: string) => {
    console.log("📤 sendMessage called with:", content);
    console.log("📊 isDemoMode:", isDemoMode);
    console.log("📊 wsConnected:", wsConnected.value);
    console.log("📊 ws instance:", getInstance());

    if (!content.trim()) return;

    // 添加用户消息
    const userMessage: TextMessage = {
      id: generateId("msg"),
      type: "text",
      role: "user",
      content: { text: content },
      createdAt: Date.now(),
      status: "pending",
    };
    chatStore.messages.push(userMessage);
    console.log("✅ User message added to messages array");

    // 创建 AI 响应占位（使用 reactive 确保响应式）
    const assistantMessage: TextMessage = reactive({
      id: generateId("msg"),
      type: "text",
      role: "assistant",
      content: { text: "" },
      createdAt: Date.now(),
    }) as TextMessage;
    chatStore.messages.push(assistantMessage);
    chatStore.setActiveMessage(assistantMessage.id);
    console.log("✅ Assistant message placeholder added");

    chatStore.isTyping = true;
    agent.value.status = "thinking";
    userMessage.status = "sent";
    currentInput.value = "";

    try {
      if (isDemoMode) {
        await new Promise((resolve) => setTimeout(resolve, config.demoDelay ?? 800));
        assistantMessage.content.text = pickDemoResponse(content);
        assistantMessage.status = "sent";
        chatStore.isTyping = false;
        agent.value.status = "idle";
      } else {
        const ws = getInstance();
        console.log("🔍 Checking WebSocket availability:", {
          "ws exists": !!ws,
          isConnected: wsConnected.value,
          "ws type": ws?.constructor?.name,
        });

        // 使用 WebSocket 进行流式对话
        if (ws && wsConnected.value) {
          console.log("✅ Using WebSocket for chat");

          // 监听 WebSocket 消息
          const unsubscribe = ws.onMessage((message: any) => {
            console.log("📥 WebSocket message:", message);

            if (message.type === "text_delta" && message.payload?.text) {
              assistantMessage.content.text += message.payload.text;
              console.log("📝 Updated text:", assistantMessage.content.text.substring(0, 50) + "...");
            } else if (message.type === "chat_complete") {
              assistantMessage.status = "sent";
              chatStore.isTyping = false;
              agent.value.status = "idle";
              unsubscribe();

              // 触发回调
              if (config.onReceive) {
                config.onReceive(assistantMessage);
              }
            } else if (message.type === "error") {
              assistantMessage.content.text = `❌ ${message.payload?.message || "发送失败"}`;
              userMessage.status = "error";
              chatStore.isTyping = false;
              agent.value.status = "idle";
              unsubscribe();
              if (config.onError) {
                config.onError(new Error(message.payload?.message));
              }
            } else if (message.type === "agent_event") {
              const ev = message.payload?.event;
              const evType = message.payload?.type || ev?.type || ev?.EventType;
              if (ev && evType) {
                handleAgentEvent(evType, ev);
              }
            }
          });

          // 发送聊天消息
          const message = {
            type: "chat",
            payload: {
              template_id: config.agentId || "chat",
              input: content,
              model_config: config.modelConfig,
            },
          };

          console.log("📤 Sending WebSocket message:", message);
          ws.send(message);
          console.log("✅ Message sent to WebSocket");

          // WebSocket 是异步的，不需要等待这里
          // 状态会在消息回调中更新
        } else {
          // 回退到 HTTP API
          console.log("⚠️ WebSocket not connected, using HTTP API");
          const response = await client.agents.chatDirect(content, config.agentId || "chat");

          assistantMessage.content.text = response.text || response.data?.text || "无响应";
          assistantMessage.status = "sent";
          chatStore.isTyping = false;
          agent.value.status = "idle";
        }
      }
    } catch (error: any) {
      console.error("Send message error:", error);

      assistantMessage.content.text = `❌ 发送失败: ${error.message || "未知错误"}`;
      userMessage.status = "error";
      chatStore.isTyping = false;
      agent.value.status = "idle";

      if (config.onError) {
        config.onError(error);
      }
    }

    // 触发回调
    if (config.onSend) {
      config.onSend(userMessage);
    }
    if (config.onReceive && assistantMessage.content.text) {
      config.onReceive(assistantMessage);
    }
  };

  // 发送图片
  const sendImage = async (file: File) => {
    // TODO: 实现图片上传
    console.log("Send image:", file.name);

    // 创建图片消息占位
    const imageMessage: Message = {
      id: generateId("msg"),
      type: "image",
      role: "user",
      content: {
        url: URL.createObjectURL(file),
        alt: file.name,
      },
      createdAt: Date.now(),
      status: "pending",
    };
    chatStore.messages.push(imageMessage);

    // TODO: 上传到服务器并获取 URL
    // 当前只是本地预览
    imageMessage.status = "sent";
  };

  // 重试消息
  const retryMessage = async (message: Message) => {
    if (message.type === "text" && message.role === "user") {
      await sendMessage(message.content.text);
    }
  };

  // 删除消息
  const deleteMessage = (messageId: string) => {
    const index = chatStore.messages.findIndex((m) => m.id === messageId);
    if (index !== -1) {
      chatStore.messages.splice(index, 1);
    }
  };

  // 清空消息
  const clearMessages = () => {
    chatStore.clearMessages();
  };

  const handleAgentEvent = (type: string, ev: any, messageId?: string) => {
    // 获取当前活跃消息 ID (如果没有提供)
    const currentMessageId = messageId || chatStore.activeMessageId || "";

    // 1. 思维事件 → thinkingStore
    if (type === "think_chunk_start") {
      thinkingStore.startThinking(currentMessageId);
      chatStore.setActiveMessage(currentMessageId);
      return;
    }
    if (type === "think_chunk") {
      thinkingStore.handleThinkChunk(ev.delta || ev.content || "");
      return;
    }
    if (type === "think_chunk_end") {
      thinkingStore.endThinking();
      return;
    }

    // 2. 工具事件 → toolsStore + thinkingStore
    if (type === "tool:start" || type === "tool_call_start" || (type.startsWith("tool") && type.includes("start"))) {
      const call = ev.Call || ev.call || {};
      const toolCall = {
        id: call.id || call.ID || call.tool_call_id || generateId("tool"),
        name: call.name || "unknown",
        state: "executing" as const,
        progress: 0,
        arguments: call.arguments || {},
        cancelable: call.cancelable ?? false,
        pausable: call.pausable ?? false,
      };

      toolsStore.handleToolStart(toolCall);

      // 同时添加到思维步骤
      thinkingStore.addStep(currentMessageId, {
        type: "tool_call",
        tool: {
          name: toolCall.name,
          args: toolCall.arguments,
        },
        timestamp: Date.now(),
      });
      return;
    }

    if (type === "tool:progress" || type === "tool_call_progress" || (type.startsWith("tool") && type.includes("progress"))) {
      const call = ev.Call || ev.call || {};
      const id = call.id || call.ID || call.tool_call_id;
      if (id) {
        toolsStore.handleToolProgress(id, ev.progress ?? call.progress ?? 0, ev.message || "");
      }
      return;
    }

    if (type === "tool:end" || type === "tool_call_end" || (type.startsWith("tool") && type.includes("end"))) {
      const call = ev.Call || ev.call || {};
      const id = call.id || call.ID || call.tool_call_id;
      if (id) {
        const toolCall = {
          id,
          name: call.name || "unknown",
          state: (call.error || ev.error ? "failed" : "completed") as "failed" | "completed",
          progress: 1,
          arguments: call.arguments || {},
          result: call.result || ev.result,
          error: call.error || ev.error,
        };

        toolsStore.handleToolEnd(toolCall);

        // 添加工具结果到思维步骤
        thinkingStore.addStep(currentMessageId, {
          type: "tool_result",
          tool: {
            name: toolCall.name,
            args: toolCall.arguments,
          },
          result: toolCall.result,
          timestamp: Date.now(),
        });
      }
      return;
    }

    // 处理旧版本工具事件 (向后兼容)
    if (type.startsWith("tool")) {
      const call = ev.Call || ev.call || {};
      const id = call.id || call.ID || call.tool_call_id;
      if (!id) return;

      const toolCall = {
        id,
        name: call.name || "unknown",
        state: (call.state || ev.state || "executing") as any,
        progress: ev.progress ?? call.progress ?? 0,
        arguments: call.arguments || {},
        result: call.result || ev.result,
        error: ev.error || call.error,
        intermediate: ev.data || call.intermediate,
        cancelable: call.cancelable ?? false,
        pausable: call.pausable ?? false,
      };

      if (type.includes("start")) {
        toolsStore.handleToolStart(toolCall);
      } else if (type.includes("end")) {
        toolsStore.handleToolEnd(toolCall);
      } else {
        toolsStore.handleToolProgress(id, toolCall.progress, "");
      }
      return;
    }

    // 3. 审批事件 → approvalStore + thinkingStore
    if (type === "permission_required") {
      const call = ev.call || {};
      const requestId = ev.request_id || generateId("approval");

      approvalStore.addApprovalRequest({
        id: requestId,
        messageId: currentMessageId,
        toolName: call.name || "",
        args: call.arguments || {},
        reason: ev.reason || "",
        timestamp: Date.now(),
      });

      // 添加审批步骤到思维过程
      thinkingStore.addStep(currentMessageId, {
        type: "approval",
        tool: {
          name: call.name,
          args: call.arguments,
        },
        timestamp: Date.now(),
      });

      console.log("Permission required for tool:", call.name);
      return;
    }

    // 4. Todo 事件 → todosStore
    if (type === "todo_update" || type === "todos_updated") {
      todosStore.updateTodos(ev.todos || []);
      return;
    }

    // 5. 工作流事件 → workflowStore
    if (type === "workflow_start" || type === "workflow:start") {
      workflowStore.loadWorkflow({
        id: ev.workflow_id || generateId("workflow"),
        name: ev.name || ev.title || "工作流",
        title: ev.title,
        steps: ev.steps || [],
      });
      return;
    }

    if (type === "workflow_step_complete" || type === "workflow:step_complete") {
      workflowStore.completeStep(ev.step_id);
      return;
    }

    if (type === "workflow_step_update" || type === "workflow:step_update") {
      workflowStore.updateStep(ev.step_id, {
        status: ev.status,
        metadata: ev.metadata,
      });
      return;
    }

    // 6. 文本消息 → chatStore (使用 RAF 批量更新)
    if (type === "text_chunk" || type === "message_delta") {
      chatStore.handleTextChunk(currentMessageId, ev.delta || ev.content || ev.text || "");
      return;
    }

    // 7. AskUser 事件
    if (type === "ask_user") {
      pendingAskUser.value = {
        requestId: ev.request_id,
        questions: ev.questions || [],
      };
      // 添加 AskUser 消息到消息列表
      const askUserMsg: AskUserMessage = {
        id: generateId("msg"),
        type: "ask-user",
        role: "assistant",
        content: {
          request_id: ev.request_id,
          questions: ev.questions || [],
          answered: false,
        },
        createdAt: Date.now(),
      };
      chatStore.messages.push(askUserMsg);
      return;
    }

    // 8. 状态变更事件
    if (type === "state_changed") {
      const state = ev.state;
      if (state === "working" || state === "running") {
        agent.value.status = "thinking";
      } else if (state === "idle" || state === "ready" || state === "completed") {
        agent.value.status = "idle";
      } else if (state === "failed") {
        agent.value.status = "error";
      }
      return;
    }

    // 9. Token 使用统计
    if (type === "token_usage") {
      console.log("Token usage:", ev);
      return;
    }

    // 10. 错误事件
    if (type === "error") {
      console.error("Agent error:", ev.message, ev.detail);
      return;
    }
  };

  // 回答 AskUser 问题
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
    const msgIndex = chatStore.messages.findIndex((m) => m.type === "ask-user" && (m as AskUserMessage).content.request_id === requestId);
    if (msgIndex !== -1) {
      const msg = chatStore.messages[msgIndex] as AskUserMessage;
      msg.content.answered = true;
      msg.content.answers = answers;
    }

    pendingAskUser.value = null;
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

  // 初始化
  onMounted(() => {
    // 添加欢迎消息
    if (config.welcomeMessage && chatStore.messages.length === 0) {
      const welcomeText = typeof config.welcomeMessage === "string" ? config.welcomeMessage : config.welcomeMessage.type === "text" ? config.welcomeMessage.content.text : "👋 你好，我是 Nomad Copilot。";

      const welcomeMsg: TextMessage = {
        id: generateId("msg"),
        type: "text",
        role: "assistant",
        content: {
          text: welcomeText,
        },
        createdAt: Date.now(),
      };
      chatStore.messages.push(welcomeMsg);
    }
  });

  return {
    // 状态 (通过 computed 从 stores 获取)
    messages: computed(() => chatStore.messages),
    isTyping: computed(() => chatStore.isTyping),
    isConnected: wsConnected,
    connectionState,
    currentInput,
    agent,
    isThinking: computed(() => thinkingStore.isThinking),
    thinkingContent: computed(() => thinkingStore.currentThought),
    currentStep: computed(() => 0), // 暂时返回 0,未来可以从 workflowStore 获取
    todos: computed(() => todosStore.todos),
    toolRunsList: computed(() => Array.from(toolsStore.toolRuns.values())),
    pendingAskUser,

    // Stores (暴露给组件使用)
    chatStore,
    thinkingStore,
    toolsStore,
    todosStore,
    approvalStore,
    workflowStore,

    // 方法
    sendMessage,
    sendImage,
    retryMessage,
    deleteMessage,
    clearMessages,
    approveAction: (requestId: string) => {
      approvalStore.approve(requestId);
      config.onApproveAction?.(requestId);
    },
    rejectAction: (requestId: string, reason?: string) => {
      approvalStore.reject(requestId, reason);
      config.onRejectAction?.(requestId);
    },
    controlTool,
    answerQuestion,

    // 暴露事件处理方法供外部组件使用
    handleAgentEvent,
  };
}
