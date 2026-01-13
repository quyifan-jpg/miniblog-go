# 🦉 MiniBlog 学习指南

## 📍 项目位置

```
advanced_ai_agents/multi_agent_apps/ai_news_and_podcast_agents/miniblog/
```

## 🎯 项目概述

**MiniBlog** 是一个完整的**生产级多 Agent 系统**，用于：
- 📰 收集和管理新闻文章
- 📱 监控社交媒体（X.com, Facebook）
- 🎙️ 自动生成播客
- 🎨 生成播客封面和脚本
- 🔊 文本转语音（TTS）

## 🏗️ 项目架构

### 核心组件

```
miniblog/
├── main.py                    # FastAPI 主应用入口
├── celery_worker.py           # Celery 异步任务处理
├── scheduler.py               # 定时任务调度
│
├── agents/                    # AI Agent 定义
│   ├── search_agent.py       # 搜索 Agent
│   ├── scrape_agent.py       # 抓取 Agent
│   ├── script_agent.py       # 脚本生成 Agent
│   ├── image_generate_agent.py  # 图像生成 Agent
│   └── audio_generate_agent.py  # 音频生成 Agent
│
├── services/                  # 业务逻辑层
│   ├── celery_tasks.py       # Celery 任务定义（核心！）
│   ├── async_podcast_agent_service.py
│   └── ...
│
├── tools/                     # Agent 工具
│   ├── web_search.py         # 网页搜索（Browser Use）
│   ├── embedding_search.py   # 向量搜索
│   ├── social/               # 社交媒体工具
│   │   ├── x_scraper.py     # X.com 抓取
│   │   └── fb_scraper.py    # Facebook 抓取
│   └── ...
│
├── processors/                # 后台处理器
│   ├── feed_processor.py     # RSS 订阅处理
│   ├── ai_analysis_processor.py  # AI 分析
│   ├── embedding_processor.py   # 向量化
│   └── podcast_generator_processor.py  # 播客生成
│
├── routers/                   # API 路由
│   ├── article_router.py
│   ├── podcast_router.py
│   └── ...
│
├── db/                        # 数据库模型
│   ├── agent_config_v2.py    # Agent 配置
│   ├── articles.py
│   └── ...
│
└── web/                       # React 前端
    └── src/
```

## 🚀 快速开始

### 1. 安装依赖

```bash
cd advanced_ai_agents/multi_agent_apps/ai_news_and_podcast_agents/miniblog
pip install -r requirements.txt
python -m playwright install
```

### 2. 配置环境变量

创建 `.env` 文件：

```bash
OPENAI_API_KEY=your_openai_api_key
ELEVENSLAB_API_KEY=your_elevenlabs_api_key  # 可选
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
```

### 3. 启动 Redis

```bash
# macOS
brew install redis
brew services start redis

# 或使用 Docker
docker run -d -p 6379:6379 redis
```

### 4. 初始化数据库（首次运行）

```bash
python main.py  # 首次运行会初始化数据库
```

### 5. 启动所有服务

需要**3个终端**同时运行：

```bash
# 终端 1: 主后端服务
python main.py

# 终端 2: Celery Worker（处理异步任务）
python celery_worker.py

# 终端 3: 定时任务调度器
python scheduler.py
```

### 6. （可选）启动前端

```bash
cd web
npm install
npm start
```

访问：`http://localhost:7000`

## 📚 核心学习点

### 1. **多 Agent 系统架构** ⭐⭐⭐⭐⭐

**文件：** `services/celery_tasks.py`

```python
# 主 Agent 配置了 10+ 个工具
_agent = Agent(
    tools=[
        search_agent_run,           # 搜索工具
        scrape_agent_run,           # 抓取工具
        podcast_script_agent_run,   # 脚本生成
        image_generation_agent_run, # 图像生成
        audio_generate_agent_run,   # 音频生成
        # ... 更多工具
    ],
    session_state=session_state,    # 会话状态管理
    add_history_to_messages=True,   # 历史记录
)
```

**学习点：**
- Agent 如何管理多个工具
- 会话状态如何持久化
- 工具之间的协作

### 2. **Celery 异步任务处理** ⭐⭐⭐⭐⭐

**文件：** `services/celery_tasks.py`, `celery_worker.py`

```python
@app.task(bind=True, max_retries=0, base=SessionLockedTask)
def agent_chat(self, session_id, message):
    # Agent 处理用户消息
    response = _agent.run(message, session_id=session_id)
    return response
```

**学习点：**
- 如何使用 Celery 处理长时间运行的 Agent 任务
- 任务队列管理
- 会话锁定机制

### 3. **浏览器自动化** ⭐⭐⭐⭐

**文件：** `tools/web_search.py`, `tools/social/x_scraper.py`

**学习点：**
- Browser Use 库的使用
- 社交媒体登录会话管理
- 浏览器自动化最佳实践

### 4. **向量搜索（FAISS）** ⭐⭐⭐⭐

**文件：** `tools/embedding_search.py`, `processors/embedding_processor.py`

**学习点：**
- 如何创建向量索引
- 语义搜索实现
- FAISS 库的使用

### 5. **RSS 和内容处理管道** ⭐⭐⭐⭐

**文件：** `processors/feed_processor.py`, `processors/url_processor.py`

**学习点：**
- RSS 订阅处理
- 内容抓取和分析
- 处理管道设计

### 6. **播客生成流程** ⭐⭐⭐⭐⭐

**完整流程：**
```
RSS/URL → 抓取 → AI 分析 → 向量化 → 搜索 → 
脚本生成 → 图像生成 → 音频生成 → 播客
```

**关键文件：**
- `processors/podcast_generator_processor.py`
- `agents/script_agent.py`
- `agents/image_generate_agent.py`
- `agents/audio_generate_agent.py`

### 7. **FastAPI + React 前后端分离** ⭐⭐⭐⭐

**后端：** FastAPI (`main.py`)
**前端：** React (`web/src/`)

**学习点：**
- RESTful API 设计
- 前后端分离架构
- 流式响应处理

## 🎓 学习路径

### 阶段 1：理解核心 Agent（1-2 天）

1. **阅读 `services/celery_tasks.py`**
   - 理解主 Agent 的配置
   - 理解工具如何注册和使用

2. **阅读 `agents/search_agent.py`**
   - 理解搜索 Agent 的实现
   - 理解结构化输出（Pydantic）

3. **阅读 `db/agent_config_v2.py`**
   - 理解 Agent 的指令和描述
   - 理解会话状态管理

### 阶段 2：理解工具系统（2-3 天）

4. **阅读 `tools/web_search.py`**
   - 理解浏览器自动化
   - 理解 Browser Use 集成

5. **阅读 `tools/embedding_search.py`**
   - 理解向量搜索
   - 理解 FAISS 使用

6. **阅读 `tools/social/x_scraper.py`**
   - 理解社交媒体抓取
   - 理解会话管理

### 阶段 3：理解处理管道（2-3 天）

7. **阅读 `processors/feed_processor.py`**
   - 理解 RSS 处理
   - 理解内容管道

8. **阅读 `processors/podcast_generator_processor.py`**
   - 理解播客生成流程
   - 理解多步骤处理

### 阶段 4：理解系统架构（3-5 天）

9. **阅读 `main.py`**
   - 理解 FastAPI 应用结构
   - 理解路由设计

10. **阅读 `scheduler.py`**
    - 理解定时任务
    - 理解任务调度

11. **阅读前端代码 `web/src/`**
    - 理解 React 组件
    - 理解 API 集成

## 🔑 关键文件详解

### 1. `services/celery_tasks.py` - 核心 Agent

**这是最重要的文件！** 定义了主 Agent 和所有工具。

**关键代码：**
```python
_agent = Agent(
    model=OpenAIChat(id=AGENT_MODEL),
    storage=SqliteStorage(...),  # 会话存储
    tools=[...],                  # 10+ 个工具
    session_state=session_state,  # 状态管理
    add_history_to_messages=True, # 历史记录
)
```

### 2. `agents/search_agent.py` - 搜索 Agent

**展示了如何创建专门的 Agent：**
- 使用结构化输出（Pydantic）
- 多个搜索工具集成
- 智能工具选择

### 3. `tools/web_search.py` - 浏览器搜索

**展示了浏览器自动化：**
- Browser Use 集成
- 浏览器会话管理
- 搜索结果处理

### 4. `processors/podcast_generator_processor.py` - 播客生成

**展示了完整的内容处理流程：**
- 多步骤处理
- Agent 调用
- 文件生成

## 💡 实践建议

### 1. 先运行起来
- 按照快速开始步骤运行
- 使用 `bootstrap_demo.py` 加载示例数据
- 在 Web UI 中体验功能

### 2. 阅读代码顺序
1. `main.py` - 了解整体架构
2. `services/celery_tasks.py` - 理解核心 Agent
3. `agents/search_agent.py` - 理解 Agent 实现
4. `tools/web_search.py` - 理解工具实现
5. `processors/` - 理解处理管道

### 3. 修改实验
- 尝试添加新的工具
- 修改 Agent 指令
- 创建新的处理器

### 4. 调试技巧
- 查看 Celery 日志
- 使用 `flower` 监控任务
- 检查数据库内容

## 🐛 常见问题

### Q: Redis 连接失败？
A: 确保 Redis 正在运行：`redis-cli ping`

### Q: Celery Worker 不工作？
A: 确保 Redis 可访问，检查 `celery_worker.py` 配置

### Q: 浏览器自动化失败？
A: 确保 Playwright 已安装：`python -m playwright install`

### Q: FAISS 安装失败？
A: 可以跳过，只在需要向量搜索时安装

## 📖 相关文档

- [MiniBlog 完整 README](../readme.md)
- [Agno 框架文档](https://github.com/agno-agi/agno)
- [Celery 文档](https://docs.celeryq.dev/)
- [Browser Use 文档](https://browser-use.com/)

## 🎯 学习目标

完成学习后，你应该能够：

1. ✅ 理解多 Agent 系统架构
2. ✅ 掌握 Celery 异步任务处理
3. ✅ 理解浏览器自动化集成
4. ✅ 掌握向量搜索实现
5. ✅ 理解内容处理管道设计
6. ✅ 理解前后端分离架构
7. ✅ 能够扩展和定制系统

---

**这是一个生产级系统，代码质量很高，非常适合深入学习！** 🚀

