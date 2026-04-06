#!/bin/bash

# AWS 一键部署脚本
# 使用方法: ./deploy-to-aws.sh

set -e

# 禁止 AWS CLI 使用分页器，避免脚本在 (END) 处卡住
export AWS_PAGER=""

echo "=========================================="
echo "🚀 MiniBlog AWS 部署向导"
echo "=========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 AWS CLI
if ! command -v aws &> /dev/null; then
    echo -e "${RED}❌ AWS CLI 未安装${NC}"
    echo "请先安装: brew install awscli"
    exit 1
fi

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装${NC}"
    exit 1
fi

# 检查 Docker 是否运行
if ! docker info &> /dev/null; then
    echo -e "${RED}❌ Docker 未运行，请先启动 Docker${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 环境检查通过${NC}"
echo ""

# 获取 AWS 账户信息
echo "📋 获取 AWS 账户信息..."
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text 2>/dev/null)
if [ -z "$AWS_ACCOUNT_ID" ]; then
    echo -e "${RED}❌ 无法获取 AWS 账户信息，请先运行: aws configure${NC}"
    exit 1
fi

echo -e "${GREEN}✅ AWS Account ID: $AWS_ACCOUNT_ID${NC}"

# 询问区域
read -p "🌍 请输入 AWS 区域 [us-east-1]: " AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}
echo -e "${GREEN}✅ 使用区域: $AWS_REGION${NC}"

# 询问 ECR 仓库名
read -p "📦 请输入 ECR 仓库名称 [miniblog-agent]: " ECR_REPO_NAME
ECR_REPO_NAME=${ECR_REPO_NAME:-miniblog-agent}
echo -e "${GREEN}✅ ECR 仓库: $ECR_REPO_NAME${NC}"

echo ""
echo "=========================================="
echo "步骤 1: 创建 ECR 仓库"
echo "=========================================="

# 检查仓库是否存在
if aws ecr describe-repositories --repository-names $ECR_REPO_NAME --region $AWS_REGION &>/dev/null; then
    echo -e "${YELLOW}⚠️  ECR 仓库已存在，跳过创建${NC}"
else
    echo "创建 ECR 仓库..."
    aws ecr create-repository --repository-name "$ECR_REPO_NAME" --region "$AWS_REGION" --image-scanning-configuration scanOnPush=true
    echo -e "${GREEN}✅ ECR 仓库创建成功${NC}"
fi

echo ""
echo "=========================================="
echo "步骤 2: 登录 ECR"
echo "=========================================="

echo "登录 ECR..."
aws ecr get-login-password --region $AWS_REGION | \
    docker login --username AWS --password-stdin \
    $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
echo -e "${GREEN}✅ ECR 登录成功${NC}"

echo ""
echo "=========================================="
echo "步骤 3: 构建 Docker 镜像"
echo "=========================================="

cd agent

echo "构建 Docker 镜像..."
docker build -t $ECR_REPO_NAME:latest .
echo -e "${GREEN}✅ 镜像构建成功${NC}"

echo ""
echo "=========================================="
echo "步骤 4: 标记并推送镜像"
echo "=========================================="

ECR_URI="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME"
IMAGE_TAG=$(date +%Y%m%d-%H%M%S)

echo "标记镜像..."
docker tag $ECR_REPO_NAME:latest $ECR_URI:latest
docker tag $ECR_REPO_NAME:latest $ECR_URI:$IMAGE_TAG

echo "推送镜像到 ECR..."
docker push $ECR_URI:latest
docker push $ECR_URI:$IMAGE_TAG

echo -e "${GREEN}✅ 镜像推送成功${NC}"

cd ..

echo ""
echo "=========================================="
echo "✅ 部署完成！"
echo "=========================================="
echo ""
echo "📝 下一步操作："
echo ""
echo "1. 创建 AWS Secrets Manager 密钥："
echo "   aws secretsmanager create-secret \\"
echo "     --name miniblog/openai-api-key \\"
echo "     --secret-string 'your-openai-api-key' \\"
echo "     --region $AWS_REGION"
echo ""
echo "2. 使用 Terraform 创建基础设施："
echo "   cd infra/terraform"
echo "   terraform init"
echo "   terraform plan"
echo "   terraform apply"
echo ""
echo "3. 创建 ECS 任务定义和服务（参考 docs/AWS_DEPLOYMENT1.md）"
echo ""
echo "📦 镜像 URI: $ECR_URI:latest"
echo ""
