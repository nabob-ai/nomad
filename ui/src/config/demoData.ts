/**
 * 演示数据配置
 * UI 独立的演示数据，不依赖后端
 *
 * 使用场景：
 * 1. DEMO_MODE=true: UI 使用本地演示数据（默认，适合展示）
 * 2. DEMO_MODE=false: UI 连接真实后端 API（适合开发/生产）
 */

export const DEMO_MODE = import.meta.env.VITE_DEMO_MODE !== "false"; // 默认启用演示模式

// 演示工作流数据
export const demoWorkflows = [
  {
    id: "demo-wf-1",
    name: "客户服务工作流",
    description: "自动处理客户咨询和问题",
    status: "idle",
    version: "1.0.0",
    steps: [
      { id: "s1", name: "接收客户咨询", status: "pending", type: "task", config: {} },
      { id: "s2", name: "分析问题类型", status: "pending", type: "task", config: {} },
      { id: "s3", name: "生成解决方案", status: "pending", type: "task", config: {} },
      { id: "s4", name: "发送回复", status: "pending", type: "task", config: {} },
    ],
    createdAt: "2024-11-20",
  },
  {
    id: "demo-wf-2",
    name: "内容生成工作流",
    description: "自动生成营销内容和文案",
    status: "idle",
    version: "1.0.0",
    steps: [
      { id: "s1", name: "分析目标受众", status: "pending", type: "task", config: {} },
      { id: "s2", name: "生成内容大纲", status: "pending", type: "task", config: {} },
      { id: "s3", name: "撰写正文", status: "pending", type: "task", config: {} },
      { id: "s4", name: "审核优化", status: "pending", type: "task", config: {} },
    ],
    createdAt: "2024-11-19",
  },
  {
    id: "demo-wf-3",
    name: "数据分析工作流",
    description: "分析业务数据并生成报告",
    status: "completed",
    version: "1.0.0",
    steps: [
      { id: "s1", name: "收集数据", status: "completed", type: "task", config: {} },
      { id: "s2", name: "数据清洗", status: "completed", type: "task", config: {} },
      { id: "s3", name: "统计分析", status: "completed", type: "task", config: {} },
      { id: "s4", name: "生成报告", status: "completed", type: "task", config: {} },
    ],
    createdAt: "2024-11-18",
  },
];

// 演示房间数据
export const demoRooms = [
  {
    id: "demo-room-1",
    name: "产品开发团队",
    description: "产品经理、设计师、开发者协作",
    members: 5,
    agents: ["产品 Agent", "设计 Agent", "开发 Agent"],
    status: "active",
    maxMembers: 10,
    createdAt: "2024-11-20",
  },
  {
    id: "demo-room-2",
    name: "营销策划组",
    description: "营销团队协作制定推广策略",
    members: 3,
    agents: ["营销 Agent", "文案 Agent"],
    status: "active",
    maxMembers: 8,
    createdAt: "2024-11-19",
  },
  {
    id: "demo-room-3",
    name: "客户支持中心",
    description: "客服团队协作处理客户问题",
    members: 8,
    agents: ["客服 Agent", "技术支持 Agent", "售后 Agent"],
    status: "active",
    maxMembers: 15,
    createdAt: "2024-11-18",
  },
];

// 演示 Agent 数据
export const demoAgents = [
  {
    id: "demo-agent-1",
    name: "写作助手",
    description: "帮助你创作优质内容",
    status: "idle",
    template_id: "chat",
    model_config: {
      provider: "anthropic",
      model: "claude-3-5-sonnet",
    },
    createdAt: "2024-11-20",
  },
  {
    id: "demo-agent-2",
    name: "代码助手",
    description: "协助编写和优化代码",
    status: "idle",
    template_id: "chat",
    model_config: {
      provider: "openai",
      model: "gpt-4",
    },
    createdAt: "2024-11-19",
  },
  {
    id: "demo-agent-3",
    name: "数据分析师",
    description: "分析数据并生成洞察",
    status: "idle",
    template_id: "chat",
    model_config: {
      provider: "deepseek",
      model: "deepseek-chat",
    },
    createdAt: "2024-11-18",
  },
];

export const demoPassport = [
  {
    position: 1,
    rank: 1,
    passport: "Singapore",
    code: "SG",
    region: "Asia",
    score: 197
  },
  {
    position: 2,
    rank: 2,
    passport: "Japan",
    code: "JP",
    region: "Asia",
    score: 193
  },
  {
    position: 3,
    rank: 3,
    passport: "Switzerland",
    code: "CH",
    region: "Europe",
    score: 192
  },
  {
    position: 4,
    rank: 4,
    passport: "Italy",
    code: "IT",
    region: "Europe",
    score: 191
  },
  {
    position: 5,
    rank: 4,
    passport: "South Korea",
    code: "KR",
    region: "Asia",
    score: 191
  },
  {
    position: 6,
    rank: 5,
    passport: "Finland",
    code: "FI",
    region: "Europe",
    score: 190
  },
  {
    position: 7,
    rank: 5,
    passport: "Germany",
    code: "DE",
    region: "Europe",
    score: 190
  },
  {
    position: 8,
    rank: 5,
    passport: "Spain",
    code: "ES",
    region: "Europe",
    score: 190
  },
  {
    position: 9,
    rank: 5,
    passport: "Sweden",
    code: "SE",
    region: "Europe",
    score: 190
  },
  {
    position: 10,
    rank: 6,
    passport: "Belgium",
    code: "BE",
    region: "Europe",
    score: 189
  },
  {
    position: 11,
    rank: 6,
    passport: "France",
    code: "FR",
    region: "Europe",
    score: 189
  },
  {
    position: 12,
    rank: 6,
    passport: "Luxembourg",
    code: "LU",
    region: "Europe",
    score: 189
  },
  {
    position: 13,
    rank: 6,
    passport: "Norway",
    code: "NO",
    region: "Europe",
    score: 189
  },
  {
    position: 14,
    rank: 7,
    passport: "United Kingdom",
    code: "GB",
    region: "Europe",
    score: 188
  },
  {
    position: 15,
    rank: 8,
    passport: "Ireland",
    code: "IE",
    region: "Europe",
    score: 187
  },
  {
    position: 16,
    rank: 8,
    passport: "Malta",
    code: "MT",
    region: "Europe",
    score: 187
  },
  {
    position: 17,
    rank: 8,
    passport: "Netherlands",
    code: "NL",
    region: "Europe",
    score: 187
  },
  {
    position: 18,
    rank: 8,
    passport: "New Zealand",
    code: "NZ",
    region: "Oceania",
    score: 187
  },
  {
    position: 19,
    rank: 8,
    passport: "Portugal",
    code: "PT",
    region: "Europe",
    score: 187
  },
  {
    position: 20,
    rank: 9,
    passport: "Austria",
    code: "AT",
    region: "Europe",
    score: 186
  },
  {
    position: 21,
    rank: 9,
    passport: "Denmark",
    code: "DK",
    region: "Europe",
    score: 186
  },
  {
    position: 22,
    rank: 9,
    passport: "Greece",
    code: "GR",
    region: "Europe",
    score: 186
  },
  {
    position: 23,
    rank: 9,
    passport: "Malaysia",
    code: "MY",
    region: "Asia",
    score: 186
  },
  {
    position: 24,
    rank: 10,
    passport: "Canada",
    code: "CA",
    region: "Americas",
    score: 185
  },
  {
    position: 25,
    rank: 11,
    passport: "Australia",
    code: "AU",
    region: "Oceania",
    score: 184
  },
  {
    position: 26,
    rank: 11,
    passport: "Lithuania",
    code: "LT",
    region: "Europe",
    score: 184
  },
  {
    position: 27,
    rank: 11,
    passport: "Slovakia",
    code: "SK",
    region: "Europe",
    score: 184
  },
  {
    position: 28,
    rank: 11,
    passport: "United States",
    code: "US",
    region: "Americas",
    score: 184
  },
  {
    position: 29,
    rank: 12,
    passport: "Czechia",
    code: "CZ",
    region: "Europe",
    score: 183
  },
  {
    position: 30,
    rank: 12,
    passport: "Estonia",
    code: "EE",
    region: "Europe",
    score: 183
  },
  {
    position: 31,
    rank: 12,
    passport: "Hungary",
    code: "HU",
    region: "Europe",
    score: 183
  },
  {
    position: 32,
    rank: 12,
    passport: "Iceland",
    code: "IS",
    region: "Europe",
    score: 183
  },
  {
    position: 33,
    rank: 12,
    passport: "Latvia",
    code: "LV",
    region: "Europe",
    score: 183
  },
  {
    position: 34,
    rank: 12,
    passport: "Poland",
    code: "PL",
    region: "Europe",
    score: 183
  },
  {
    position: 35,
    rank: 13,
    passport: "Liechtenstein",
    code: "LI",
    region: "Europe",
    score: 182
  },
  {
    position: 36,
    rank: 13,
    passport: "Slovenia",
    code: "SI",
    region: "Europe",
    score: 182
  },
  {
    position: 37,
    rank: 14,
    passport: "Croatia",
    code: "HR",
    region: "Europe",
    score: 181
  },
  {
    position: 38,
    rank: 15,
    passport: "Monaco",
    code: "MC",
    region: "Europe",
    score: 180
  },
  {
    position: 39,
    rank: 15,
    passport: "United Arab Emirates",
    code: "AE",
    region: "Middle East",
    score: 180
  },
  {
    position: 40,
    rank: 16,
    passport: "Cyprus",
    code: "CY",
    region: "Europe",
    score: 176
  },
  {
    position: 41,
    rank: 16,
    passport: "Romania",
    code: "RO",
    region: "Europe",
    score: 176
  },
  {
    position: 42,
    rank: 17,
    passport: "Chile",
    code: "CL",
    region: "Americas",
    score: 175
  },
  {
    position: 43,
    rank: 18,
    passport: "Argentina",
    code: "AR",
    region: "Americas",
    score: 174
  },
  {
    position: 44,
    rank: 18,
    passport: "Bulgaria",
    code: "BG",
    region: "Europe",
    score: 174
  },
  {
    position: 45,
    rank: 19,
    passport: "Andorra",
    code: "AD",
    region: "Europe",
    score: 173
  },
  {
    position: 46,
    rank: 20,
    passport: "Brazil",
    code: "BR",
    region: "Americas",
    score: 172
  },
  {
    position: 47,
    rank: 20,
    passport: "Hong Kong (SAR China)",
    code: "HK",
    region: "Asia",
    score: 172
  },
  {
    position: 48,
    rank: 21,
    passport: "San Marino",
    code: "SM",
    region: "Europe",
    score: 171
  },
  {
    position: 49,
    rank: 22,
    passport: "Barbados",
    code: "BB",
    region: "Caribbean",
    score: 167
  },
  {
    position: 50,
    rank: 22,
    passport: "Brunei",
    code: "BN",
    region: "Asia",
    score: 167
  }
];

