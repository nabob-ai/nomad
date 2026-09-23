/**
 * NomadClient 完整使用示例
 * 展示主客户端类的所有功能
 */

import { NomadClient } from "@nomad/client-js";

async function main() {
  console.log("=".repeat(70));
  console.log("NomadClient 完整功能演示");
  console.log("=".repeat(70));

  // ========================================================================
  // 1. 初始化客户端
  // ========================================================================
  console.log("\n🚀 1. 初始化 NomadClient 客户端");
  console.log("-".repeat(70));

  const client = new NomadClient({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
    timeout: 30000,
    retry: {
      maxRetries: 3,
      retryDelay: 1000,
    },
  });

  console.log("✅ 客户端已初始化");
  console.log("   Base URL:", "http://localhost:8080");

  // 健康检查
  try {
    const health = await client.healthCheck();
    console.log("💚 健康检查:", health.status);
    console.log("   组件状态:");
    Object.entries(health.components).forEach(([name, status]) => {
      const icon = status.status === "healthy" ? "✅" : "⚠️";
      console.log(`     ${icon} ${name}: ${status.status}`);
    });
  } catch (error: any) {
    console.log("⚠️  健康检查失败:", error.message);
  }

  // ========================================================================
  // 2. Memory 系统
  // ========================================================================
  console.log("\n🧠 2. Memory 系统");
  console.log("-".repeat(70));

  // Working Memory
  await client.memory.working.set(
    "user_preference",
    {
      theme: "dark",
      language: "zh-CN",
      notifications: true,
    },
    {
      scope: "resource", // 全局级别（跨会话）
      ttl: 3600,
    },
  );
  console.log("✅ Working Memory 已设置");

  const preference = await client.memory.working.get("user_preference");
  console.log("📝 获取 Working Memory:", preference?.value);

  // Semantic Memory
  await client.memory.semantic.store("NomadClient is a powerful framework for building AI agents", { source: "documentation", category: "introduction" });
  console.log("✅ Semantic Memory 已添加");

  const searchResults = await client.memory.semantic.search("What is NomadClient?", {
    limit: 3,
  });
  console.log(`🔍 搜索结果: ${searchResults.length} 条`);

  // ========================================================================
  // 3. Session 管理
  // ========================================================================
  console.log("\n💬 3. Session 管理");
  console.log("-".repeat(70));

  const session = await client.sessions.create({
    agentId: "assistant-agent",
    templateId: "chat-template",
    userId: "user-123",
    enableCheckpoints: true,
    checkpointInterval: 5,
  });
  console.log("✅ Session 已创建:", session.id);

  await client.sessions.addMessage(session.id, {
    role: "user",
    content: "Hello, how can you help me today?",
  });
  console.log("📨 用户消息已添加");

  await client.sessions.addMessage(session.id, {
    role: "assistant",
    content: "I can help you with various tasks. What do you need?",
  });
  console.log("🤖 助手响应已添加");

  const messages = await client.sessions.getMessages(session.id);
  console.log(`📋 消息列表: ${messages.items.length} 条消息`);

  // ========================================================================
  // 4. Workflow 编排
  // ========================================================================
  console.log("\n🔄 4. Workflow 编排");
  console.log("-".repeat(70));

  const workflow = await client.workflows.create({
    type: "sequential",
    name: "Document Processing Pipeline",
    description: "文档处理流水线",
    steps: [
      { agent: "reader", action: "read_document" },
      { agent: "analyzer", action: "analyze_content" },
      { agent: "summarizer", action: "generate_summary" },
    ],
  });
  console.log("✅ Workflow 已创建:", workflow.id);

  const run = await client.workflows.execute(workflow.id, {
    input: { documentUrl: "https://example.com/doc.pdf" },
  });
  console.log("▶️  Workflow 已启动:", run.id);
  console.log("   状态:", run.status);
  console.log("   进度:", run.progress, "%");

  // ========================================================================
  // 5. MCP 协议
  // ========================================================================
  console.log("\n🔌 5. MCP 协议");
  console.log("-".repeat(70));

  try {
    await client.mcp.addServer({
      serverId: "example-server",
      name: "Example MCP Server",
      endpoint: "http://localhost:8090/mcp",
      enabled: true,
    });
    console.log("✅ MCP Server 已添加");

    const servers = await client.mcp.listServers();
    console.log(`📋 MCP Servers: ${servers.length} 个`);
  } catch (error: any) {
    console.log("⚠️  MCP Server 操作失败:", error.message);
  }

  // ========================================================================
  // 6. Middleware 配置
  // ========================================================================
  console.log("\n🧅 6. Middleware 配置");
  console.log("-".repeat(70));

  const middlewares = await client.middleware.list();
  console.log(`📋 总共 ${middlewares.length} 个 Middlewares`);

  // 配置 Summarization
  await client.middleware.updateConfig("summarization", {
    threshold: 170000,
    keepMessages: 6,
    llmProvider: "anthropic",
    llmModel: "claude-sonnet-4",
  });
  console.log("✅ Summarization 已配置");

  // 配置 Cost Tracker
  await client.middleware.updateConfig("cost_tracker", {
    enabled: true,
    costModel: "token_based",
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
  console.log("✅ Cost Tracker 已配置");

  // ========================================================================
  // 7. Tool 执行
  // ========================================================================
  console.log("\n🔧 7. Tool 执行");
  console.log("-".repeat(70));

  const tools = await client.tools.list({ type: "builtin" });
  console.log(`📋 内置工具: ${tools.length} 个`);

  try {
    // 执行 Bash 工具
    const result = await client.tools.execute("bash", {
      command: 'echo "Hello from NomadClient"',
      timeout: 10,
    });
    console.log("✅ Bash 执行成功:");
    console.log("   耗时:", result.executionTime, "ms");
    console.log("   结果:", result.result);

    // 执行 HTTP 请求工具
    const httpResult = await client.tools.execute("http_request", {
      url: "https://api.github.com/zen",
      method: "GET",
    });
    console.log("✅ HTTP 请求成功:");
    console.log("   响应:", httpResult.result);
  } catch (error: any) {
    console.log("⚠️  工具执行失败:", error.message);
  }

  // ========================================================================
  // 8. Telemetry 监控
  // ========================================================================
  console.log("\n📊 8. Telemetry 监控");
  console.log("-".repeat(70));

  try {
    // 获取性能指标
    const performance = await client.getPerformanceMetrics({
      start: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
      end: new Date().toISOString(),
    });
    console.log("📈 性能指标（过去24小时）:");
    console.log("   总请求数:", performance.requests.total);
    console.log("   成功率:", ((performance.requests.successful / performance.requests.total) * 100).toFixed(1), "%");
    console.log("   平均延迟:", performance.requests.avgLatency.toFixed(2), "ms");
    console.log("   P95 延迟:", performance.requests.p95Latency.toFixed(2), "ms");
    console.log("   P99 延迟:", performance.requests.p99Latency.toFixed(2), "ms");

    if (performance.tokens) {
      console.log("   总 Tokens:", performance.tokens.total.toLocaleString());
    }
    if (performance.cost) {
      console.log("   总成本:", performance.cost.currency, performance.cost.total.toFixed(2));
    }

    // 获取使用统计
    const usage = await client.getUsageStatistics({
      start: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
      end: new Date().toISOString(),
    });
    console.log("\n📊 使用统计（过去7天）:");
    if (usage.sessions) {
      console.log("   Sessions: 总计", usage.sessions.total, "| 活跃", usage.sessions.active);
      console.log("   平均时长:", usage.sessions.avgDuration.toFixed(0), "秒");
    }
    if (usage.workflows) {
      console.log("   Workflows: 成功", usage.workflows.successful, "| 失败", usage.workflows.failed);
    }
    if (usage.tools) {
      console.log("   工具调用:", usage.tools.total, "次");
      if (usage.tools.topTools && usage.tools.topTools.length > 0) {
        console.log("   最常用工具:");
        usage.tools.topTools.slice(0, 3).forEach((tool, i) => {
          console.log(`     ${i + 1}. ${tool.toolName} - ${tool.callCount} 次`);
        });
      }
    }

    // 查询 Metrics
    const metrics = await client.telemetry.listMetrics();
    console.log(`\n📊 Metrics: ${metrics.length} 个`);
    metrics.slice(0, 5).forEach((metric, i) => {
      console.log(`   ${i + 1}. ${metric.name} (${metric.type}): ${metric.value} ${metric.unit || ""}`);
    });

    // 查询 Traces
    const traces = await client.telemetry.queryTraces({
      timeRange: {
        start: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
        end: new Date().toISOString(),
      },
      limit: 5,
    });
    console.log(`\n🔍 Traces（过去1小时）: ${traces.length} 个`);
    traces.forEach((trace, i) => {
      console.log(`   ${i + 1}. ${trace.operationName} - ${trace.duration}ms (${trace.status})`);
    });
  } catch (error: any) {
    console.log("⚠️  Telemetry 查询失败:", error.message);
  }

  // ========================================================================
  // 9. 导出数据
  // ========================================================================
  console.log("\n💾 9. 数据导出");
  console.log("-".repeat(70));

  try {
    // 导出 Session
    const sessionExport = await client.sessions.export(session.id, {
      format: "json",
      includeMetadata: true,
      includeStats: true,
    });
    console.log("✅ Session 已导出:");
    console.log("   格式:", sessionExport.format);
    console.log("   文件名:", sessionExport.suggestedFilename);

    // 导出 Metrics
    const metricsExport = await client.telemetry.exportMetrics("json", {
      start: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
      end: new Date().toISOString(),
    });
    console.log("✅ Metrics 已导出:");
    console.log("   格式:", metricsExport.format);
    console.log("   文件名:", metricsExport.suggestedFilename);
  } catch (error: any) {
    console.log("⚠️  导出失败:", error.message);
  }

  // ========================================================================
  // 10. 清理
  // ========================================================================
  console.log("\n🧹 10. 清理");
  console.log("-".repeat(70));

  await client.sessions.complete(session.id);
  console.log("✅ Session 已完成");

  await client.workflows.archiveBatch([workflow.id]);
  console.log("✅ Workflow 已归档");

  console.log("\n" + "=".repeat(70));
  console.log("✅ 演示完成！");
  console.log("=".repeat(70));

  console.log("\n📝 总结:");
  console.log("本示例展示了 NomadClient 的核心功能：");
  console.log("  1. ✅ 客户端初始化和健康检查");
  console.log("  2. ✅ Memory 系统（Working + Semantic）");
  console.log("  3. ✅ Session 管理和消息历史");
  console.log("  4. ✅ Workflow 编排和执行");
  console.log("  5. ✅ MCP 协议集成");
  console.log("  6. ✅ Middleware 配置");
  console.log("  7. ✅ Tool 执行");
  console.log("  8. ✅ Telemetry 监控");
  console.log("  9. ✅ 数据导出");
  console.log("  10. ✅ 资源清理");
}

// 运行示例
main().catch((error) => {
  console.error("❌ 错误:", error);
  process.exit(1);
});
