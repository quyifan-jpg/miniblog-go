# 🚀 AWS 部署快速开始

## 前置要求

1. **AWS 账户** - 确保有 AWS 账户并已配置凭证
2. **AWS CLI** - 安装并配置
3. **Docker** - 确保 Docker Desktop 正在运行

## 一键部署

### 1. 安装 AWS CLI（如果未安装）

```bash
# macOS
brew install awscli

# 配置 AWS 凭证
aws configure
# 输入你的 AWS Access Key ID, Secret Access Key, 区域等
```

### 2. 运行部署脚本

```bash
# 在项目根目录
./deploy-to-aws.sh
```

脚本会自动：
- ✅ 检查环境（AWS CLI, Docker）
- ✅ 创建 ECR 仓库
- ✅ 构建 Docker 镜像
- ✅ 推送镜像到 ECR

### 3. 创建 Secrets Manager 密钥

```bash
# 设置变量
export AWS_REGION=us-east-1  # 或你选择的区域

# 创建 OpenAI API Key
aws secretsmanager create-secret \
  --name miniblog/openai-api-key \
  --secret-string "sk-your-openai-api-key" \
  --region $AWS_REGION

# 创建 Redis 密码
aws secretsmanager create-secret \
  --name miniblog/redis-password \
  --secret-string "your-secure-redis-password" \
  --region $AWS_REGION

# 创建数据库密码
aws secretsmanager create-secret \
  --name miniblog/db-password \
  --secret-string "your-secure-db-password" \
  --region $AWS_REGION
```

### 4. 使用 Terraform 创建基础设施

```bash
cd infra/terraform

# 创建配置文件
cat > terraform.tfvars << EOF
aws_region = "us-east-1"
vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]
redis_password = "your-redis-password"
db_username = "miniblog"
db_password = "your-db-password"
EOF

# 初始化并部署
terraform init
terraform plan
terraform apply
```

### 5. 创建 ECS 任务定义和服务

参考 `docs/AWS_DEPLOYMENT1.md` 中的详细步骤创建：
- ECS 任务定义（FastAPI、Celery、Scheduler）
- ECS 服务
- ECS Scheduled Tasks

## 验证部署

部署完成后，检查：

```bash
# 检查 ECS 服务状态
aws ecs list-services --cluster miniblog-cluster --region us-east-1

# 检查任务运行状态
aws ecs list-tasks --cluster miniblog-cluster --region us-east-1

# 查看日志
aws logs tail /ecs/miniblog --follow --region us-east-1
```

## 需要帮助？

查看详细文档：
- `docs/AWS_DEPLOYMENT1.md` - 完整部署指南
- `docs/REDIS_USAGE.md` - Redis 使用说明
- `DEPLOY_GUIDE.md` - 分步部署教程
