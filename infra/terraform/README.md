# Terraform 基础设施配置

这个目录包含用于在 AWS 上部署 MiniBlog 的 Terraform 配置。

## 使用方法

### 1. 初始化 Terraform

```bash
cd infra/terraform
terraform init
```

### 2. 创建 terraform.tfvars 文件

```bash
cat > terraform.tfvars << EOF
aws_region = "us-east-1"
vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]
redis_password = "your-redis-password"
db_username = "miniblog"
db_password = "your-db-password"
EOF
```

### 3. 规划部署

```bash
terraform plan
```

### 4. 应用配置

```bash
terraform apply
```

### 5. 查看输出

```bash
terraform output
```

## 模块结构

- `main.tf`: 主要资源定义（VPC、ECS 集群、ECR、S3）
- `variables.tf`: 变量定义
- `elasticache.tf`: ElastiCache Redis 配置（需要单独创建）
- `rds.tf`: RDS PostgreSQL 配置（需要单独创建）
- `ecs.tf`: ECS 任务定义和服务配置（需要单独创建）
- `scheduled-tasks.tf`: ECS Scheduled Tasks 配置（需要单独创建）

## 注意事项

1. 确保 AWS 凭证已配置（`aws configure` 或环境变量）
2. 根据实际需求调整资源大小和数量
3. 生产环境建议使用 S3 backend 存储 Terraform state
4. 敏感信息（密码、API 密钥）应使用 AWS Secrets Manager
