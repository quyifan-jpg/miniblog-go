#!/bin/bash

# ============================================================
# MiniBlog 本地一键启动脚本
# 用法: ./start.sh [local|docker|stop]
# ============================================================

set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
AGENT_DIR="$ROOT_DIR/agent"
WEB_DIR="$ROOT_DIR/web"
LOG_DIR="$ROOT_DIR/.logs"
PID_FILE="$ROOT_DIR/.pids"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[!]${NC} $1"; }
log_error() { echo -e "${RED}[✗]${NC} $1"; }
log_step()  { echo -e "${CYAN}[→]${NC} $1"; }

print_banner() {
  echo -e "${BLUE}"
  echo "  ╔══════════════════════════════════╗"
  echo "  ║        MiniBlog Launcher         ║"
  echo "  ╚══════════════════════════════════╝"
  echo -e "${NC}"
}

# ─────────────────────────────────────────
# STOP: 停止所有进程
# ─────────────────────────────────────────
stop_all() {
  log_step "正在停止所有服务..."

  if [ -f "$PID_FILE" ]; then
    while IFS= read -r pid; do
      if kill -0 "$pid" 2>/dev/null; then
        kill "$pid" && log_info "已停止进程 $pid"
      fi
    done < "$PID_FILE"
    rm -f "$PID_FILE"
  else
    log_warn "没有找到运行中的进程记录"
  fi

  # 兜底：杀掉匹配的进程
  pkill -f "python main.py"      2>/dev/null && log_info "已停止 main.py"      || true
  pkill -f "python scheduler.py" 2>/dev/null && log_info "已停止 scheduler.py" || true
  pkill -f "celery_worker.py"    2>/dev/null && log_info "已停止 celery_worker" || true

  log_info "所有服务已停止"
}

# ─────────────────────────────────────────
# DOCKER 模式
# ─────────────────────────────────────────
start_docker() {
  log_step "检查 Docker..."
  if ! command -v docker &>/dev/null; then
    log_error "Docker 未安装，请先安装 Docker Desktop"
    exit 1
  fi
  if ! docker info &>/dev/null; then
    log_error "Docker 未运行，请启动 Docker Desktop"
    exit 1
  fi

  # 检查 .env
  if [ ! -f "$AGENT_DIR/.env" ]; then
    log_warn "未找到 $AGENT_DIR/.env，正在从模板创建..."
    cat > "$AGENT_DIR/.env" <<EOF
OPENAI_API_KEY=your_openai_api_key_here
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=root
REDIS_DB=0
PORT=8000
EOF
    log_warn "请编辑 agent/.env 填入你的 API Key，然后重新运行"
    open "$AGENT_DIR/.env" 2>/dev/null || nano "$AGENT_DIR/.env"
    exit 0
  fi

  log_step "启动 Docker 服务..."
  cd "$AGENT_DIR"
  docker compose up --build -d

  log_info "所有服务已启动（Docker 模式）"
  echo ""
  echo -e "  ${GREEN}后端 API：${NC}  http://localhost:8000"
  echo -e "  ${GREEN}API 文档：${NC}  http://localhost:8000/docs"
  echo ""
  echo -e "  ${YELLOW}查看日志：${NC}  docker compose logs -f   (在 agent/ 目录下)"
  echo -e "  ${YELLOW}停止服务：${NC}  docker compose down       (在 agent/ 目录下)"

  sleep 2
  open "http://localhost:8000" 2>/dev/null || true
}

# ─────────────────────────────────────────
# LOCAL 模式
# ─────────────────────────────────────────
start_local() {
  mkdir -p "$LOG_DIR"
  > "$PID_FILE"

  # 1. 检查 Redis
  log_step "检查 Redis..."
  if ! command -v redis-cli &>/dev/null; then
    log_error "Redis 未安装，请运行: brew install redis"
    exit 1
  fi
  if ! redis-cli ping &>/dev/null; then
    log_warn "Redis 未运行，正在启动..."
    brew services start redis
    sleep 1
  fi
  log_info "Redis 正在运行"

  # 2. 找 Python 虚拟环境
  log_step "查找 Python 虚拟环境..."

  VENV_ACTIVATE=""
  for VENV_PATH in "$AGENT_DIR/venv" "$ROOT_DIR/venv" "$ROOT_DIR/.venv"; do
    if [ -f "$VENV_PATH/bin/activate" ]; then
      VENV_ACTIVATE="$VENV_PATH/bin/activate"
      log_info "使用虚拟环境: $VENV_PATH"
      break
    fi
  done

  if [ -z "$VENV_ACTIVATE" ]; then
    log_warn "未找到虚拟环境，正在创建..."
    python3 -m venv "$AGENT_DIR/venv"
    VENV_ACTIVATE="$AGENT_DIR/venv/bin/activate"
    source "$VENV_ACTIVATE"
    log_step "安装依赖（首次需要几分钟）..."
    pip install -r "$AGENT_DIR/requirements.txt" -q
    log_step "安装 Playwright 浏览器驱动..."
    playwright install chromium 2>/dev/null || true
  else
    source "$VENV_ACTIVATE"
  fi

  # 3. 检查 .env
  if [ ! -f "$AGENT_DIR/.env" ]; then
    log_warn "未找到 $AGENT_DIR/.env，正在创建模板..."
    cat > "$AGENT_DIR/.env" <<EOF
OPENAI_API_KEY=your_openai_api_key_here
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
PORT=8000
EOF
    log_warn "请编辑 agent/.env 填入你的 API Key，然后重新运行"
    ${EDITOR:-nano} "$AGENT_DIR/.env"
    exit 0
  fi

  # 检查 API Key 是否还是占位符
  if grep -q "your_openai_api_key_here" "$AGENT_DIR/.env"; then
    log_warn "⚠️  agent/.env 中的 OPENAI_API_KEY 还未配置！"
    log_warn "    请编辑 agent/.env 并重新运行"
    exit 1
  fi

  # 4. 启动主后端
  log_step "启动后端 (main.py)..."
  cd "$AGENT_DIR"
  source "$VENV_ACTIVATE"
  nohup python main.py > "$LOG_DIR/main.log" 2>&1 &
  echo $! >> "$PID_FILE"
  log_info "后端已启动 → 日志: .logs/main.log"

  # 等待后端就绪
  log_step "等待后端就绪..."
  for i in {1..20}; do
    if curl -s "http://localhost:8000/health" &>/dev/null || \
       curl -s "http://localhost:8000/api/articles/" &>/dev/null; then
      log_info "后端已就绪"
      break
    fi
    sleep 1
    printf "."
  done
  echo ""

  # 5. 启动调度器
  log_step "启动调度器 (scheduler.py)..."
  nohup python "$AGENT_DIR/scheduler.py" > "$LOG_DIR/scheduler.log" 2>&1 &
  echo $! >> "$PID_FILE"
  log_info "调度器已启动 → 日志: .logs/scheduler.log"

  # 6. 启动 Celery Worker
  log_step "启动 Celery Worker..."
  nohup python "$AGENT_DIR/celery_worker.py" > "$LOG_DIR/celery.log" 2>&1 &
  echo $! >> "$PID_FILE"
  log_info "Celery Worker 已启动 → 日志: .logs/celery.log"

  # 7. 完成
  echo ""
  echo -e "  ${GREEN}════════════════════════════════════${NC}"
  echo -e "  ${GREEN}  🚀 MiniBlog 已启动！${NC}"
  echo -e "  ${GREEN}════════════════════════════════════${NC}"
  echo ""
  echo -e "  ${CYAN}后端 API：${NC}  http://localhost:8000"
  echo -e "  ${CYAN}API 文档：${NC}  http://localhost:8000/docs"
  echo ""
  echo -e "  ${YELLOW}查看日志：${NC}"
  echo -e "    tail -f .logs/main.log"
  echo -e "    tail -f .logs/scheduler.log"
  echo -e "    tail -f .logs/celery.log"
  echo ""
  echo -e "  ${YELLOW}停止所有服务：${NC}  ./start.sh stop"
  echo ""

  open "http://localhost:8000" 2>/dev/null || true
}

# ─────────────────────────────────────────
# 入口
# ─────────────────────────────────────────
print_banner

MODE="${1:-local}"

case "$MODE" in
  local)   start_local  ;;
  docker)  start_docker ;;
  stop)    stop_all     ;;
  *)
    echo "用法: $0 [local|docker|stop]"
    echo ""
    echo "  local   本地直接启动（默认）"
    echo "  docker  使用 Docker Compose 启动"
    echo "  stop    停止所有本地服务"
    exit 1
    ;;
esac
