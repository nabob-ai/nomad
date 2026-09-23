/**
 * Agent 使用示例
 * 展示 Agent 的创建、管理、对话等功能
 */

import { NomadClient } from "@nomad/client-js";

async function main() {
  console.log("=".repeat(70));
  console.log("NomadClient Agent 功能演示");
  console.log("=".repeat(70));

  // 初始化客户端
  const client = new NomadClient({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // ========================================================================
  // 1. Agent 模板
  // ========================================================================
  console.log("\n📋 1. Agent 模板");
  console.log("-".repeat(70));

  // 列出所有模板
  const templates = await client.agents.listTemplates();
  console.log(`✅ 可用模板: ${templates.length} 个`);
  templates.forEach((template, i) => {
    const icon = template.builtin ? "🔧" : "✨";
    console.log(`   ${icon} ${i + 1}. ${template.name} (${template.type})`);
    console.log(`      ${template.description}`);
  });

  // ========================================================================
  // 2. 创建 Agent
  // ========================================================================
  console.log("\n🤖 2. 创建 Agent");
  console.log("-".repeat(70));

  // 从模板创建
  const assistant = await client.agents.createFromTemplate("assistant", {
    name: "My Assistant",
    description: "我的智能助手",
    llmProvider: "anthropic",
    llmModel: "claude-sonnet-4",
    llmParams: {
      temperature: 0.7,
      maxTokens: 4096,
    },
  });
  console.log("✅ 从模板创建 Agent:", assistant.id);
  console.log("   名称:", assistant.name);
  console.log("   状态:", assistant.status);
  console.log("   LLM:", `${assistant.llmProvider}/${assistant.llmModel}`);

  // 直接创建
  const researcher = await client.agents.create({
    name: "Research Agent",
    description: "AI 研究专家",
    templateId: "researcher",
    llmProvider: "openai",
    llmModel: "gpt-4-turbo",
    systemPrompt: "You are an expert AI researcher. Provide detailed, accurate information.",
    tools: ["http_request", "web_scraper"],
    middlewares: ["summarization", "cost_tracker"],
  });
  console.log("✅ 直接创建 Agent:", researcher.id);

  // ========================================================================
  // 3. Agent 管理
  // ========================================================================
  console.log("\n📂 3. Agent 管理");
  console.log("-".repeat(70));

  // 列出所有 Agents
  const agents = await client.agents.list({
    status: "active",
    page: 1,
    pageSize: 10,
    sortBy: "createdAt",
    sortOrder: "desc",
  });
  console.log(`📋 总共 ${agents.total} 个 Agents (显示 ${agents.items.length} 个):`);
  agents.items.forEach((agent, i) => {
    console.log(`   ${i + 1}. ${agent.name} (${agent.id})`);
    console.log(`      状态: ${agent.status} | LLM: ${agent.llmProvider}/${agent.llmModel}`);
  });

  // 获取 Agent 详情
  const agentDetail = await client.agents.get(assistant.id);
  console.log("\n🔍 Agent 详情:");
  console.log("   名称:", agentDetail.name);
  console.log("   模板:", agentDetail.templateId);
  console.log("   工具:", agentDetail.tools?.join(", ") || "无");
  console.log("   中间件:", agentDetail.middlewares?.join(", ") || "无");
  console.log("   版本:", agentDetail.version);

  // 更新 Agent
  await client.agents.update(assistant.id, {
    description: "我的智能助手 - 已更新",
    llmParams: {
      temperature: 0.8,
    },
    tools: ["bash", "http_request", "file_read"],
  });
  console.log("✅ Agent 已更新");

  // ========================================================================
  // 4. Agent 对话（同步）
  // ========================================================================
  console.log("\n💬 4. Agent 对话（同步）");
  console.log("-".repeat(70));

  const chatResponse = await client.agents.chat(assistant.id, {
    input: "Hello! Can you help me understand what NomadClient is?",
    userId: "user-123",
  });

  console.log("🤖 Agent 响应:");
  console.log(`   Session ID: ${chatResponse.sessionId}`);
  console.log(`   响应: ${chatResponse.response}`);

  if (chatResponse.usage) {
    console.log("   Token 使用:");
    console.log(`     Prompt: ${chatResponse.usage.promptTokens}`);
    console.log(`     Completion: ${chatResponse.usage.completionTokens}`);
    console.log(`     Total: ${chatResponse.usage.totalTokens}`);
  }

  if (chatResponse.cost) {
    console.log(`   成本: ${chatResponse.cost.currency} ${chatResponse.cost.amount.toFixed(4)}`);
  }

  console.log(`   执行时间: ${chatResponse.executionTime}ms`);

  // 继续对话（复用 Session）
  const followUp = await client.agents.chat(assistant.id, {
    input: "Can you give me an example use case?",
    sessionId: chatResponse.sessionId,
    userId: "user-123",
  });
  console.log("\n💬 继续对话:");
  console.log(`   响应: ${followUp.response}`);

  // ========================================================================
  // 5. Agent 对话（流式）
  // ========================================================================
  console.log("\n🌊 5. Agent 对话（流式）");
  console.log("-".repeat(70));

  console.log("🤖 开始流式对话...");
  let streamResponse = "";

  try {
    for await (const event of client.agents.chatStream(assistant.id, {
      input: "Tell me about the benefits of using NomadClient",
      userId: "user-123",
    })) {
      switch (event.type) {
        case "start":
          console.log(`   Session: ${event.sessionId}`);
          break;
        case "token":
          streamResponse += event.token;
          process.stdout.write(event.token);
          break;
        case "tool_call":
          console.log(`\n   🔧 工具调用: ${event.toolCall.name}`);
          break;
        case "end":
          console.log("\n   ✅ 完成");
          console.log(`   总耗时: ${event.response.executionTime}ms`);
          break;
        case "error":
          console.log(`\n   ❌ 错误: ${event.error}`);
          break;
      }
    }
  } catch (error: any) {
    console.log("\n⚠️  流式对话失败:", error.message);
  }

  // ========================================================================
  // 6. Agent 统计
  // ========================================================================
  console.log("\n📊 6. Agent 统计");
  console.log("-".repeat(70));

  const stats = await client.agents.getStats(assistant.id, {
    start: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
    end: new Date().toISOString(),
  });

  console.log("Agent 统计（过去7天）:");
  console.log(`   总请求数: ${stats.totalRequests}`);
  console.log(`   成功率: ${((stats.successfulRequests / stats.totalRequests) * 100).toFixed(1)}%`);
  console.log(`   平均响应时间: ${stats.avgResponseTime.toFixed(2)}ms`);
  console.log(`   Token 使用: ${stats.tokenUsage.totalTokens.toLocaleString()}`);
  console.log(`   总成本: ${stats.cost.currency} ${stats.cost.total.toFixed(4)}`);

  if (stats.toolCalls) {
    console.log(`   工具调用: ${stats.toolCalls.total} 次`);
    const topTools = Object.entries(stats.toolCalls.byTool)
      .sort(([, a], [, b]) => b - a)
      .slice(0, 3);
    if (topTools.length > 0) {
      console.log("   最常用工具:");
      topTools.forEach(([tool, count]) => {
        console.log(`     - ${tool}: ${count} 次`);
      });
    }
  }

  // 汇总统计
  const aggregated = await client.agents.getAggregatedStats({
    start: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    end: new Date().toISOString(),
  });
  console.log("\n📈 所有 Agents 汇总（过去24小时）:");
  console.log(`   总 Agents: ${aggregated.totalAgents} | 活跃: ${aggregated.activeAgents}`);
  console.log(`   总请求数: ${aggregated.totalRequests.toLocaleString()}`);
  console.log(`   总 Tokens: ${aggregated.totalTokens.toLocaleString()}`);
  console.log(`   总成本: ${aggregated.currency} ${aggregated.totalCost.toFixed(2)}`);

  // ========================================================================
  // 7. Agent 克隆
  // ========================================================================
  console.log("\n📋 7. Agent 克隆");
  console.log("-".repeat(70));

  const cloned = await client.agents.clone(assistant.id, "My Assistant (Clone)");
  console.log("✅ Agent 已克隆:", cloned.id);
  console.log("   原始 Agent:", assistant.id);
  console.log("   克隆 Agent:", cloned.id);
  console.log("   名称:", cloned.name);

  // ========================================================================
  // 8. Agent 状态管理
  // ========================================================================
  console.log("\n⚙️  8. Agent 状态管理");
  console.log("-".repeat(70));

  // 禁用 Agent
  await client.agents.disable(cloned.id);
  console.log("⏸️  Agent 已禁用:", cloned.id);

  // 激活 Agent
  await client.agents.activate(cloned.id);
  console.log("▶️  Agent 已激活:", cloned.id);

  // 归档 Agent
  await client.agents.archive(cloned.id);
  console.log("📦 Agent 已归档:", cloned.id);

  // ========================================================================
  // 9. Agent 验证
  // ========================================================================
  console.log("\n✅ 9. Agent 验证");
  console.log("-".repeat(70));

  const validation = await client.agents.validate({
    name: "Test Agent",
    templateId: "assistant",
    llmProvider: "openai",
    llmModel: "gpt-4",
    llmParams: {
      temperature: 0.7,
    },
  });

  console.log("验证结果:");
  console.log(`   有效: ${validation.valid ? "✅" : "❌"}`);
  if (validation.errors && validation.errors.length > 0) {
    console.log("   错误:");
    validation.errors.forEach((err) => console.log(`     - ${err}`));
  }
  if (validation.warnings && validation.warnings.length > 0) {
    console.log("   警告:");
    validation.warnings.forEach((warn) => console.log(`     - ${warn}`));
  }

  // ========================================================================
  // 10. 批量操作
  // ========================================================================
  console.log("\n🗂️  10. 批量操作");
  console.log("-".repeat(70));

  // 批量归档
  await client.agents.archiveBatch([researcher.id]);
  console.log("✅ 批量归档完成");

  // 批量激活
  await client.agents.activateBatch([researcher.id]);
  console.log("✅ 批量激活完成");

  // ========================================================================
  // 11. 清理
  // ========================================================================
  console.log("\n🧹 11. 清理");
  console.log("-".repeat(70));

  await client.agents.delete(assistant.id);
  console.log("✅ Agent 已删除:", assistant.id);

  await client.agents.deleteBatch([researcher.id, cloned.id]);
  console.log("✅ 批量删除完成");

  console.log("\n" + "=".repeat(70));
  console.log("✅ 演示完成！");
  console.log("=".repeat(70));

  console.log("\n📝 总结:");
  console.log("本示例展示了 Agent 的完整功能：");
  console.log("  1. ✅ Agent 模板浏览和使用");
  console.log("  2. ✅ Agent 创建（从模板/直接创建）");
  console.log("  3. ✅ Agent 管理（列表、详情、更新）");
  console.log("  4. ✅ Agent 对话（同步）");
  console.log("  5. ✅ Agent 对话（流式）");
  console.log("  6. ✅ Agent 统计和汇总");
  console.log("  7. ✅ Agent 克隆");
  console.log("  8. ✅ Agent 状态管理");
  console.log("  9. ✅ Agent 验证");
  console.log("  10. ✅ 批量操作");
  console.log("  11. ✅ 资源清理");
}

// 运行示例
main().catch((error) => {
  console.error("❌ 错误:", error);
  process.exit(1);
});
