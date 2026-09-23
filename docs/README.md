# nomad 文档

nomad 文档站点源码。

## 🚀 本地开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问: http://localhost:3000/nomad/

## 📦 构建

```bash
# 生成静态文件
npm run generate

# 预览构建结果
npm run preview
```

## 📝 文档结构

```
docs/
├── content/
│   ├── index.md                 # 首页
│   ├── 01.introduction/         # 介绍
│   ├── 02.core-concepts/        # 核心概念
│   ├── 03.providers/            # 模型与 Provider
│   ├── 04.memory/               # 记忆
│   ├── 05.tools/                # 工具系统
│   ├── 06.middleware/           # 中间件
│   ├── 07.workflows/            # 工作流
│   ├── 08.multi-agent/          # 多Agent系统
│   ├── 09.deployment/           # 部署指南
│   ├── 10.observability/        # 可观测性
│   ├── 11.evals/                # 评估系统
│   ├── 12.examples/             # 代码示例
│   ├── 13.guides/               # 实战指南
│   ├── 14.api-reference/        # API参考
│   └── 15.best-practices/       # 最佳实践
├── public/
│   └── images/                  # 图片资源
├── components/                  # Vue组件
├── layouts/                     # 布局模板
└── nuxt.config.ts              # 配置文件
```

## 🛠️ 技术栈

- [Nuxt 3](https://nuxt.com/) - Vue框架
- [Nuxt Content](https://content.nuxt.com/) - 内容管理
- [Tailwind CSS](https://tailwindcss.com/) - 样式框架
