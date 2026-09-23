# AI Usage

## 使用的工具

| 工具 | 用途 | 使用范围 |
| --- | --- | --- |
| chatgpt | 搜集移民、护照相关数据  | 数据整理 |
| qwen3-coder:30b | 辅助修改代码  | 技术实现 |

## 关键交互摘要

### 交互 1

#### 我的目标或提示词摘要：

抓取 https://multipassrank.com/抓取上述网站 全球护照排名数据,梳理成 json 数据格式文件。

#### AI 给出的建议： 
给出了核心 json 格式的数据字段，以及 json 数据文件。

```json
{
  "source": "MultiPass Rank",
  "source_url": "https://multipassrank.com/",
  "dataset_type": "global_passport_ranking",
  "checked_date": "2026-09-18",
  "retrieved_date": "2026-09-23",
  "methodology_summary": "Ranked by destinations available without prior visa approval: visa-free, visa on arrival, or ETA.",
  "passport_count": 199,
  "tracked_destinations": 227,
  "records": [
    {
      "position": 1,
      "rank": 1,
      "passport": "Singapore",
      "code": "SG",
      "region": "Asia",
      "score": 197
    },
    {...}
    ]
 }
```

#### 我接受、修改或拒绝了什么：

只接受了 json 数据结构里的 records 部分，删减了其它不需要的数据。

### 交互 2

#### 我的目标或提示词摘要：

ollama 本地环境，运行 gpt-oss:20b，与模型交互过程中出现警告信息 `forcing full prompt re-processing` 应该如何优化？

#### AI 给出的建议：

根据 llama.cpp 和新版 Ollama 底层的流式缓存机制：iSWA 架构的特性： 为了极大地节省运行 gpt-oss:20b 这类大模型时的显存（KV-Cache）消耗，模型内部某些层只保留最近一个固定窗口内（例如最近 2048 或 4096 个 Token）的记忆。缓存失效冲突： 当你的 Agent 携带多轮工具链历史（Tool Calls）或者执行多轮 Agent 循环时，每次请求的上下文都在动态增加或被局部修改。由于 iSWA 会在窗口滑动时“自动遗忘/裁剪”掉较老的 Token，这导致大模型服务端发现无法直接无缝复用（Reuse）上一轮的完整 KV 缓存。后果： Ollama 只能被迫打印出这行警告，提示它必须丢弃当前不完整的缓存，对你发过去的全部 Context（提示词+前几轮历史）进行全量的重新计算（Full Re-processing）。这就解释了为什么当对话变长后，模型在吐字前会突然卡住很长一段时间。


#### 我接受、修改或拒绝了什么：

根据建议在 Ollama docker 容器启动时加入环境变量参数：

```
-e OLLAMA_MAX_QUEUE=512 
-e OLLAMA_KEEP_ALIVE=-1 
-e FLASH_ATTENTION=1 
-e SWA_FULL=1 
-e OLLAMA_NO_KV_UNIFIED=1
```

## 发现并修正的问题

由于 iSWA 模型遇到极长的上下文时必定会由于滑窗引发抖动,让 AI 给出优化建议。
刚开始给出的错误答案是在调整 golang 代码，修改 websocket 永不超时。

人工修复方法：发起模型执行循环前，限制历史消息的最大长度。

## 验证方式

### 我运行的自动化检查：
- 代码检查： `make lint`
- 测试代码： `make test`

### 我执行的人工场景：

在 golang 关键代码中逐步打印运行日志，定位 bug。

### 我特别检查的安全或数据风险：

特别检查的问题是：大模型接口密钥泄漏问题。起初使用的商用大模型供应 openrouter，又来切换到 ollama 本地大模型服务，这个问题暂时不需要考虑了。
