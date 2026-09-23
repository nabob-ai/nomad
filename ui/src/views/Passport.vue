<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-6 py-8">
      <div class="mb-8 flex items-center justify-between">
        <div>
          <router-link to="/" class="text-blue-600 dark:text-blue-400 hover:underline mb-4 inline-block"> ← 返回首页 </router-link>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">全球护照排名</h1>
          <p class="text-gray-600 dark:text-gray-400 mt-2"> 全球护照排名 TOP 50 </p>
        </div>
      </div>

      <PassportList :passports="passports" :loading="loading" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import PassportList from "../components/Passport/PassportList.vue";
import { DEMO_MODE, demoPassport } from "../config/demoData";

const passports = ref<any[]>([]);

console.log("🎭 演示模式:", DEMO_MODE ? "启用（使用本地数据）" : "禁用（连接后端API）");

const loading = ref(false);

// Read - 加载房间列表
const loadPassports = async () => {
  try {
    loading.value = true;

    if (DEMO_MODE) {
      // 演示模式：使用 UI 本地数据
      passports.value = JSON.parse(JSON.stringify(demoPassport));
    } else {
      // 生产模式：从后端 API 获取
    }
  } catch (error: any) {
    console.error("加载房间失败:", error);
    // 失败时使用演示数据作为后备
    passports.value = JSON.parse(JSON.stringify(demoPassport));
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadPassports();
});
</script>
