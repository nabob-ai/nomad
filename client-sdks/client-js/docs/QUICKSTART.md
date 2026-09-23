# nomad Client JS - 快速开始

> 5 分钟快速上手 nomad JavaScript 客户端

---

## 📦 安装

```bash
npm install @nomad/client-js
# 或
yarn add @nomad/client-js
# 或
pnpm add @nomad/client-js
```

---

## 🚀 快速开始

### 1. 初始化客户端

```typescript
import { nomad } from "@nomad/client-js";

const client = new nomad({
  baseUrl: "http://localhost:8080",
  apiKey: "your-api-key", // 可选
});
```

### 2. 健康检查

```typescript
const health = await client.healthCheck();
console.log("状态:", health.status); // healthy
```

---

## 💡 核心功能

### Memory 系统

#### Working Memory（工作记忆）

```typescript
// 设置
await client.memory.working.set(
  "user_preference",
  {
    theme: "dark",
    language: "zh-CN",
  },
  {
    scope: "resource", // 跨会话
    ttl: 3600, // 1小时过期
  },
);

// 获取
const item = await client.memory.working.get("user_preference");
console.log(item?.value); // { theme: 'dark', language: 'zh-CN' }
```

#### Semantic Memory（语义记忆）

```typescript
// 存储
await client.memory.semantic.store("nomad is a powerful framework for building AI agents", { category: "introduction" });

// 搜索
const results = await client.memory.semantic.search("What is nomad?", {
  limit: 5,
});

results.forEach((chunk) => {
  console.log(chunk.content, chunk.similarity);
});
```

---

### Session 管理

```typescript
// 创建 Session
const session = await client.sessions.create({
  agentId: "assistant-agent",
  templateId: "chat",
  userId: "user-123",
  enableCheckpoints: true, // 启用断点恢复
});

// 添加消息
await client.sessions.addMessage(session.id, {
  role: "user",
  content: "Hello!",
});

await client.sessions.addMessage(session.id, {
  role: "assistant",
  content: "Hi! How can I help you?",
});

// 获取消息历史
const messages = await client.sessions.getMessages(session.id);
console.log(`${messages.items.length} 条消息`);

// 完成 Session
await client.sessions.complete(session.id);
```

---

### Workflow 编排

#### Parallel Workflow（并行）

```typescript
// 创建并行工作流
const workflow = await client.workflows.create({
  type: "parallel",
  name: "Multi-Agent Research",
  agents: [
    { id: "researcher-1", task: "Research AI trends" },
    { id: "researcher-2", task: "Research quantum computing" },
  ],
  maxConcurrency: 2,
});

// 执行
const run = await client.workflows.execute(workflow.id, {
  input: "Start research",
});

// 等待完成
const result = await client.workflows.waitForCompletion(workflow.id, run.id);

console.log("状态:", result.status); // completed
```

#### Sequential Workflow（顺序）

```typescript
const workflow = await client.workflows.create({
  type: "sequential",
  name: "Document Pipeline",
  steps: [
    { agent: "reader", action: "read_document" },
    { agent: "analyzer", action: "analyze" },
    { agent: "summarizer", action: "summarize" },
  ],
});
```

---

### Tool 执行

```typescript
// 同步执行（快速工具）
const result = await client.tools.execute("bash", {
  command: 'echo "Hello nomad"',
  timeout: 10,
});

console.log(result.result); // Hello nomad
console.log("耗时:", result.executionTime, "ms");

// 异步执行（长时运行工具）
const task = await client.tools.executeAsync("web_scraper", {
  url: "https://example.com",
});

// 等待完成
const taskResult = await client.tools.waitForTask(task.taskId);
console.log("抓取结果:", taskResult.result);
```

---

### MCP 协议

```typescript
// 添加 MCP Server
await client.mcp.addServer({
  serverId: "my-server",
  name: "My MCP Server",
  endpoint: "http://localhost:8090/mcp",
});

// 调用远程工具
const result = await client.mcp.call("my-server", "calculator", {
  operation: "add",
  numbers: [1, 2, 3],
});

console.log("结果:", result.result); // 6
```

---

### Middleware 配置

```typescript
// 配置上下文总结
await client.middleware.updateConfig("summarization", {
  threshold: 170000, // 170K tokens 触发
  keepMessages: 6,
});

// 配置成本追踪
await client.middleware.updateConfig("cost_tracker", {
  enabled: true,
  pricing: {
    promptTokenPrice: 0.003,
    completionTokenPrice: 0.015,
    currency: "USD",
  },
  budget: {
    daily: 100,
    monthly: 2000,
  },
});

// 列出所有 Middleware
const middlewares = await client.middleware.list();
middlewares.forEach((mw) => {
  console.log(`${mw.displayName}: ${mw.enabled ? "ON" : "OFF"}`);
});
```

---

### Telemetry 监控

```typescript
// 获取性能指标
const perf = await client.getPerformanceMetrics({
  start: "2024-01-01T00:00:00Z",
  end: "2024-01-02T00:00:00Z",
});

console.log("总请求数:", perf.requests.total);
console.log("平均延迟:", perf.requests.avgLatency, "ms");
console.log("P95 延迟:", perf.requests.p95Latency, "ms");

// 获取使用统计
const usage = await client.getUsageStatistics({
  start: "2024-01-01T00:00:00Z",
  end: "2024-01-08T00:00:00Z",
});

console.log("总 Sessions:", usage.sessions?.total);
console.log("最常用工具:", usage.tools?.topTools);

// 查询 Metrics
const metrics = await client.telemetry.listMetrics();
metrics.forEach((m) => {
  console.log(`${m.name}: ${m.value} ${m.unit || ""}`);
});
```

---

## 🎯 完整示例

### 示例 1：智能对话 Agent

```typescript
import { nomad } from "@nomad/client-js";

async function chatExample() {
  const client = new nomad({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // 1. 创建 Session
  const session = await client.sessions.create({
    agentId: "chat-agent",
    templateId: "assistant",
    userId: "user-123",
    enableCheckpoints: true,
  });

  // 2. 存储用户偏好到 Working Memory
  await client.memory.working.set(
    "user_context",
    {
      name: "Alice",
      preferences: { language: "zh-CN" },
    },
    { scope: "resource" },
  );

  // 3. 对话
  await client.sessions.addMessage(session.id, {
    role: "user",
    content: "帮我总结一下 AI 的最新趋势",
  });

  // 4. 存储重要信息到 Semantic Memory
  await client.memory.semantic.store("User is interested in AI trends", {
    userId: "user-123",
    topic: "AI",
  });

  // 5. 获取统计
  const stats = await client.sessions.getStats(session.id);
  console.log("总 Tokens:", stats.totalTokens);
  console.log("总成本:", stats.totalCost, stats.currency);

  // 6. 完成
  await client.sessions.complete(session.id);
}

chatExample().catch(console.error);
```

### 示例 2：并行研究工作流

```typescript
import { nomad } from "@nomad/client-js";

async function researchWorkflow() {
  const client = new nomad({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // 1. 创建并行工作流
  const workflow = await client.workflows.create({
    type: "parallel",
    name: "Tech Research",
    agents: [
      { id: "ai-researcher", task: "Research AI developments" },
      { id: "quantum-researcher", task: "Research quantum computing" },
      { id: "climate-researcher", task: "Research climate tech" },
    ],
    maxConcurrency: 3,
    timeout: 300,
  });

  // 2. 执行
  const run = await client.workflows.execute(workflow.id, {
    input: "Research latest developments in 2024",
  });

  // 3. 监控进度
  const interval = setInterval(async () => {
    const status = await client.workflows.getRunStatus(workflow.id, run.id);
    console.log(`进度: ${status.progress}%`);

    if (status.status !== "running") {
      clearInterval(interval);
    }
  }, 2000);

  // 4. 等待完成
  const result = await client.workflows.waitForCompletion(workflow.id, run.id, {
    timeout: 300000,
  });

  console.log("工作流完成:", result.status);
  console.log("结果:", result.output);

  // 5. 获取详细结果
  const details = await client.workflows.getRunDetails(workflow.id, run.id);
  details.steps.forEach((step, i) => {
    console.log(`步骤 ${i + 1}: ${step.status} (${step.duration}ms)`);
  });
}

researchWorkflow().catch(console.error);
```

### 示例 3：使用 MCP 和 Tools

```typescript
import { nomad } from "@nomad/client-js";

async function mcpToolsExample() {
  const client = new nomad({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // 1. 添加 MCP Server
  await client.mcp.addServer({
    serverId: "file-server",
    name: "File Operations Server",
    endpoint: "http://localhost:8090/mcp",
  });

  // 2. 调用 MCP 工具
  const fileList = await client.mcp.call("file-server", "list_files", {
    directory: "/tmp",
  });

  // 3. 使用本地工具
  const httpResult = await client.tools.execute("http_request", {
    url: "https://api.github.com/zen",
    method: "GET",
  });

  console.log("GitHub Zen:", httpResult.result);

  // 4. 异步执行长时运行工具
  const task = await client.tools.executeAsync("web_scraper", {
    url: "https://example.com",
    selectors: ["h1", "p"],
  });

  const scrapeResult = await client.tools.waitForTask(task.taskId);
  console.log("抓取完成:", scrapeResult.result);

  // 5. 获取工具统计
  const stats = await client.tools.getStats("http_request");
  console.log("HTTP 请求调用次数:", stats.totalCalls);
  console.log("平均耗时:", stats.avgExecutionTime, "ms");
}

mcpToolsExample().catch(console.error);
```

---

## 📚 下一步

- 📖 查看完整 [API 文档](./API.md)
- 💡 浏览 [示例代码](../examples/)
- 🔧 了解[架构设计](../../ARCHITECTURE.md)
- 📊 查看[进度报告](../PROGRESS.md)

---

## 🛠️ 调试

### 启用日志

```typescript
// 配置日志 Middleware
await client.middleware.updateConfig("logging", {
  level: "debug",
  logRequests: true,
  logResponses: true,
  format: "json",
  outputs: ["console"],
});
```

### 错误处理

```typescript
try {
  await client.tools.execute("bash", {
    command: "invalid-command",
  });
} catch (error: any) {
  console.error("错误:", error.message);
  console.error("状态码:", error.statusCode);

  if (error.statusCode === 429) {
    console.error("速率限制，请稍后重试");
  }
}
```

---

## ❓ 常见问题

### Q: 如何设置 API Key？

A: 在初始化客户端时传入：

```typescript
const client = new nomad({
  baseUrl: "http://localhost:8080",
  apiKey: process.env.NOMAD_API_KEY,
});
```

### Q: 如何处理超时？

A: 可以在配置中设置全局超时，或在单个请求中设置：

```typescript
// 全局超时
const client = new nomad({
  baseUrl: "http://localhost:8080",
  timeout: 60000, // 60秒
});

// 单个请求超时
await client.tools.execute("bash", {
  command: "long-running-command",
  timeout: 120, // 120秒
});
```

### Q: 如何重试失败的请求？

A: 配置重试选项：

```typescript
const client = new nomad({
  baseUrl: "http://localhost:8080",
  retry: {
    maxRetries: 3,
    retryDelay: 1000,
  },
});
```

### Q: 支持哪些 Node.js 版本？

A: Node.js 16+ 或最新的 LTS 版本。

---