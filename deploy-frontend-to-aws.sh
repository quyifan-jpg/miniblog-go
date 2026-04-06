#!/bin/bash

# 前端静态站点部署到 AWS S3
# 使用前：确保后端 API 已部署，并准备好后端地址（如 ECS 负载均衡器 URL）
# 使用：./deploy-frontend-to-aws.sh

set -e

export AWS_PAGER=""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "🌐 MiniBlog 前端部署到 S3"
echo "=========================================="
echo ""

# 检查 Node.js
if ! command -v node &>/dev/null; then
    echo -e "${RED}❌ 未安装 Node.js，请先安装 (如 brew install node)${NC}"
    exit 1
fi

# 检查 AWS CLI
if ! command -v aws &>/dev/null; then
    echo -e "${RED}❌ 未安装 AWS CLI${NC}"
    exit 1
fi

# 获取 AWS 账户与区域
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text 2>/dev/null)
if [ -z "$AWS_ACCOUNT_ID" ]; then
    echo -e "${RED}❌ 无法获取 AWS 账户，请先运行 aws configure${NC}"
    exit 1
fi

read -p "🌍 AWS 区域 [us-east-1]: " AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}

# 后端 API 地址（前端会通过 REACT_APP_API_URL 请求该地址）
read -p "🔗 后端 API 地址 (如 https://your-alb-xxx.us-east-1.elb.amazonaws.com 或 http://IP:8000): " BACKEND_URL
if [ -z "$BACKEND_URL" ]; then
    echo -e "${RED}❌ 必须填写后端 API 地址${NC}"
    exit 1
fi
# 去掉末尾斜杠
BACKEND_URL="${BACKEND_URL%/}"

# S3 桶名（全局唯一，建议带账户 ID 或随机后缀）
DEFAULT_BUCKET="miniblog-web-${AWS_ACCOUNT_ID}"
read -p "📦 S3 桶名称 [${DEFAULT_BUCKET}]: " S3_BUCKET
S3_BUCKET=${S3_BUCKET:-$DEFAULT_BUCKET}

echo ""
echo -e "${GREEN}✅ 区域: $AWS_REGION | 后端: $BACKEND_URL | 桶: $S3_BUCKET${NC}"
echo ""

# 构建前端
echo "=========================================="
echo "步骤 1: 构建前端"
echo "=========================================="
cd "$(dirname "$0")/web"
echo "REACT_APP_API_URL=$BACKEND_URL npm run build"
REACT_APP_API_URL="$BACKEND_URL" npm run build
cd ..
echo -e "${GREEN}✅ 构建完成${NC}"
echo ""

# 创建 S3 桶（若不存在）
echo "=========================================="
echo "步骤 2: 准备 S3 桶"
echo "=========================================="
if ! aws s3api head-bucket --bucket "$S3_BUCKET" --region "$AWS_REGION" 2>/dev/null; then
    echo "创建 S3 桶: $S3_BUCKET"
    if [ "$AWS_REGION" = "us-east-1" ]; then
        aws s3api create-bucket --bucket "$S3_BUCKET" --region "$AWS_REGION"
    else
        aws s3api create-bucket --bucket "$S3_BUCKET" --region "$AWS_REGION" --create-bucket-configuration "LocationConstraint=$AWS_REGION"
    fi
    echo -e "${YELLOW}⚠️  新桶默认禁止公开访问。若要用 S3 静态网站地址访问，请到控制台 S3 → 桶 $S3_BUCKET → 权限 → 阻止公开访问 → 关闭，并保存桶策略（见下方）${NC}"
else
    echo "桶已存在: $S3_BUCKET"
fi

# 设置静态网站托管
aws s3 website "s3://${S3_BUCKET}" --index-document index.html --error-document index.html 2>/dev/null || true

# 生成并应用桶策略（允许公开读，用于静态网站）
POLICY_FILE=$(mktemp)
cat > "$POLICY_FILE" << POL
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${S3_BUCKET}/*"
    }
  ]
}
POL
aws s3api put-bucket-policy --bucket "$S3_BUCKET" --policy "file://${POLICY_FILE}"
rm -f "$POLICY_FILE"
echo "已设置桶策略（允许公开读取对象）。若仍无法访问，请在控制台关闭该桶的「阻止全部公共访问」。"
echo ""

# 上传构建产物
echo "=========================================="
echo "步骤 3: 上传到 S3"
echo "=========================================="
aws s3 sync web/build/ "s3://${S3_BUCKET}/" --delete --region "$AWS_REGION"
echo -e "${GREEN}✅ 上传完成${NC}"
echo ""

echo "=========================================="
echo "✅ 前端部署完成"
echo "=========================================="
echo ""
echo "📌 S3 静态网站地址（需桶允许公开读）："
echo "   http://${S3_BUCKET}.s3-website-${AWS_REGION}.amazonaws.com"
echo ""
echo "若未配置桶策略允许公开读，请在 IAM/桶策略中允许 s3:GetObject 公开访问，"
echo "或使用 CloudFront 分发并设置 S3 为源（推荐生产环境）。"
echo ""
