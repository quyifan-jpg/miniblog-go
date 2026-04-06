FROM python:3.11-slim

WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    make \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY requirements.txt .

# 安装 Python 依赖
RUN pip install --no-cache-dir -r requirements.txt

# 复制应用代码
COPY . .

# 创建必要的目录
RUN mkdir -p databases podcasts/audio podcasts/images podcasts/recordings browsers

# 暴露端口
EXPOSE 8000

# 默认命令（可在 ECS Task Definition 或 docker-compose 中覆盖）
CMD ["python", "main.py"]
