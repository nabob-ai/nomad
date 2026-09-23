---
title: 部署指南总览
description: 将 nomad 应用部署到不同环境
navigation: false
---

# 部署指南总览

本指南涵盖从本地开发到生产部署的完整流程。

## 📚 部署选项

### [本地部署](/deployment/local)

- 开发环境配置
- HTTP Server 启动
- 工作流 HTTP API

### [Docker 部署](/deployment/docker)

- Dockerfile 配置
- 容器化最佳实践
- Docker Compose 编排

### [Kubernetes 部署](/deployment/kubernetes)

- K8s 配置文件
- 服务发现
- 自动扩缩容

### [Serverless 部署](/deployment/serverless)

- Lambda/Cloud Functions
- 冷启动优化
- 成本控制

### [云端沙箱](/deployment/cloud-sandbox)

- 阿里云 AgentBay
- 火山引擎集成

## 🚀 快速开始

```bash
# 本地运行
go run main.go

# Docker 部署
docker build -t my-agent .
docker run -p 8080:8080 my-agent

# K8s 部署
kubectl apply -f deployment.yaml
```

## 📖 相关文档

- [最佳实践：部署](/best-practices/deployment)
- [可观测性：监控](/observability/monitoring)
