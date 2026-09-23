<template>
  <div class="min-h-screen bg-background text-text dark:bg-background-dark dark:text-text-dark selection:bg-sky-500/40" :style="{ background: 'radial-gradient(circle at top, rgba(59, 130, 246, 0.12), transparent 60%), var(--tw-bg-opacity, 1)' }">
    <!-- Hero -->
    <section id="live" class="relative overflow-hidden px-6 py-10 sm:py-16 lg:px-12">
      <div class="absolute inset-0 pointer-events-none">
        <div class="absolute inset-y-0 inset-x-10 bg-gradient-to-r from-sky-500/20 via-blue-500/10 to-transparent blur-3xl"></div>
        <div class="absolute right-0 top-10 w-96 h-96 bg-purple-500/20 blur-3xl rounded-full"></div>
      </div>
      <div class="relative z-10 max-w-6xl mx-auto">
        <header class="flex flex-wrap items-center justify-between gap-6 pb-10 border-b border-white/10">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-sky-500 to-blue-600 flex items-center justify-center font-black text-xl">A</div>
            <div>
              <p class="text-xl font-semibold">Nomad Agent Studio</p>
              <p class="text-sm text-white/60">UI · Demo · Docs</p>
            </div>
          </div>
          <nav class="flex flex-wrap items-center gap-4 text-sm text-white/70">
            <a v-for="link in navLinks" :key="link.href" class="hover:text-white transition-colors" :href="link.href">{{ link.label }}</a>
            <a href="https://github.com/nabob-ai/nomad" target="_blank" class="inline-flex items-center gap-2 rounded-full bg-white text-slate-900 px-4 py-2 font-semibold hover:-translate-y-0.5 transition-all">
              GitHub
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 17L17 7M7 7h10v10" />
              </svg>
            </a>
          </nav>
        </header>

        <div class="grid gap-10 lg:grid-cols-[1.1fr,0.9fr] items-start pt-10">
          <div class="space-y-6">
            <p class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-white/10 border border-white/20 w-fit">
              <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              Nomad 0.14 · Chat + Workflow + Docs
            </p>
            <h1 class="text-4xl sm:text-5xl font-semibold leading-tight tracking-tight">
              构建一站式 Agent UI，
              <span class="text-sky-400">既能演示体验又能交付文档</span>
            </h1>
            <p class="text-lg text-white/70 max-w-xl">参考阿里云 ChatUI 的信息架构，我们将 Nomad UI 目录重塑为一个集「演示、构建、文档」于一体的站点：统一入口、双模式体验、丰富的组件库说明以及嵌入式 Agent Playground。</p>
            <div class="flex flex-wrap gap-4">
              <a href="#demo" class="inline-flex items-center gap-2 rounded-2xl bg-sky-500 hover:bg-sky-400 px-6 py-3 font-semibold transition-colors">
                即刻体验
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </a>
              <a href="#docs" class="inline-flex items-center gap-2 rounded-2xl border border-white/30 px-6 py-3 font-semibold text-white/80 hover:text-white hover:border-white transition-colors"> 阅读文档 </a>
            </div>
            <div class="grid grid-cols-2 gap-6 pt-8 border-t border-white/10">
              <div v-for="metric in heroMetrics" :key="metric.label">
                <p class="text-3xl font-black">{{ metric.value }}</p>
                <p class="text-sm text-white/60">{{ metric.label }}</p>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <div class="rounded-3xl bg-white/5 backdrop-blur border border-white/10 p-4 shadow-[0_20px_80px_rgba(15,23,42,0.6)] h-[500px] flex flex-col">
              <div class="flex items-center justify-between mb-4 flex-shrink-0">
                <div>
                  <p class="text-sm text-white/60">Live Agent Preview</p>
                  <p class="font-semibold">Nomad Copilot</p>
                </div>
                <div class="flex items-center gap-2 text-xs text-emerald-400">
                  <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
                  在线
                </div>
              </div>
              <div class="flex-1 min-h-0">
                <NomadChat :config="heroDemoConfig" :show-header="false" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div v-for="card in heroCapability" :key="card.title" class="rounded-2xl border border-white/10 bg-gradient-to-br from-white/5 to-transparent p-4">
                <p class="text-sm text-white/60">{{ card.label }}</p>
                <p class="text-xl font-semibold mt-1">{{ card.title }}</p>
                <p class="text-xs text-white/50 mt-1">{{ card.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Feature Spectrum -->
    <section class="px-6 py-16 lg:px-12 bg-slate-800/30 border-y border-white/5">
      <div class="max-w-6xl mx-auto space-y-10">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p class="text-sm uppercase tracking-[0.3em] text-white/50">Capability Spectrum</p>
            <h2 class="text-3xl font-semibold mt-2">多角色、全形态的 Agent UI 设计体系</h2>
          </div>
          <span class="text-white/60 text-sm">该区块在上线环境中可配置为客户案例或产品亮点。</span>
        </div>
        <div class="grid gap-6 md:grid-cols-3">
          <article v-for="feature in featureGroups" :key="feature.title" class="rounded-3xl border border-white/10 bg-gradient-to-br from-white/5 to-transparent p-6 space-y-3">
            <div class="flex items-center gap-3 text-sm text-white/70">
              <span class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center">{{ feature.icon }}</span>
              {{ feature.badge }}
            </div>
            <h3 class="text-2xl font-semibold">{{ feature.title }}</h3>
            <p class="text-white/70 text-sm leading-relaxed">{{ feature.description }}</p>
            <ul class="text-sm text-white/60 space-y-1">
              <li v-for="item in feature.items" :key="item" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 rounded-full bg-sky-400"></span>
                {{ item }}
              </li>
            </ul>
          </article>
        </div>
      </div>
    </section>

    <!-- Playground -->
    <section id="demo" class="px-6 py-20 lg:px-12">
      <div class="max-w-6xl mx-auto space-y-10">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p class="text-sm uppercase tracking-[0.3em] text-white/50">Live Playground</p>
            <h2 class="text-3xl font-semibold mt-2">一键切换 HTTP / Streaming 的体验面板</h2>
            <p class="text-white/60 text-sm mt-2">通过双模式预置配置，演示如何使用 Nomad SDK 迅速搭建多 Agent / 流式体验。默认开启 demoMode，无需后端即可演练。</p>
          </div>
          <span class="text-white/60 text-sm">切换 demoMode=false 即可连上真实的 Nomad Server。</span>
        </div>
        <div class="grid gap-8 lg:grid-cols-[0.8fr,1.2fr]">
          <div class="space-y-6">
            <div class="flex flex-wrap gap-3">
              <button
                v-for="mode in demoModes"
                :key="mode.id"
                @click="selectMode(mode.id)"
                :class="['px-4 py-2 rounded-2xl text-sm font-semibold border transition-all', mode.id === selectedModeId ? 'border-sky-400 bg-sky-500/10 text-white' : 'border-white/10 text-white/60 hover:text-white hover:border-white/30']"
              >
                {{ mode.label }}
              </button>
            </div>
            <div class="rounded-3xl border border-white/10 bg-white/5 p-6 space-y-4">
              <div>
                <p class="text-sm text-white/50">当前模式</p>
                <p class="text-2xl font-semibold">{{ selectedMode.tagline }}</p>
              </div>
              <ul class="space-y-3 text-sm text-white/70">
                <li v-for="point in selectedMode.highlights" :key="point" class="flex items-start gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-sky-400 mt-1"></span>
                  {{ point }}
                </li>
              </ul>
              <div class="text-xs text-white/50 border-t border-white/10 pt-4">
                Call <span class="font-mono text-sky-300">{{ selectedMode.api }}</span>
              </div>
            </div>
          </div>
          <div class="rounded-[28px] border border-white/10 bg-slate-800/80 p-4 shadow-[0_40px_90px_rgba(30,41,59,0.8)] h-[600px]">
            <NomadChat :config="demoConfig" />
          </div>
        </div>
      </div>
    </section>

    <!-- Builder -->
    <section class="px-6 py-20 lg:px-12 bg-slate-800/50 border-y border-white/5">
      <div class="max-w-6xl mx-auto grid gap-10 lg:grid-cols-[1.1fr,0.9fr]">
        <div class="space-y-5">
          <p class="text-sm uppercase tracking-[0.3em] text-white/50">Builder Stories</p>
          <h2 class="text-3xl font-semibold">UI 目录不再是散乱的 demo，而是一座 <span class="text-sky-400">可扩展的 Agent 展示站</span></h2>
          <p class="text-white/70 leading-relaxed">我们将 ChatUI 风格的设计语言融入到 Nomad 的组件体系，提供完成度极高的「产品化体验」：导航、亮点、可切换的 Playground、构建步骤、以及面向开发者的文档卡片。整个目录可以直接对外展示，兼具业务拓展价值。</p>
          <div class="rounded-3xl border border-white/10 bg-gradient-to-br from-white/5 to-transparent p-6 space-y-4">
            <h3 class="text-xl font-semibold">构建 Nomad Agent UI 的四步法</h3>
            <div class="space-y-4">
              <div v-for="step in builderSteps" :key="step.title" class="flex gap-3">
                <div class="w-10 h-10 rounded-2xl bg-white/10 flex items-center justify-center font-semibold text-sky-300">
                  {{ step.step }}
                </div>
                <div>
                  <p class="font-semibold">{{ step.title }}</p>
                  <p class="text-sm text-white/60">{{ step.description }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="space-y-4">
          <div v-for="block in capabilityBlocks" :key="block.title" class="rounded-3xl border border-white/10 bg-white/[0.04] p-5">
            <p class="text-sm text-white/50">{{ block.badge }}</p>
            <h3 class="text-2xl font-semibold mt-2">{{ block.title }}</h3>
            <p class="text-sm text-white/60 leading-relaxed">{{ block.description }}</p>
            <ul class="mt-4 space-y-1 text-white/70 text-sm">
              <li v-for="item in block.points" :key="item" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                {{ item }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- Documentation -->
    <section id="docs" class="px-6 py-20 lg:px-12">
      <div class="max-w-6xl mx-auto space-y-10">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p class="text-sm uppercase tracking-[0.3em] text-white/50">Documentation & Assets</p>
            <h2 class="text-3xl font-semibold mt-2">演示、文档、开源资源，一站式呈现</h2>
          </div>
          <a href="https://chatbot.console.aliyun.com/ChatSDK/chatui" target="_blank" class="text-sky-400 text-sm flex items-center gap-2">
            参考 ChatUI
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </a>
        </div>
        <div class="grid gap-6 md:grid-cols-3">
          <article v-for="doc in docEntries" :key="doc.title" class="rounded-3xl border border-white/10 bg-white/5 p-6 flex flex-col gap-4">
            <div class="text-sm text-white/60">{{ doc.category }}</div>
            <h3 class="text-2xl font-semibold">{{ doc.title }}</h3>
            <p class="text-sm text-white/60 flex-1">{{ doc.description }}</p>
            <ul class="text-sm text-white/60 space-y-1">
              <li v-for="item in doc.points" :key="item" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 rounded-full bg-white/40"></span>
                {{ item }}
              </li>
            </ul>
            <a :href="doc.href" class="text-sky-400 text-sm flex items-center gap-2">
              {{ doc.cta }}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </a>
          </article>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="px-6 py-10 lg:px-12 border-t border-white/5 text-sm text-white/50">
      <div class="max-w-6xl mx-auto flex flex-wrap justify-between gap-4">
        <p>© {{ new Date().getFullYear() }} Nomad Cloud. 由 Nomad Agent Studio 驱动。</p>
        <div class="flex items-center gap-4">
          <a href="https://github.com/nabob-ai/nomad" target="_blank" class="hover:text-white transition-colors">GitHub</a>
          <a href="mailto:hello@nomadcloud.com" class="hover:text-white transition-colors">Contact</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import NomadChat from "@/components/NomadChat.vue";

interface DemoMode {
  id: string;
  label: string;
  tagline: string;
  highlights: string[];
  api: string;
  agentProfile: {
    name: string;
    description: string;
  };
  responses: string[];
  enableImage?: boolean;
  enableVoice?: boolean;
  delay?: number;
}

const navLinks = [
  { label: "解决方案", href: "#demo" },
  { label: "构建流程", href: "#docs" },
  { label: "体验站", href: "#live" },
];

const heroMetrics = [
  { label: "组件 / Hooks", value: "24" },
  { label: "演示模式", value: "HTTP · Streaming" },
];

const heroCapability = [
  { label: "体验", title: "Streaming + Markdown", description: "完整的内容分发体验，支持代码块渲染、引用和审批。" },
  { label: "操作", title: "工具联动", description: "内建思考块、工具调用结果、人工审批环节。" },
];

const featureGroups = [
  {
    title: "产品级体验",
    description: "ChatUI 式的结构，不只是 demo，而是对外可呈现的产品页面。",
    badge: "Experience Layer",
    icon: "✨",
    items: ["自适应导航与亮点介绍", "多模式 Playground", "预制指标与客户语境"],
  },
  {
    title: "Agent Builder",
    description: "用组件和 hooks 快速拼装 Agent 运行面板，支持多 Agent 与房间模式。",
    badge: "Builder Layer",
    icon: "🧱",
    items: ["NomadChat、Agent、Room 组件", "useChat / useNomadClient composable", "Workflow 状态流"],
  },
  {
    title: "Docs + Assets",
    description: "围绕演示整理的文档、脚手架与二开指南，适合销售&合作介绍。",
    badge: "Documentation Layer",
    icon: "📚",
    items: ["代码示例与 API 链接", "Ready-to-share 的 HTML 页面", "Tailwind 主题变量"],
  },
];

const demoModes: DemoMode[] = [
  {
    id: "http",
    label: "HTTP Chat",
    tagline: "标准 REST 模式 · 直连 Agent API",
    highlights: ["单回合/多回合历史自动挂载", "Markdown + 附件渲染", "热键 Enter 发送，支持 Quick Reply"],
    api: "POST /agents/{id}/chat",
    agentProfile: {
      name: "Nomad Architect",
      description: "负责整理上下文、编排工具的 Planner Agent",
    },
    responses: ["我已根据 {question} 的需求生成一套执行计划，包括调研、工具调用与风控审批。有什么要微调的吗？", "Nomad Builder 已更新沙箱策略，并将关键信息同步到记忆中心。你可以继续追问 {question} 相关的细节。"],
    enableImage: true,
  },
  {
    id: "stream",
    label: "Streaming",
    tagline: "SSE / WebSocket 流式体验 · 实时反馈",
    highlights: ["Token 级刷新，平均延迟 220ms", "思考块实时输出，支持人工审批", "支持语音输入与多模态锁定"],
    api: "GET /agents/{id}/events",
    agentProfile: {
      name: "Nomad Realtime Copilot",
      description: "专注多模态输入与 streaming 输出的执行 Agent",
    },
    responses: ["正在监听 streaming 频道……\n\n1. 已解析指令 {question}\n2. 正在触发工具链\n3. 预计 0.2s 内返回 token", "这类实时指令会使用 SSE 管道，以保持响应流的持续。所有输出在收敛后会自动写回消息列表。"],
    enableVoice: true,
    delay: 300,
  },
];

const selectedModeId = ref(demoModes[0]?.id ?? "http");
const selectedMode = computed(() => demoModes.find((mode) => mode.id === selectedModeId.value) ?? demoModes[0]!);
const selectMode = (modeId: string) => {
  selectedModeId.value = modeId;
};

const defaultDemoConfig = {
  placeholder: "输入你的 Agent 指令，例如：生成一个有审批的多代理工作流",
  enableThinking: true,
  enableApproval: true,
  enableQuickReplies: true,
  enableSuggestions: true,
  welcomeMessage: "👋 你好，我是 Nomad Copilot，负责把 Agent 能力产品化展示。试着问我：如何搭建一个多 Agent 工单系统？",
};

const heroDemoConfig = computed(() => ({
  ...defaultDemoConfig,
  agentProfile: {
    name: "Nomad Copilot",
    description: "负责产品演示的说明 Agent",
  },
  demoMode: false,
  apiUrl: "http://localhost:8080",
  modelConfig: {
    provider: "ollama",
    model: "gpt-oss:20b",
  },
}));

const demoConfig = computed(() => ({
  ...defaultDemoConfig,
  agentProfile: selectedMode.value?.agentProfile,
  demoMode: false,
  apiUrl: "http://localhost:8080",
  modelConfig: {
    provider: "ollama",
    model: "gpt-oss:20b",
  },
  enableImage: selectedMode.value?.enableImage,
  enableVoice: selectedMode.value?.enableVoice,
}));

const builderSteps = [
  { step: "01", title: "统一入口", description: "重新组织 UI 目录，形成导航 + Hero + Playground + Docs 的站点结构。" },
  { step: "02", title: "组件接入", description: "引入 NomadChat、Agent、Workflow 等组件，形成可复用模块。" },
  { step: "03", title: "模式扩展", description: "内建 demoMode，可离线演练；只需配置即可切换到真实 API。" },
  { step: "04", title: "文档沉淀", description: "卡片化呈现快速开始、API、主题定制，方便协作团队查阅。" },
];

const capabilityBlocks = [
  {
    badge: "Ops Team",
    title: "统一演示面板",
    description: "销售/合作团队可直接展示 UI 站点，包含 demo Chat、Streaming、审批和多语言文案。",
    points: ["参考 ChatUI 的信息密度", "Tailwind + CSS 变量即可换肤", "支持部署到任意静态托管"],
  },
  {
    badge: "Developer Experience",
    title: "组件 + Composable",
    description: "通过 @nomad/ui 暴露的组件和 hooks，开发者能快速在自己的业务系统中复用同样的体验。",
    points: ["NomadChat / Room / Workflow", "useChat / useNomadClient / useWebSocket", "类型完备的 TypeScript 定义"],
  },
];

const docEntries = [
  {
    category: "快速开始",
    title: "安装与基础用法",
    description: "一行命令安装 SDK，三行代码引入 Chat 组件，开箱即可使用。",
    points: ["npm install @nomad/ui", "import { NomadChat } ...", "支持 Vite / Nuxt / SPA"],
    cta: "查看指南",
    href: "https://github.com/nabob-ai/nomad",
  },
  {
    category: "API & Hooks",
    title: "Nomad Client 能力",
    description: "通过 useNomadClient 接入 HTTP / WS，结合 useChat 管理状态与消息。",
    points: ["agents.chat / rooms.join", "Streaming 订阅管理", "错误与审批回调"],
    cta: "阅读 API",
    href: "https://github.com/nabob-ai/nomad",
  },
  {
    category: "主题与资产",
    title: "自定义主题 / 动画 / 图标",
    description: "Tailwind 配置 + CSS 变量即可完成品牌定制，配套提供插画、动效与布局示例。",
    points: ["tailwind.config.js 主题", "可复用的 Vue 组件", "插图与动效示例"],
    cta: "获取素材",
    href: "https://github.com/nabob-ai/nomad/tree/main/ui",
  },
];
</script>
