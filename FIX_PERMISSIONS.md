# 🔐 修复 AWS IAM 权限错误

## 错误说明

你遇到的错误是：
```
User: arn:aws:iam::285604254102:user/miniblog is not authorized to perform: ecr:CreateRepository
```

**原因**：IAM 用户 `miniblog` 没有创建 ECR 仓库的权限。

## 解决方案

### 方案 1: 使用 AWS Console 添加权限（推荐）

1. **登录 AWS Console**
   - 访问：https://console.aws.amazon.com/
   - 使用有管理员权限的账户登录

2. **进入 IAM 控制台**
   - 搜索 "IAM" 或访问：https://console.aws.amazon.com/iam/

3. **找到用户 `miniblog`**
   - 点击左侧 "Users"
   - 搜索并点击用户 `miniblog`

4. **添加权限策略**
   - 点击 "Add permissions" 或 "Add inline policy"
   - 选择 "JSON" 标签
   - 复制 `infra/iam-policy.json` 的内容并粘贴
   - 点击 "Review policy"
   - 命名策略为：`MiniblogDeploymentPolicy`
   - 点击 "Create policy"

### 方案 2: 使用 AWS CLI 添加权限

如果你有管理员权限，可以运行：

```bash
# 创建策略
aws iam create-policy \
  --policy-name MiniblogDeploymentPolicy \
  --policy-document file://infra/iam-policy.json \
  --region us-east-1

# 获取策略 ARN（从输出中复制）
POLICY_ARN="arn:aws:iam::285604254102:policy/MiniblogDeploymentPolicy"

# 将策略附加到用户
aws iam attach-user-policy \
  --user-name miniblog \
  --policy-arn $POLICY_ARN
```

### 方案 3: 使用管理员账户（临时方案）

如果你有另一个有管理员权限的 AWS 账户，可以：

1. **切换到管理员账户**
   ```bash
   aws configure
   # 输入管理员账户的 Access Key
   ```

2. **创建 ECR 仓库（手动）**
   ```bash
   aws ecr create-repository \
     --repository-name miniblog-agent \
     --region us-east-1
   ```

3. **然后切换回 miniblog 用户继续部署**

## 最小权限策略（如果只需要 ECR）

如果只想快速解决 ECR 权限问题，可以使用这个最小策略：

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:CreateRepository",
        "ecr:DescribeRepositories",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "*"
    }
  ]
}
```

## 验证权限

添加权限后，验证是否生效：

```bash
# 测试 ECR 权限
aws ecr describe-repositories --region us-east-1

# 如果成功，说明权限已添加
```

## 完整部署所需权限

完整的部署需要以下 AWS 服务的权限：

- ✅ **ECR** - 容器镜像仓库
- ✅ **ECS** - 容器服务
- ✅ **Secrets Manager** - 密钥管理
- ✅ **CloudWatch Logs** - 日志
- ✅ **IAM** - 角色管理（PassRole）
- ✅ **VPC/Networking** - 网络配置（如果使用 Terraform）
- ✅ **ElastiCache** - Redis（如果使用 Terraform）
- ✅ **RDS** - 数据库（如果使用 Terraform）
- ✅ **S3** - 存储（如果使用 Terraform）

## 下一步

添加权限后，重新运行部署脚本：

```bash
./deploy-to-aws.sh
```

如果仍有权限问题，请告诉我具体的错误信息。
