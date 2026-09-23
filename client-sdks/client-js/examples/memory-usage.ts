/**
 * Memory 三层记忆系统使用示例
 */

import { MemoryResource } from "@nomad/client-js";

async function main() {
  // 创建 Memory 资源
  const memory = new MemoryResource({
    baseUrl: "http://localhost:8080",
    apiKey: process.env.NOMAD_API_KEY,
  });

  console.log("=".repeat(60));
  console.log("NomadClient 三层记忆系统演示");
  console.log("=".repeat(60));

  // ========================================================================
  // 1. Working Memory 演示
  // ========================================================================
  console.log("\n📝 1. Working Memory（工作记忆）");
  console.log("-".repeat(60));

  // 设置 Thread 作用域的记忆（会话级别）
  await memory.working.set(
    "user_preference",
    {
      theme: "dark",
      language: "zh-CN",
      notifications: true,
    },
    {
      scope: "thread",
      ttl: 3600, // 1小时后过期
    },
  );
  console.log("✅ 设置 Thread 作用域记忆: user_preference");

  // 设置 Resource 作用域的记忆（全局）
  await memory.working.set(
    "app_config",
    {
      version: "1.0.0",
      features: ["chat", "workflow", "memory"],
    },
    {
      scope: "resource", // 全局作用域
      ttl: 0, // 永不过期
    },
  );
  console.log("✅ 设置 Resource 作用域记忆: app_config");

  // 带 JSON Schema 验证的记忆
  await memory.working.set(
    "validated_data",
    {
      count: 42,
      name: "test",
    },
    {
      schema: {
        type: "object",
        properties: {
          count: { type: "number" },
          name: { type: "string" },
        },
        required: ["count", "name"],
      },
    },
  );
  console.log("✅ 设置带 Schema 验证的记忆: validated_data");

  // 获取记忆
  const preference = await memory.working.get("user_preference", "thread");
  console.log("📖 读取记忆:", preference);

  // 列出所有 Thread 作用域的记忆
  const threadMemories = await memory.working.list("thread");
  console.log("📋 Thread 作用域记忆数:", Object.keys(threadMemories).length);

  // ========================================================================
  // 2. Semantic Memory 演示
  // ========================================================================
  console.log("\n🔍 2. Semantic Memory（语义记忆）");
  console.log("-".repeat(60));

  // 存储知识
  const chunk1 = await memory.semantic.store("Paris is the capital of France.", { source: "wikipedia", category: "geography", language: "en" });
  console.log("✅ 存储记忆块 1:", chunk1);

  const chunk2 = await memory.semantic.store("The Eiffel Tower is located in Paris.", { source: "wikipedia", category: "landmarks", language: "en" });
  console.log("✅ 存储记忆块 2:", chunk2);

  const chunk3 = await memory.semantic.store("France is a country in Western Europe.", { source: "wikipedia", category: "geography", language: "en" });
  console.log("✅ 存储记忆块 3:", chunk3);

  // 语义搜索
  console.log('\n🔎 搜索: "What is the capital of France?"');
  const results = await memory.semantic.search("What is the capital of France?", {
    limit: 5,
    threshold: 0.7,
    filter: { category: "geography" },
  });

  results.forEach((chunk, index) => {
    console.log(`  ${index + 1}. [Score: ${chunk.score?.toFixed(2)}] ${chunk.content}`);
    console.log(`     Metadata:`, chunk.metadata);
  });

  // ========================================================================
  // 3. Memory Provenance 演示
  // ========================================================================
  console.log("\n🔗 3. Memory Provenance（记忆溯源）");
  console.log("-".repeat(60));

  // 查询记忆溯源
  try {
    const provenance = await memory.getProvenance(chunk1);
    console.log("📊 记忆溯源信息:");
    console.log("  - 来源:", provenance.provenance.source);
    console.log("  - 置信度:", provenance.provenance.confidence);
    console.log("  - 时间:", provenance.provenance.timestamp);

    // 查询谱系链
    if (provenance.provenance.parentId) {
      const lineage = await memory.getLineage(chunk1);
      console.log("  - 谱系链长度:", lineage.length);
    }
  } catch (error: any) {
    console.log("⚠️  溯源功能未启用或记忆未找到:", error.message);
  }

  // ========================================================================
  // 4. Memory Consolidation 演示
  // ========================================================================
  console.log("\n🔄 4. Memory Consolidation（记忆合并）");
  console.log("-".repeat(60));

  try {
    // 触发记忆合并
    const consolidation = await memory.consolidate({
      strategy: "summarize",
      llmProvider: "anthropic",
      llmModel: "claude-sonnet-4",
    });

    console.log("✅ 合并任务已启动:", consolidation.jobId);
    console.log("   状态:", consolidation.status);

    // 查询合并状态
    const status = await memory.getConsolidationStatus(consolidation.jobId);
    console.log("📊 合并进度:", status.progress, "%");

    // 如果需要取消
    // await memory.cancelConsolidation(consolidation.jobId);
  } catch (error: any) {
    console.log("⚠️  合并功能未启用:", error.message);
  }

  // ========================================================================
  // 5. 统计信息
  // ========================================================================
  console.log("\n📊 5. 统计信息");
  console.log("-".repeat(60));

  try {
    const stats = await memory.getStats();
    console.log("Working Memory:");
    console.log("  - Thread 记忆数:", stats.workingMemory.threadCount);
    console.log("  - Resource 记忆数:", stats.workingMemory.resourceCount);
    console.log("  - 总大小:", stats.workingMemory.totalSize, "bytes");

    console.log("Semantic Memory:");
    console.log("  - 记忆块数:", stats.semanticMemory.chunkCount);
    console.log("  - 总大小:", stats.semanticMemory.totalSize, "bytes");
  } catch (error: any) {
    console.log("⚠️  统计信息获取失败:", error.message);
  }

  // ========================================================================
  // 6. 清理
  // ========================================================================
  console.log("\n🧹 6. 清理（可选）");
  console.log("-".repeat(60));

  // 删除单个记忆
  await memory.working.delete("validated_data", "thread");
  console.log("✅ 删除 Working Memory: validated_data");

  // 删除 Semantic Memory
  await memory.semantic.delete(chunk3);
  console.log("✅ 删除 Semantic Memory:", chunk3);

  // 批量删除
  await memory.semantic.deleteBatch([chunk1, chunk2]);
  console.log("✅ 批量删除 Semantic Memory");

  // 清空 Thread 作用域的所有记忆
  await memory.working.clear("thread");
  console.log("✅ 清空 Thread 作用域记忆");

  // 危险：清空所有记忆
  // await memory.clearAll(true);

  console.log("\n" + "=".repeat(60));
  console.log("✅ 演示完成！");
  console.log("=".repeat(60));
}

// 运行示例
main().catch((error) => {
  console.error("❌ 错误:", error);
  process.exit(1);
});
