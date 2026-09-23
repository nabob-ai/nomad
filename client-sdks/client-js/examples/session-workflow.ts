/**
 * Session 和 Workflow 使用示例
 */

import { SessionResource, WorkflowResource, ParallelWorkflowDefinition, SequentialWorkflowDefinition, LoopWorkflowDefinition } from "@nomad/client-js";

async function main() {
  console.log("=".repeat(60));
  console.log("NomadClient Session + Workflow 演示");
  console.log("=".repeat(60));

  // 创建资源实例
  const session = new SessionResource({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  const workflow = new WorkflowResource({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // ========================================================================
  // 1. Session 管理演示
  // ========================================================================
  console.log("\n📝 1. Session 管理");
  console.log("-".repeat(60));

  // 创建 Session
  const newSession = await session.create({
    agentId: "agent-123",
    templateId: "assistant",
    userId: "user-456",
    metadata: {
      project: "demo",
      environment: "development",
    },
    enableCheckpoints: true,
    checkpointInterval: 5, // 每5条消息创建一个断点
  });
  console.log("✅ 创建 Session:", newSession.id);

  // 添加消息
  await session.addMessage(newSession.id, {
    role: "user",
    content: "Hello! Can you help me with something?",
  });
  console.log("✅ 添加用户消息");

  await session.addMessage(newSession.id, {
    role: "assistant",
    content: "Of course! I'd be happy to help. What do you need assistance with?",
  });
  console.log("✅ 添加助手消息");

  // 获取消息列表
  const messages = await session.getMessages(newSession.id, {
    page: 1,
    pageSize: 10,
    sort: "asc",
  });
  console.log(`📋 获取消息列表: ${messages.items.length} 条消息`);

  // ========================================================================
  // 2. Checkpoint 断点恢复演示
  // ========================================================================
  console.log("\n🔄 2. Checkpoint 断点恢复");
  console.log("-".repeat(60));

  // 创建手动 checkpoint
  const checkpoint = await session.createCheckpoint(newSession.id, "before-important-action");
  console.log("✅ 创建 Checkpoint:", checkpoint.id);

  // 获取所有 checkpoints
  const checkpoints = await session.getCheckpoints(newSession.id);
  console.log(`📊 总共 ${checkpoints.length} 个 Checkpoints:`);
  checkpoints.forEach((cp, index) => {
    console.log(`  ${index + 1}. [${cp.type}] Sequence: ${cp.sequence}, Time: ${cp.timestamp}`);
  });

  // 从 checkpoint 恢复
  if (checkpoints.length > 0) {
    console.log("\n🔄 从最新 Checkpoint 恢复...");
    const resumed = await session.resume(newSession.id, {
      checkpointId: checkpoints[checkpoints.length - 1].id,
      keepSubsequentMessages: false,
    });
    console.log("✅ Session 已恢复:", resumed.status);
  }

  // ========================================================================
  // 3. Session 统计演示
  // ========================================================================
  console.log("\n📊 3. Session 统计");
  console.log("-".repeat(60));

  const stats = await session.getStats(newSession.id);
  console.log("统计信息:");
  console.log(`  - 总消息数: ${stats.totalMessages}`);
  console.log(`  - 用户消息: ${stats.userMessages}`);
  console.log(`  - 助手消息: ${stats.assistantMessages}`);
  console.log(`  - 总 Tokens: ${stats.totalTokens}`);
  console.log(`  - 总成本: ${stats.totalCost} ${stats.currency}`);
  console.log(`  - 持续时间: ${stats.duration} 秒`);

  // ========================================================================
  // 4. Parallel Workflow 演示
  // ========================================================================
  console.log("\n🔀 4. Parallel Workflow（并行执行）");
  console.log("-".repeat(60));

  const parallelWorkflow: ParallelWorkflowDefinition = {
    type: "parallel",
    name: "Multi-Agent Research",
    description: "多个 Agent 并行研究不同主题",
    agents: [
      { id: "researcher-1", task: "Research AI trends in 2024" },
      { id: "researcher-2", task: "Research quantum computing developments" },
      { id: "researcher-3", task: "Research climate tech innovations" },
    ],
    maxConcurrency: 3,
    timeout: 300,
  };

  const parallelWf = await workflow.create(parallelWorkflow);
  console.log("✅ 创建 Parallel Workflow:", parallelWf.id);

  // 执行 Workflow
  const parallelRun = await workflow.execute(parallelWf.id, {
    input: "Please provide comprehensive research summaries",
    options: { async: false },
  });
  console.log("▶️  执行 Workflow, Run ID:", parallelRun.id);

  // 等待完成（模拟）
  try {
    const finalRun = await workflow.waitForCompletion(parallelWf.id, parallelRun.id, {
      pollInterval: 2000,
      timeout: 60000,
    });
    console.log("✅ Workflow 完成:", finalRun.status);
  } catch (error: any) {
    console.log("⚠️  等待超时或失败:", error.message);
  }

  // ========================================================================
  // 5. Sequential Workflow 演示
  // ========================================================================
  console.log("\n➡️  5. Sequential Workflow（顺序执行）");
  console.log("-".repeat(60));

  const sequentialWorkflow: SequentialWorkflowDefinition = {
    type: "sequential",
    name: "Document Processing Pipeline",
    description: "文档处理流水线",
    steps: [
      {
        agent: "reader",
        action: "read_document",
        params: { source: "https://example.com/doc.pdf" },
      },
      {
        agent: "analyzer",
        action: "analyze_content",
        params: { depth: "detailed" },
      },
      {
        agent: "summarizer",
        action: "generate_summary",
        params: { length: "medium" },
      },
      {
        agent: "translator",
        action: "translate",
        params: { targetLang: "zh-CN" },
        condition: 'previousStep.language === "en"',
      },
    ],
    continueOnError: false,
  };

  const sequentialWf = await workflow.create(sequentialWorkflow);
  console.log("✅ 创建 Sequential Workflow:", sequentialWf.id);

  const sequentialRun = await workflow.execute(sequentialWf.id, {
    input: { documentUrl: "https://example.com/doc.pdf" },
  });
  console.log("▶️  执行 Workflow, Run ID:", sequentialRun.id);

  // ========================================================================
  // 6. Loop Workflow 演示
  // ========================================================================
  console.log("\n🔁 6. Loop Workflow（循环执行）");
  console.log("-".repeat(60));

  const loopWorkflow: LoopWorkflowDefinition = {
    type: "loop",
    name: "Iterative Code Optimizer",
    description: "迭代优化代码直到达到质量标准",
    agent: "optimizer",
    condition: "result.quality < 0.95", // 质量 < 95% 则继续
    maxIterations: 10,
    initialInput: {
      code: "function add(a, b) { return a + b; }",
      targetQuality: 0.95,
    },
  };

  const loopWf = await workflow.create(loopWorkflow);
  console.log("✅ 创建 Loop Workflow:", loopWf.id);

  const loopRun = await workflow.execute(loopWf.id, {
    input: { code: "function example() { /* needs optimization */ }" },
  });
  console.log("▶️  执行 Workflow, Run ID:", loopRun.id);

  // ========================================================================
  // 7. Workflow 控制操作
  // ========================================================================
  console.log("\n⏯️  7. Workflow 控制操作");
  console.log("-".repeat(60));

  // 获取执行详情
  const runDetails = await workflow.getRunDetails(parallelWf.id, parallelRun.id);
  console.log("📊 执行详情:");
  console.log(`  - 状态: ${runDetails.status}`);
  console.log(`  - 进度: ${runDetails.progress}%`);
  console.log(`  - 步骤: ${runDetails.currentStep}/${runDetails.totalSteps}`);
  console.log(`  - 成功步骤: ${runDetails.stats.successfulSteps}`);
  console.log(`  - 总耗时: ${runDetails.stats.totalDuration}ms`);

  // 暂停执行（如果正在运行）
  if (loopRun.status === "running") {
    await workflow.suspend(loopWf.id, {
      runId: loopRun.id,
      reason: "User requested pause",
    });
    console.log("⏸️  已暂停 Workflow");

    // 恢复执行
    await workflow.resume(loopWf.id, {
      runId: loopRun.id,
    });
    console.log("▶️  已恢复 Workflow");
  }

  // ========================================================================
  // 8. Workflow 历史查询
  // ========================================================================
  console.log("\n📜 8. Workflow 历史查询");
  console.log("-".repeat(60));

  const runs = await workflow.getRuns(parallelWf.id, {
    page: 1,
    pageSize: 10,
  });
  console.log(`📋 执行历史: 共 ${runs.total} 次执行`);
  runs.items.forEach((run, index) => {
    console.log(`  ${index + 1}. [${run.status}] ${run.startedAt} - Progress: ${run.progress}%`);
  });

  // ========================================================================
  // 9. Session 导出
  // ========================================================================
  console.log("\n💾 9. Session 导出");
  console.log("-".repeat(60));

  const exported = await session.export(newSession.id, {
    format: "json",
    includeMetadata: true,
    includeStats: true,
  });
  console.log("✅ Session 已导出:");
  console.log(`  - 格式: ${exported.format}`);
  console.log(`  - 文件名: ${exported.suggestedFilename}`);
  console.log(`  - 导出时间: ${exported.exportedAt}`);

  // ========================================================================
  // 10. 清理
  // ========================================================================
  console.log("\n🧹 10. 清理");
  console.log("-".repeat(60));

  // 完成 Session
  await session.complete(newSession.id);
  console.log("✅ Session 已完成");

  // 归档 Workflows
  await workflow.archiveBatch([parallelWf.id, sequentialWf.id, loopWf.id]);
  console.log("✅ Workflows 已归档");

  console.log("\n" + "=".repeat(60));
  console.log("✅ 演示完成！");
  console.log("=".repeat(60));
}

// 运行示例
main().catch((error) => {
  console.error("❌ 错误:", error);
  process.exit(1);
});
