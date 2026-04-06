# AWS 部署指南 - 分步教程

## 📋 部署前检查清单

### 1. 必需工具
- [ ] AWS 账户
- [ ] AWS CLI 已安装并配置
- [ ] Docker 已安装并运行
- [ ] 已准备好环境变量（API Keys）

### 2. AWS 资源准备
- [ ] 选择 AWS 区域（推荐：us-east-1）
- [ ] 准备 Redis 密码
- [ ] 准备数据库密码
- [ ] OpenAI API Key

---

## 🚀 部署步骤

### 步骤 1: 安装和配置 AWS CLI

```bash
# 检查是否已安装
aws --version

# 如果未安装，macOS 使用 Homebrew：
brew install awscli

# 配置 AWS 凭证
aws configure
# 输入：
# - AWS Access Key ID
# - AWS Secret Access Key
# - Default region name (例如: us-east-1)
# - Default output format (json)
```

### 步骤 2: 获取 AWS Account ID

```bash
aws sts get-caller-identity --query Account --output text
```

保存这个 Account ID，后续会用到。

### 步骤 3: 创建 ECR 仓库

```bash
# 设置变量
export AWS_REGION=us-east-1
export ECR_REPO_NAME=miniblog-agent
export AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# 创建 ECR 仓库
aws ecr create-repository \
  --repository-name $ECR_REPO_NAME \
  --region $AWS_REGION \
  --image-scanning-configuration scanOnPush=true

# 登录 ECR
aws ecr get-login-password --region $AWS_REGION | \
  docker login --username AWS --password-stdin \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
```

### 步骤 4: 构建和推送 Docker 镜像

```bash
cd agent

# 使用部署脚本（推荐）
./deploy-aws.sh $AWS_REGION $ECR_REPO_NAME $AWS_ACCOUNT_ID

# 或手动执行：
docker build -t $ECR_REPO_NAME:latest .
docker tag $ECR_REPO_NAME:latest \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:latest
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:latest
```

### 步骤 5: 创建 AWS Secrets Manager 密钥

```bash
# 创建 OpenAI API Key 密钥
aws secretsmanager create-secret \
  --name miniblog/openai-api-key \
  --secret-string "your-openai-api-key-here" \
  --region $AWS_REGION

# 创建 Redis 密码密钥
aws secretsmanager create-secret \
  --name miniblog/redis-password \
  --secret-string "your-redis-password-here" \
  --region $AWS_REGION

# 创建数据库密码密钥
aws secretsmanager create-secret \
  --name miniblog/db-password \
  --secret-string "your-db-password-here" \
  --region $AWS_REGION
```

### 步骤 6: 使用 Terraform 创建基础设施（推荐）

```bash
cd infra/terraform

# 创建 terraform.tfvars
cat > terraform.tfvars << EOF
aws_region = "$AWS_REGION"
vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]
redis_password = "your-redis-password"
db_username = "miniblog"
db_password = "your-db-password"
EOF

# 初始化 Terraform
terraform init

# 查看计划
terraform plan

# 应用配置（创建资源）
terraform apply
```

### 步骤 7: 创建 ECS 任务定义和服务

参考 `docs/AWS_DEPLOYMENT1.md` 中的详细配置。

---

## 🎯 快速部署脚本

我为你准备了一个自动化部署脚本，可以一键完成大部分工作。
