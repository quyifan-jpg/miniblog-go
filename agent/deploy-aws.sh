#!/bin/bash

# AWS 部署脚本
# 使用方法: ./deploy-aws.sh <aws-region> <ecr-repo-name> <account-id>

set -e

REGION=${1:-us-east-1}
REPO_NAME=${2:-miniblog-agent}
ACCOUNT_ID=${3}

if [ -z "$ACCOUNT_ID" ]; then
    echo "错误: 请提供 AWS Account ID"
    echo "使用方法: ./deploy-aws.sh <region> <repo-name> <account-id>"
    exit 1
fi

ECR_URI="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO_NAME}"

echo "=========================================="
echo "开始部署到 AWS ECR"
echo "Region: $REGION"
echo "Repository: $REPO_NAME"
echo "ECR URI: $ECR_URI"
echo "=========================================="

# 1. 登录 ECR
echo "步骤 1: 登录 ECR..."
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_URI

# 2. 检查仓库是否存在，不存在则创建
echo "步骤 2: 检查 ECR 仓库..."
if ! aws ecr describe-repositories --repository-names $REPO_NAME --region $REGION &>/dev/null; then
    echo "创建 ECR 仓库..."
    aws ecr create-repository --repository-name $REPO_NAME --region $REGION
else
    echo "ECR 仓库已存在"
fi

# 3. 构建 Docker 镜像
echo "步骤 3: 构建 Docker 镜像..."
docker build -t $REPO_NAME:latest .

# 4. 标记镜像
echo "步骤 4: 标记镜像..."
docker tag $REPO_NAME:latest $ECR_URI:latest
docker tag $REPO_NAME:latest $ECR_URI:$(date +%Y%m%d-%H%M%S)

# 5. 推送镜像
echo "步骤 5: 推送镜像到 ECR..."
docker push $ECR_URI:latest
docker push $ECR_URI:$(date +%Y%m%d-%H%M%S)

echo "=========================================="
echo "部署完成！"
echo "镜像 URI: $ECR_URI:latest"
echo "=========================================="
echo ""
echo "下一步:"
echo "1. 更新 ECS 任务定义中的镜像 URI"
echo "2. 创建新的任务定义修订版"
echo "3. 更新 ECS 服务以使用新任务定义"
