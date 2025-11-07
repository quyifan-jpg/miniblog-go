#!/bin/bash

# Simple IM API 测试脚本

BASE_URL="http://localhost:8080"

echo "========================================="
echo "Simple IM API 测试脚本"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local token=$5
    
    echo -e "${YELLOW}测试: ${name}${NC}"
    
    if [ -z "$token" ]; then
        response=$(curl -s -X $method "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -X $method "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data")
    fi
    
    echo "响应: $response"
    echo ""
    
    echo "$response"
}

# 1. 测试健康检查
echo -e "${GREEN}=== 1. 健康检查 ===${NC}"
curl -s "${BASE_URL}/health" | jq '.'
echo ""

# 2. 注册用户 Alice
echo -e "${GREEN}=== 2. 注册用户 Alice ===${NC}"
register_alice=$(test_api "注册 Alice" POST "/api/auth/register" \
'{
  "username": "alice",
  "password": "password123",
  "nickname": "Alice",
  "email": "alice@example.com"
}')
echo "$register_alice" | jq '.'

# 3. 注册用户 Bob
echo -e "${GREEN}=== 3. 注册用户 Bob ===${NC}"
register_bob=$(test_api "注册 Bob" POST "/api/auth/register" \
'{
  "username": "bob",
  "password": "password123",
  "nickname": "Bob",
  "email": "bob@example.com"
}')
echo "$register_bob" | jq '.'

# 4. Alice 登录
echo -e "${GREEN}=== 4. Alice 登录 ===${NC}"
login_alice=$(test_api "Alice 登录" POST "/api/auth/login" \
'{
  "username": "alice",
  "password": "password123"
}')
echo "$login_alice" | jq '.'

alice_token=$(echo "$login_alice" | jq -r '.data.token')
alice_id=$(echo "$login_alice" | jq -r '.data.user.id')
echo "Alice Token: $alice_token"
echo "Alice ID: $alice_id"
echo ""

# 5. Bob 登录
echo -e "${GREEN}=== 5. Bob 登录 ===${NC}"
login_bob=$(test_api "Bob 登录" POST "/api/auth/login" \
'{
  "username": "bob",
  "password": "password123"
}')
echo "$login_bob" | jq '.'

bob_token=$(echo "$login_bob" | jq -r '.data.token')
bob_id=$(echo "$login_bob" | jq -r '.data.user.id')
echo "Bob Token: $bob_token"
echo "Bob ID: $bob_id"
echo ""

# 6. Alice 获取个人信息
echo -e "${GREEN}=== 6. Alice 获取个人信息 ===${NC}"
curl -s -X GET "${BASE_URL}/api/user/profile" \
    -H "Authorization: Bearer $alice_token" | jq '.'
echo ""

# 7. Alice 搜索 Bob
echo -e "${GREEN}=== 7. Alice 搜索 Bob ===${NC}"
curl -s -X GET "${BASE_URL}/api/user/search?keyword=bob" \
    -H "Authorization: Bearer $alice_token" | jq '.'
echo ""

# 8. Alice 添加 Bob 为好友
echo -e "${GREEN}=== 8. Alice 添加 Bob 为好友 ===${NC}"
curl -s -X POST "${BASE_URL}/api/friend/add" \
    -H "Authorization: Bearer $alice_token" \
    -H "Content-Type: application/json" \
    -d "{
        \"friend_id\": $bob_id,
        \"remark\": \"我的朋友 Bob\"
    }" | jq '.'
echo ""

# 9. Bob 查看好友请求
echo -e "${GREEN}=== 9. Bob 查看好友请求 ===${NC}"
curl -s -X GET "${BASE_URL}/api/friend/requests" \
    -H "Authorization: Bearer $bob_token" | jq '.'
echo ""

# 10. Bob 接受好友请求
echo -e "${GREEN}=== 10. Bob 接受好友请求 ===${NC}"
curl -s -X POST "${BASE_URL}/api/friend/handle" \
    -H "Authorization: Bearer $bob_token" \
    -H "Content-Type: application/json" \
    -d "{
        \"friend_id\": $alice_id,
        \"accept\": true
    }" | jq '.'
echo ""

# 11. Alice 查看好友列表
echo -e "${GREEN}=== 11. Alice 查看好友列表 ===${NC}"
curl -s -X GET "${BASE_URL}/api/friend/list" \
    -H "Authorization: Bearer $alice_token" | jq '.'
echo ""

# 12. Alice 发送消息给 Bob
echo -e "${GREEN}=== 12. Alice 发送消息给 Bob ===${NC}"
curl -s -X POST "${BASE_URL}/api/message/send" \
    -H "Authorization: Bearer $alice_token" \
    -H "Content-Type: application/json" \
    -d "{
        \"to_user_id\": $bob_id,
        \"msg_type\": 1,
        \"content\": \"Hi Bob! 你好！\"
    }" | jq '.'
echo ""

# 13. Bob 查看未读消息
echo -e "${GREEN}=== 13. Bob 查看未读消息 ===${NC}"
curl -s -X GET "${BASE_URL}/api/message/unread" \
    -H "Authorization: Bearer $bob_token" | jq '.'
echo ""

# 14. Bob 查看与 Alice 的聊天历史
echo -e "${GREEN}=== 14. Bob 查看与 Alice 的聊天历史 ===${NC}"
curl -s -X GET "${BASE_URL}/api/message/history?friend_id=$alice_id&offset=0&limit=20" \
    -H "Authorization: Bearer $bob_token" | jq '.'
echo ""

# 15. Alice 创建群组
echo -e "${GREEN}=== 15. Alice 创建群组 ===${NC}"
create_group=$(curl -s -X POST "${BASE_URL}/api/group/create" \
    -H "Authorization: Bearer $alice_token" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"测试群组\",
        \"description\": \"这是一个测试群组\",
        \"member_ids\": [$bob_id]
    }")
echo "$create_group" | jq '.'

group_id=$(echo "$create_group" | jq -r '.data.id')
echo "Group ID: $group_id"
echo ""

# 16. Alice 查看群组列表
echo -e "${GREEN}=== 16. Alice 查看群组列表 ===${NC}"
curl -s -X GET "${BASE_URL}/api/group/list" \
    -H "Authorization: Bearer $alice_token" | jq '.'
echo ""

# 17. 查看群组成员
echo -e "${GREEN}=== 17. 查看群组成员 ===${NC}"
curl -s -X GET "${BASE_URL}/api/group/${group_id}/members" \
    -H "Authorization: Bearer $alice_token" | jq '.'
echo ""

# 18. Alice 发送群消息
echo -e "${GREEN}=== 18. Alice 发送群消息 ===${NC}"
curl -s -X POST "${BASE_URL}/api/message/send" \
    -H "Authorization: Bearer $alice_token" \
    -H "Content-Type: application/json" \
    -d "{
        \"to_user_id\": 0,
        \"group_id\": $group_id,
        \"msg_type\": 1,
        \"content\": \"大家好！这是群消息。\"
    }" | jq '.'
echo ""

echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo -e "${GREEN}=========================================${NC}"

