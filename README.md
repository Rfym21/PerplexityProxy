# Pplx2Api

> 基于原 [原作者yushangxiao](https://github.com/yushangxiao) 的版本基础上添加代理池功能

## ✨ 特性
- 🖼️ **图像识别** - 发送图像给Ai进行分析
- 📝 **隐私模式** - 对话不保存在官网，可选择关闭
- 🌊 **流式响应** - 获取实时流式输出
- 📁 **文件上传支持** - 上传长文本内容
- 🧠 **思考过程** - 访问思考模型的逐步推理，自动输出`<think>`标签
- 🔄 **聊天历史管理** - 控制对话上下文长度，超出将上传为文件
- 🌐 **代理支持** - 通过您首选的代理路由请求
- 🔄 **代理池管理** - 自动从代理池API获取和轮换代理，提高稳定性
- 🔐 **API密钥认证** - 保护您的API端点
- 🔍 **搜索模式**- 访问 -search 结尾的模型，连接网络且返回搜索内容
- 📊 **模型监控** - 跟踪响应的实际模型，如果模型不一致会返回实际使用的模型
- 🔄 **自动刷新** 每天自动刷新cookie，持续可用
- 🖼️ **绘图模型** - 在搜索模式，支持模型绘图，文生图，图生图
 ## 📋 前提条件
 - Go 1.23+（从源代码构建）
 - Docker（用于容器化部署）

## ✨ 关于环境变量SESSIONS
  为 https://www.perplexity.ai 官网cookie中 __Secure-next-auth.session-token 的值
  
  环境变量SESSIONS可以设置多个账户轮询或重试，使用英文逗号分割即可


## 项目效果

 识图：
 
![image](https://github.com/user-attachments/assets/3bb823e0-4232-4c6c-93cd-76d6c329ede3)

搜索：

![image](https://github.com/user-attachments/assets/26f7b6f7-ef00-499b-be32-c5dbc6e80ea6)

思考：

![image](https://github.com/user-attachments/assets/a075584a-ab49-4bf9-857b-6436b34bd363)

模型检测：

![image](https://github.com/user-attachments/assets/06013dd7-31ff-4bdd-bc5a-746ecaa8e922)

文生图：

![image](https://github.com/user-attachments/assets/bae2fd09-c738-4078-81a3-993c0b805943)

图生图：

![image](https://github.com/user-attachments/assets/f1866af5-5558-4fbb-83d7-b753035628bd)


 ## 🚀 部署选项
 
 ### Docker
 ```bash
 docker run -d \
   -p 8080:8080 \
   -e SESSIONS=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0**,eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0** \
   -e APIKEY=sk-123456 \
   -e IS_INCOGNITO=true \
   -e MAX_CHAT_HISTORY_LENGTH=10000 \
   -e NO_ROLE_PREFIX=false \
   -e SEARCH_RESULT_COMPATIBLE=false \
   -e ENABLE_PROXY_POOL=false \
   --name pplx-proxy \
   ghcr.io/rfym21/pplx-proxy:latest
 ```
 
 ### Docker Compose
 创建一个`docker-compose.yml`文件：
 ```yaml
 version: '3'
 services:
   pplx-proxy:
     image: ghcr.io/rfym21/pplx-proxy:latest
     container_name: pplx-proxy
     ports:
       - "8080:8080"
     environment:
       - SESSIONS=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0**,eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0**
       - ADDRESS=0.0.0.0:8080
       - APIKEY=sk-123456
       - PROXY=http://proxy:2080
       - ENABLE_PROXY_POOL=false
       - PROXY_POOL_API=https://your-proxy-api.com/get
       - PROXY_POOL_SIZE=10
       - MAX_CHAT_HISTORY_LENGTH=10000
       - NO_ROLE_PREFIX=false
       - IS_INCOGNITO=true
       - SEARCH_RESULT_COMPATIBLE=false
     restart: always
 ```
 然后运行：
 ```bash
 docker-compose up -d
 ```

 ### 启用代理池的 Docker 示例
 ```bash
 docker run -d \
   -p 8080:8080 \
   -e SESSIONS=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0** \
   -e APIKEY=sk-123456 \
   -e ENABLE_PROXY_POOL=true \
   -e PROXY_POOL_API=https://your-proxy-api.com/get \
   -e PROXY_POOL_SIZE=20 \
   -e IS_INCOGNITO=true \
   --name pplx-proxy-with-pool \
   ghcr.io/rfym21/pplx-proxy:latest
 ```
 
 ## ⚙️ 配置

### � 配置文件
项目提供了 `.env.example` 配置模板文件，包含所有可用的环境变量及其说明。

**使用步骤**：
1. 复制配置模板：`cp .env.example .env`
2. 编辑 `.env` 文件，填入您的配置值
3. 重启服务使配置生效

### �🔐 基础认证配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `SESSIONS` | Perplexity 会话令牌，支持多个账户（英文逗号分隔） | - | ✅ |
| `APIKEY` | API 访问密钥，用于保护端点安全 | - | ✅ |

### 🌐 服务器配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `ADDRESS` | 服务器监听地址和端口 | `0.0.0.0:8080` | ❌ |

### 🔄 代理配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `PROXY` | 单一代理地址 (http://user:pass@host:port) | - | ❌ |
| `ENABLE_PROXY_POOL` | 启用代理池功能 | `false` | ❌ |
| `PROXY_POOL_API` | 代理池 API 地址 | - | ❌ |
| `PROXY_POOL_SIZE` | 代理池大小（建议 10-50） | `10` | ❌ |

### 💬 聊天功能配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `IS_INCOGNITO` | 隐私模式，对话不保存在官网 | `true` | ❌ |
| `MAX_CHAT_HISTORY_LENGTH` | 聊天历史长度限制，超出转为文件 | `10000` | ❌ |
| `NO_ROLE_PREFIX` | 禁用消息角色前缀 | `false` | ❌ |
| `PROMPT_FOR_FILE` | 文件上传时的系统提示词 | 见示例 | ❌ |

### 🔍 搜索功能配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `IGNORE_SEARCH_RESULT` | 忽略搜索结果内容 | `false` | ❌ |
| `SEARCH_RESULT_COMPATIBLE` | 搜索结果兼容模式，禁用折叠块 | `false` | ❌ |

### 🎯 高级功能配置
| 环境变量 | 描述 | 默认值 | 必填 |
|----------|------|--------|------|
| `IGNORE_MODEL_MONITORING` | 忽略模型一致性监控 | `false` | ❌ |
| `IS_MAX_SUBSCRIBE` | 启用 Max 订阅功能 | `false` | ❌ |

 ## 📝 API使用
 ### 认证
 在请求头中包含您的API密钥：
 ```
 Authorization: Bearer YOUR_API_KEY
 ```
 
 ### 聊天完成
 ```bash
 curl -X POST http://localhost:8080/v1/chat/completions \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer YOUR_API_KEY" \
   -d '{
     "model": "claude-4.0-sonnet",
     "messages": [
       {
         "role": "user",
         "content": "你好，Claude！"
       }
     ],
     "stream": true
   }'
 ```
 
 ### 图像分析
 ```bash
 curl -X POST http://localhost:8080/v1/chat/completions \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer YOUR_API_KEY" \
   -d '{
     "model": "claude-4.0-sonnet",
     "messages": [
       {
         "role": "user",
         "content": [
           {
             "type": "text",
             "text": "这张图片里有什么？"
           },
           {
             "type": "image_url",
             "image_url": {
               "url": "data:image/jpeg;base64,..."
             }
           }
         ]
       }
     ]
   }'
 ```

## 🔄 代理池功能

代理池功能为 Pplx2Api 提供了强大的代理管理能力，通过自动获取、轮换和监控代理，显著提高请求的稳定性和成功率。

### ✨ 功能特性
- **🔄 自动获取代理**: 启动时从代理池API自动获取指定数量的代理
- **⚡ 轮询分配**: 每次请求使用不同的代理，提高请求成功率和并发性能
- **🔄 智能轮换**: 代理使用6小时后自动轮换新的一批代理，保持代理池活跃
- **♻️ 重复使用**: 代理节点可重复使用，不是一次性消耗，提高资源利用率
- **📊 状态监控**: 提供API端点查看代理池状态和详细信息
- **🔧 手动刷新**: 支持手动触发代理池刷新，灵活管理
- **⚖️ 负载均衡**: 自动在可用代理间分配请求负载

### ⚙️ 配置说明
```bash
# 启用代理池功能
ENABLE_PROXY_POOL=true

# 代理池API地址（返回纯文本格式：http://user:pass@ip:port）
PROXY_POOL_API=https://your-proxy-api.com/get

# 代理池大小（启动时获取的代理数量，建议10-50）
PROXY_POOL_SIZE=10
```

### 📋 代理池API要求
您的代理池API需要返回以下纯文本格式：
```text
http://username:password@proxy-ip:port
```

**示例响应**：
```text
http://user123:pass456@192.168.1.100:8080
```

**API要求**：
- HTTP GET 请求
- 返回 `Content-Type: text/plain`
- 每次调用返回一个可用的代理地址
- 支持HTTP/HTTPS代理协议

### 🔧 管理API
| 端点 | 方法 | 描述 | 响应格式 |
|------|------|------|----------|
| `/proxy/status` | GET | 查看代理池状态和统计信息 | JSON |
| `/proxy/refresh` | POST | 手动刷新代理池 | JSON |

**状态API响应示例**：
```json
{
  "enabled": true,
  "total_proxies": 10,
  "active_proxies": 8,
  "failed_proxies": 2,
  "last_refresh": "2024-01-15T10:30:00Z",
  "next_refresh": "2024-01-15T16:30:00Z",
  "proxies": [
    {
      "url": "http://user:***@192.168.1.100:8080",
      "status": "active",
      "error_count": 0,
      "last_used": "2024-01-15T10:25:00Z"
    }
  ]
}
```

### 🔄 代理生命周期
1. **🚀 获取阶段**: 启动时从API获取指定数量的代理
2. **⚡ 使用阶段**: 轮询使用代理，可重复使用同一代理
3. **🛡️ 错误处理**:
   - **其他错误**: 增加错误计数，超过5次后删除
   - **超时错误**: 标记为临时不可用，稍后重试
4. **🔄 轮换阶段**: 6小时后自动获取新的一批代理替换旧代理

### 📊 使用优先级
代理选择按以下优先级顺序：
1. **请求参数中的代理**（如果在请求中指定）
2. **代理池中的代理**（如果启用代理池功能）
3. **环境变量PROXY**（全局代理配置）
4. **直连**（无代理）

### 💡 最佳实践
- **代理池大小**: 建议设置为10-50个，根据并发需求调整
- **API稳定性**: 确保代理池API的高可用性和响应速度
- **监控**: 定期检查 `/proxy/status` 端点监控代理池健康状态
- **错误处理**: 代理池会自动处理失效代理，无需手动干预
- **资源管理**: 代理会自动轮换，避免长期占用同一代理资源

