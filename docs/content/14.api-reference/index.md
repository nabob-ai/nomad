---
title: API 参考
description: 完整的 API 文档和类型定义
navigation: false
---

# API 参考

完整的 nomad API 文档。

## 📚 分类

### [Agent API](/api-reference/agent)

- agent.Create
- agent.Chat
- agent.RegisterTool
- agent.RegisterMiddleware
- agent.Subscribe
- agent.Unsubscribe

### [Provider API](/api-reference/provider)

- provider.Anthropic
- provider.OpenAI
- provider.DeepSeek
- provider.Gemini
- provider.Custom

### [Tools API](/api-reference/tools)

- tools.BashTool
- tools.FileSystemTool
- tools.HTTPTool
- tools.MCPTool
- tools.CreateCustomTool

### [Middleware API](/api-reference/middleware)

- middleware.Filesystem
- middleware.SubAgent
- middleware.Summarization
- middleware.Memory
- middleware.Custom

### [Memory API](/api-reference/memory)

- memory.SessionStore
- memory.VectorStore
- memory.Backends

### [Workflow API](/api-reference/workflow)

- workflow.Create
- workflow.Execute
- workflow.Monitor

### [Events API](/api-reference/events)

- events.ProgressEvent
- events.ControlEvent
- events.MonitorEvent

### [Types](/api-reference/types)

- types.AgentConfig
- types.ModelConfig
- types.ToolConfig

## 🚀 快速查找

使用搜索功能快速定位你需要的 API。

## 📖 相关文档

- [代码示例](/examples)
- [教程指南](/guides)
