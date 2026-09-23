# 技术架构

nomad 采用四层架构设计，将后端核心、Server 层、HTTP 层和客户端 SDK 完全解耦：

```
┌─────────────────────────────────────────────────────────┐
│                      Agent Layer                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │             Middleware Stack                     │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  SummarizationMiddleware (priority: 40)   │  │  │
│  │  │  - Auto summarize (>170k tokens)         │  │  │
│  │  │  - Keep recent 6 messages                 │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  FilesystemMiddleware (priority: 100)     │  │  │
│  │  │  - Tools: [fs_*, bash_*]                  │  │  │
│  │  │  - Auto eviction (>20k tokens)            │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  SubAgentMiddleware (priority: 200)       │  │  │
│  │  │  - Tools: [task]                          │  │  │
│  │  │  - Context isolation                      │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  CustomMiddleware (priority: 500+)        │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────┘  │
│         │ WrapModelCall      │ WrapToolCall            │
│         ▼                    ▼                         │
│  ┌─────────────┐      ┌──────────────┐                │
│  │  Provider   │      │ Tool Executor│                │
│  │  (Stream)   │      │  (Execute)   │                │
│  └─────────────┘      └──────────────┘                │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ StateBackend │  │ StoreBackend │  │FilesystemBE  │  │
│  │  (临时内存)  │  │  (持久化)    │  │  (真实FS)    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                ┌──────────────────┐                     │
│                │ CompositeBackend │                     │
│                │   (路由组合)     │                     │
│                └──────────────────┘                     │
└─────────────────────────────────────────────────────────┘
```

# 技术选型

- 后端编程语言： golang ,1.24 版本或更高版本
- 数据库：sqlite、postgres + redis
- 前端：vite + typescript + vue3.x

# 部署运行

## 本地开发环境说明

### 系统环境：

- 操作系统：Ubuntu 24.04.4 LTS
- 处理器: 13th Gen Intel® Core™ i9-13900HX × 32
- 显卡：NVIDIA GeForce RTX™ 4090 Laptop GPU
- 内存： 64 GiB

### 开发环境：

- golang: 1.27.1 版本
- nodejs: v24.21.0
- npm: 11.19.0
- pnpm: 10.12.1
- 集成 ide: GoLand by JetBrains
- 调试工具: Delve：Go调试器; pprof：性能分析
- 代码质量： golangci-lint：代码检查; gofmt：代码格式化（Go自带）; vite test

## Ollama 本地大模型环境


### 配置英伟达显卡对 docker 环境的支持

执行下述命令前确保 docker 环境已经安装。

```bash

curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey \
    | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
curl -fsSL https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list \
    | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' \
    | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt-get update

sudo apt-get install -y nvidia-container-toolkit

sudo nvidia-ctk runtime configure --runtime=docker

sudo systemctl restart docker

```

### 拉取 Docker Ollama 镜像

执行下述命令前确保系统当前用户属于 docker 用户组以及 root 用户组。

```bash

docker pull ollama/ollama:latest

```

### 拉取 Ollama 支持的本地模型

```bash

docker exec -it ollama ollama pull llama3.2:3b

docker exec -it ollama ollama pull llama4:17b

docker exec -it ollama ollama pull qwen3.5:9b

docker exec -it ollama ollama pull gemma4:12b

docker exec -it ollama ollama pull devstral-small-2:24b

docker exec -it ollama ollama pull gpt-oss:20b

docker exec -it ollama ollama pull gemma4:26b

docker exec -it ollama ollama pull qwen3.5:27b

```

上述本地模型，根据本地计算机设备配置，拉取合适的即可，本程序开发测试选择的 gpt-oss:20b 模型。

### 启动 Ollama 支持的本地模型服务

```bash

docker run -d --gpus all -e OLLAMA_NUM_PARALLEL=4 -e OLLAMA_MAX_QUEUE=512 -e OLLAMA_KEEP_ALIVE=-1 -e FLASH_ATTENTION=1 -e SWA_FULL=1 -e OLLAMA_NO_KV_UNIFIED=1 -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama:latest

```

Ollama 启动参数说明：

- OLLAMA_NUM_PARALLEL=4：并行请求数。允许 Ollama 同时处理最多 4 个 推理请求。默认通常为 1。设置为 4 适合多用户或多 Agent 同时调用。
- OLLAMA_MAX_QUEUE=512：最大排队队列。当并发请求超过 4 个时，最多允许 512 个 请求在队列中等待，超过此数量的请求将被拒绝。
- OLLAMA_KEEP_ALIVE=-1：模型常驻显存。设置模型加载后的常驻时间。-1 表示永久有效，即模型加载到显存后永远不会因为超时而被释放，从而彻底消除下一次调用时的“首次加载延迟”。
- FLASH_ATTENTION=1：启用闪速注意力机制。激活 FlashAttention 优化算法，能大幅减少大模型在处理长文本（长上下文）时的显存占用并提升生成速度。
- SWA_FULL=1：启用滑动窗口注意力（Sliding Window Attention）。一种针对特定模型（如 Mistral）的长文本优化机制，允许模型处理超越传统上下文长度限制的文本。
- OLLAMA_NO_KV_UNIFIED=1：禁用统一 KV 缓存。针对特定高并发场景的底层微调参数，用于调整键值（KV）缓存的管理策略，通常配合多并行（Parallel）设置来尝试优化显存分配。


运行模型服务

```bash

docker exec -it ollama ollama run gpt-oss:20b

```

查看 Ollama 服务运行日志

```bash

docker logs -f ollama

```

## 构建运行 nomad 程序

需要启动 2 同目录的个命名行窗口，前后端各一个。

### 构建运行后端服务程序

在项目根目录下执行下述命令

```bash

make build

PROVIDER=ollama MODEL=gpt-oss:20b ./bin/nomad-server

```

nomad-server 启动参数说明：

- PROVIDER=ollama：设置模型供应商为 ollama，如果不是 ollama 本地模型服务供应商，也可以定义设置，比如：deepseek、anthropic、glm、openai、gemini、openrouter。
- MODEL=gpt-oss:20b：指定模型供应商提供的具体模型型号名称。比如：llama3.2:3b、llama4:17b、qwen3.5:9b、gemma4:12b、devstral-small-2:24b、gpt-oss:20b、gemma4:26b、qwen3.5:27b
- {PROVIDER}_API_KEY="大模型供应商平台的接口访问密钥”: 设置商用大模型供应商平台接口密钥，比如供应是 deepseek，那么该参数就是 DEEPSEEK_API_KEY; 供应商是 openai，那么参数就是 OPENAI_API_KEY,因为 ollama 本地模型服务访问接口不需要密钥，所以不用设置。

### 构建运行前端程序

在项目根目录执行下述命令

```bash

cd ui

pnpm run dev

```

打开浏览器访问地址：http://localhost:3000

### 清理

如果是重新运行前后端程序，请执行下述命令。

```bash

make clean

```

## 运行效果图

### 1. 首页

![home](/screenshot/home.png)

### 2. 人工智能助手

![aidea](/screenshot/aidea.png)

### 2. 人工智能问答

![assistant](/screenshot/assistant.png)

### 3. 护照排名

![passport](/screenshot/passport.png)

## 程序详细设计文档

# nomad 文档

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