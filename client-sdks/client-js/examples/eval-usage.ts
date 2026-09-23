/**
 * Eval 使用示例
 * 展示 Agent 评估、测试和基准测试功能
 */

import { NomadClient } from "@nomad/client-js";

async function main() {
  console.log("=".repeat(70));
  console.log("NomadClient Eval 功能演示");
  console.log("=".repeat(70));

  // 初始化客户端
  const client = new NomadClient({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  // ========================================================================
  // 1. 创建测试用例集
  // ========================================================================
  console.log("\n📝 1. 创建测试用例集");
  console.log("-".repeat(70));

  const testCaseSet = await client.evals.createTestCaseSet(
    "Q&A Test Cases",
    [
      {
        id: "test-1",
        name: "Basic Greeting",
        input: "Hello, how are you?",
        expectedOutput: "I am doing well, thank you for asking!",
        tags: ["greeting", "simple"],
      },
      {
        id: "test-2",
        name: "Technical Question",
        input: "What is the difference between HTTP and HTTPS?",
        expectedOutput: "HTTPS is the secure version of HTTP. It uses SSL/TLS encryption to protect data in transit.",
        tags: ["technical", "security"],
      },
      {
        id: "test-3",
        name: "Complex Query",
        input: "Explain how machine learning models are trained",
        expectedOutput: "Machine learning models are trained by feeding them data and adjusting their parameters to minimize prediction errors.",
        tags: ["ml", "complex"],
      },
    ],
    "A collection of Q&A test cases for agent evaluation",
  );

  console.log("✅ 测试用例集已创建:", testCaseSet.id);
  console.log(`   名称: ${testCaseSet.name}`);
  console.log(`   测试用例数: ${testCaseSet.testCases.length}`);

  // ========================================================================
  // 2. 快速单次评估
  // ========================================================================
  console.log("\n⚡ 2. 快速单次评估");
  console.log("-".repeat(70));

  // 创建一个测试 Agent
  const agent = await client.agents.createFromTemplate("assistant", {
    name: "Test Assistant",
    llmProvider: "openai",
    llmModel: "gpt-4",
  });
  console.log("✅ 测试 Agent 已创建:", agent.id);

  // 执行快速评估
  const quickResult = await client.evals.quickEval(agent.id, "What is AI?", "Artificial Intelligence (AI) refers to computer systems that can perform tasks requiring human intelligence.", [
    { type: "semantic_similarity", weight: 0.5, params: { threshold: 0.7 } },
    { type: "keyword_coverage", weight: 0.3 },
    { type: "coherence", weight: 0.2 },
  ]);

  console.log("\n📊 快速评估结果:");
  console.log(`   状态: ${quickResult.status}`);
  console.log(`   总分: ${quickResult.summary.avgScore.toFixed(2)}`);
  console.log(`   通过率: ${(quickResult.summary.passRate * 100).toFixed(1)}%`);

  const testResult = quickResult.testCaseResults[0];
  console.log("\n   详细结果:");
  console.log(`   - Agent 输出: ${testResult.output.substring(0, 100)}...`);
  console.log(`   - 评分: ${testResult.overallScore.toFixed(2)}`);
  console.log(`   - 通过: ${testResult.passed ? "✅" : "❌"}`);
  console.log(`   - 执行时间: ${testResult.executionTime}ms`);

  // ========================================================================
  // 3. 批量评估
  // ========================================================================
  console.log("\n📦 3. 批量评估");
  console.log("-".repeat(70));

  console.log("开始批量评估...");
  const batchResult = await client.evals.batchEval(
    agent.id,
    testCaseSet.testCases,
    [
      { type: "semantic_similarity", weight: 0.4 },
      {
        type: "llm_judge",
        weight: 0.4,
        params: {
          provider: "openai",
          model: "gpt-4",
          criteria: ["accuracy", "completeness", "clarity"],
        },
      },
      { type: "coherence", weight: 0.2 },
    ],
    2, // 并发数
  );

  console.log("\n📊 批量评估结果:");
  console.log(`   总测试用例: ${batchResult.summary.totalTestCases}`);
  console.log(`   通过: ${batchResult.summary.passed} | 失败: ${batchResult.summary.failed}`);
  console.log(`   通过率: ${(batchResult.summary.passRate * 100).toFixed(1)}%`);
  console.log(`   平均分数: ${batchResult.summary.avgScore.toFixed(2)}`);
  console.log(`   平均执行时间: ${batchResult.summary.avgExecutionTime.toFixed(0)}ms`);

  if (batchResult.summary.totalTokenUsage) {
    console.log(`   总 Tokens: ${batchResult.summary.totalTokenUsage.totalTokens.toLocaleString()}`);
  }

  if (batchResult.summary.totalCost) {
    console.log(`   总成本: ${batchResult.summary.totalCost.currency} ${batchResult.summary.totalCost.amount.toFixed(4)}`);
  }

  console.log("\n   各 Scorer 平均分:");
  Object.entries(batchResult.summary.avgScoresByScorer).forEach(([scorer, score]) => {
    console.log(`     ${scorer}: ${score.toFixed(2)}`);
  });

  console.log("\n   各测试用例结果:");
  batchResult.testCaseResults.forEach((result, i) => {
    const icon = result.passed ? "✅" : "❌";
    console.log(`     ${icon} ${result.testCaseName}: ${result.overallScore.toFixed(2)}`);
  });

  // ========================================================================
  // 4. Benchmark 多个 Agents
  // ========================================================================
  console.log("\n🏆 4. Benchmark 多个 Agents");
  console.log("-".repeat(70));

  // 创建第二个 Agent
  const agent2 = await client.agents.createFromTemplate("assistant", {
    name: "Test Assistant 2",
    llmProvider: "anthropic",
    llmModel: "claude-sonnet-4",
  });
  console.log("✅ 第二个测试 Agent 已创建:", agent2.id);

  // 执行 Benchmark
  console.log("\n开始 Benchmark...");
  const benchmark = await client.evals.createBenchmark({
    name: "Agent Comparison Benchmark",
    description: "比较两个 Agents 的性能",
    agentIds: [agent.id, agent2.id],
    testCaseSetId: testCaseSet.id,
    scorers: [
      { type: "semantic_similarity", weight: 0.5 },
      { type: "keyword_coverage", weight: 0.3 },
      { type: "coherence", weight: 0.2 },
    ],
    concurrency: 2,
  });

  // 等待 Benchmark 完成
  const benchmarkResult = await client.evals.waitForBenchmarkCompletion(benchmark.id);

  console.log("\n📊 Benchmark 结果:");
  console.log(`   状态: ${benchmarkResult.status}`);
  console.log("\n   排行榜:");
  benchmarkResult.leaderboard.forEach((entry) => {
    const medal = entry.rank === 1 ? "🥇" : entry.rank === 2 ? "🥈" : "🥉";
    console.log(`     ${medal} #${entry.rank} ${entry.agentName}`);
    console.log(`        平均分数: ${entry.avgScore.toFixed(2)}`);
    console.log(`        通过率: ${(entry.passRate * 100).toFixed(1)}%`);
    console.log(`        平均响应时间: ${entry.avgExecutionTime.toFixed(0)}ms`);
  });

  // ========================================================================
  // 5. A/B 测试
  // ========================================================================
  console.log("\n🔬 5. A/B 测试");
  console.log("-".repeat(70));

  console.log("开始 A/B 测试...");
  const abTestResult = await client.evals.compareAgents(agent.id, agent2.id, testCaseSet.id, [
    { type: "semantic_similarity", weight: 0.5 },
    { type: "keyword_coverage", weight: 0.5 },
  ]);

  console.log("\n📊 A/B 测试结果:");
  console.log(`   状态: ${abTestResult.status}`);

  const stats = abTestResult.statisticalAnalysis;
  console.log("\n   统计分析:");
  console.log(`     Agent A 平均分: ${stats.agentAAvgScore.toFixed(2)}`);
  console.log(`     Agent B 平均分: ${stats.agentBAvgScore.toFixed(2)}`);
  console.log(`     差异: ${stats.difference > 0 ? "+" : ""}${stats.difference.toFixed(2)} (${stats.differencePercent > 0 ? "+" : ""}${stats.differencePercent.toFixed(1)}%)`);
  console.log(`     p-value: ${stats.pValue.toFixed(4)}`);
  console.log(`     显著性: ${stats.isSignificant ? "✅ 显著" : "❌ 不显著"}`);

  if (stats.winner) {
    const winnerIcon = stats.winner === "A" ? "🏆" : stats.winner === "B" ? "🏆" : "🤝";
    console.log(`     胜者: ${winnerIcon} Agent ${stats.winner}`);
  }

  // ========================================================================
  // 6. 生成报告
  // ========================================================================
  console.log("\n📄 6. 生成报告");
  console.log("-".repeat(70));

  // 生成 HTML 报告
  const htmlReport = await client.evals.generateReport({
    evalId: batchResult.evalId,
    format: "html",
    includeDetails: true,
    includeVisualization: true,
  });
  console.log("✅ HTML 报告已生成");
  console.log(`   格式: ${htmlReport.format}`);
  console.log(`   生成时间: ${htmlReport.generatedAt}`);
  console.log(`   内容长度: ${htmlReport.content.length} 字符`);

  // 生成 Markdown 报告
  const mdReport = await client.evals.generateReport({
    evalId: batchResult.evalId,
    format: "markdown",
    includeDetails: false,
  });
  console.log("\n✅ Markdown 报告已生成");
  console.log("   预览:");
  console.log(mdReport.content.substring(0, 200) + "...");

  // 导出 JSON
  const jsonExport = await client.evals.exportResult(batchResult.evalId, "json");
  console.log("\n✅ JSON 导出已完成");
  console.log(`   大小: ${jsonExport.length} 字符`);

  // ========================================================================
  // 7. Eval 管理
  // ========================================================================
  console.log("\n📂 7. Eval 管理");
  console.log("-".repeat(70));

  // 列出所有 Evals
  const evals = await client.evals.list({
    agentId: agent.id,
    status: "completed",
    page: 1,
    pageSize: 10,
    sortBy: "createdAt",
    sortOrder: "desc",
  });

  console.log(`📋 找到 ${evals.total} 个 Evals (显示 ${evals.items.length} 个):`);
  evals.items.forEach((evalInfo, i) => {
    console.log(`   ${i + 1}. ${evalInfo.name} (${evalInfo.type})`);
    console.log(`      状态: ${evalInfo.status} | 进度: ${evalInfo.progress}%`);
    console.log(`      测试用例: ${evalInfo.completedTestCases}/${evalInfo.totalTestCases}`);
  });

  // 获取 Eval 详情
  const evalDetail = await client.evals.get(batchResult.evalId);
  console.log("\n🔍 Eval 详情:");
  console.log(`   ID: ${evalDetail.id}`);
  console.log(`   名称: ${evalDetail.name}`);
  console.log(`   类型: ${evalDetail.type}`);
  console.log(`   状态: ${evalDetail.status}`);
  console.log(`   Agent: ${evalDetail.agentId}`);

  // ========================================================================
  // 8. 测试用例集管理
  // ========================================================================
  console.log("\n📚 8. 测试用例集管理");
  console.log("-".repeat(70));

  // 列出所有测试用例集
  const testCaseSets = await client.evals.listTestCaseSets();
  console.log(`📋 找到 ${testCaseSets.length} 个测试用例集:`);
  testCaseSets.forEach((set, i) => {
    console.log(`   ${i + 1}. ${set.name} (${set.testCases.length} 个用例)`);
  });

  // 更新测试用例集
  await client.evals.updateTestCaseSet(testCaseSet.id, {
    description: "Updated: " + testCaseSet.description,
  });
  console.log("\n✅ 测试用例集已更新");

  // ========================================================================
  // 9. 清理
  // ========================================================================
  console.log("\n🧹 9. 清理");
  console.log("-".repeat(70));

  // 删除 Evals
  await client.evals.delete(quickResult.evalId);
  await client.evals.delete(batchResult.evalId);
  console.log("✅ Evals 已删除");

  // 删除测试用例集
  await client.evals.deleteTestCaseSet(testCaseSet.id);
  console.log("✅ 测试用例集已删除");

  // 删除 Agents
  await client.agents.delete(agent.id);
  await client.agents.delete(agent2.id);
  console.log("✅ Agents 已删除");

  console.log("\n" + "=".repeat(70));
  console.log("✅ 演示完成！");
  console.log("=".repeat(70));

  console.log("\n📝 总结:");
  console.log("本示例展示了 Eval 的完整功能：");
  console.log("  1. ✅ 测试用例集管理");
  console.log("  2. ✅ 快速单次评估");
  console.log("  3. ✅ 批量评估");
  console.log("  4. ✅ Benchmark 多个 Agents");
  console.log("  5. ✅ A/B 测试和统计分析");
  console.log("  6. ✅ 报告生成（HTML/Markdown/JSON）");
  console.log("  7. ✅ Eval 管理和查询");
  console.log("  8. ✅ 测试用例集管理");
  console.log("  9. ✅ 资源清理");
}

// 运行示例
main().catch((error) => {
  console.error("❌ 错误:", error);
  process.exit(1);
});
