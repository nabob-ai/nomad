<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
    <!-- Header -->
    <header class="border-b border-gray-200 dark:border-gray-700 bg-white/80 dark:bg-gray-800/80 backdrop-blur-sm sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-xl">A</span>
            </div>
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Nomad Concierge</h1>
              <p class="text-sm text-gray-600 dark:text-gray-400">一个面向计划到海外生活的用户的服务入口</p>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Navigation Cards -->
    <section class="max-w-7xl mx-auto px-6 pb-16">
      <h3 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">快速导航</h3>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- AI Assistant -->
        <router-link to="/aidea" class="group bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-blue-500 dark:hover:border-blue-400 transition-all">
          <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <span class="text-2xl">🤖</span>
          </div>
          <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">人工智能助手</h4>
          <p class="text-gray-600 dark:text-gray-400 mb-4">使用人工智能助手帮您规划海外生活</p>
          <div class="text-blue-600 dark:text-blue-400 font-medium group-hover:translate-x-1 transition-transform inline-block">立即体验 →</div>
        </router-link>

        <!-- Component passport -->
        <router-link to="/passport" class="group bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-purple-500 dark:hover:border-purple-400 transition-all">
          <div class="w-12 h-12 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <span class="text-2xl">📖</span>
          </div>
          <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">护照排名</h4>
          <p class="text-gray-600 dark:text-gray-400 mb-4">查看多国护照比较</p>
          <div class="text-purple-600 dark:text-purple-400 font-medium group-hover:translate-x-1 transition-transform inline-block">查看护照信息 →</div>
        </router-link>

        <!-- Agent Management -->
        <router-link to="/agents" class="group bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-orange-500 dark:hover:border-orange-400 transition-all">
          <div class="w-12 h-12 bg-orange-100 dark:bg-orange-900/30 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <span class="text-2xl">⚙️</span>
          </div>
          <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Agent 管理</h4>
          <p class="text-gray-600 dark:text-gray-400 mb-4">创建、配置和管理你的 AI Agents</p>
          <div class="text-orange-600 dark:text-orange-400 font-medium group-hover:translate-x-1 transition-transform inline-block">管理 Agents →</div>
        </router-link>

        <!-- Workflows -->
        <router-link to="/workflows" class="group bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-pink-500 dark:hover:border-pink-400 transition-all">
          <div class="w-12 h-12 bg-pink-100 dark:bg-pink-900/30 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <span class="text-2xl">🔄</span>
          </div>
          <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">工作流</h4>
          <p class="text-gray-600 dark:text-gray-400 mb-4">可视化 Agent 工作流和执行步骤</p>
          <div class="text-pink-600 dark:text-pink-400 font-medium group-hover:translate-x-1 transition-transform inline-block">查看工作流 →</div>
        </router-link>

        <!-- Rooms -->
        <router-link to="/rooms" class="group bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-indigo-500 dark:hover:border-indigo-400 transition-all">
          <div class="w-12 h-12 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <span class="text-2xl">👥</span>
          </div>
          <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">协作房间</h4>
          <p class="text-gray-600 dark:text-gray-400 mb-4">多 Agent 协作和团队工作空间</p>
          <div class="text-indigo-600 dark:text-indigo-400 font-medium group-hover:translate-x-1 transition-transform inline-block">进入房间 →</div>
        </router-link>
      </div>
    </section>

    <!-- Features -->
    <section class="max-w-7xl mx-auto px-6 pb-16">
      <h3 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">功能介绍</h3>

      <div class="grid md:grid-cols-2 gap-8">
        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">🤖</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">人工智能助手</h4>
            <p class="text-gray-600 dark:text-gray-400">支持与本地或云端大模型服务对话</p>
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">💬</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">护照排名</h4>
            <p class="text-gray-600 dark:text-gray-400">简单的护照排名</p>
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">🔄</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">工作流可视化</h4>
            <p class="text-gray-600 dark:text-gray-400">实时展示 Agent 执行步骤，支持状态追踪和错误处理</p>
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-orange-100 dark:bg-orange-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">🎨</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">现代设计</h4>
            <p class="text-gray-600 dark:text-gray-400">基于 Tailwind CSS，支持深色模式，响应式布局</p>
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-pink-100 dark:bg-pink-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">💪</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">TypeScript 支持</h4>
            <p class="text-gray-600 dark:text-gray-400">完整的类型定义，提供优秀的开发体验</p>
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-shrink-0 w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center">
            <span class="text-xl">⚡</span>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">高性能</h4>
            <p class="text-gray-600 dark:text-gray-400">基于 Vue 3 Composition API，优化的渲染性能</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-gray-200 dark:border-gray-700 bg-white/80 dark:bg-gray-800/80 backdrop-blur-sm">
      <div class="max-w-7xl mx-auto px-6 py-8">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <div class="text-gray-600 dark:text-gray-400">© 2024 Nabob AI. Apache-2.0 license.</div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// Home page - Navigation hub for Nomad Agent UI
</script>
