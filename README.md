# RainIM

## 使用
```
cd proto
# 生成代码
go generate ./...
# 格式检查
alias protofmt="find . -name '*.proto' -print0 | xargs -0 clang-format -i"
protofmt
```

## 特点

- 方便已有的系统接入消息服务

## 模块划分
- msg_gateway: websocket服务
- base: 注册消息服务的客户端

## 技术选型

- RPC框架：ConnectRPC
> 之前使用的是gRPC,但

## 设计选择
### 一、消息的存储结构与查询
> 使用源目的方式
1. 一条消息存储两份，增加is_sender字段, 查询更具源目查询一次即可
2. 消息存储一份，查询的时候源目,正反查询两次
3. 创建conversation表，两个client第一次会创建一个conversation,并且把存储id到消息中，查询一次

** 选择第2种 **
core专注与消息,不要而外的conversation。存储两份占用空间和IO性能。
bun orm 有or查询，和添加符合索引功能。

### 二、消息服务的连接认证
client需要先调用认证接口,发送自己的信息，如终端信息,认证信息。
然后服务端生成sessionId存储在redis中，连接websocket会需要携带sessionId，认证一次。
sessionId有一定有效期，短时间断线后不用再调用认证接口

### 三、多类型终端支持
对于msg_gateway,不需要知道client是什么类型。认证的时候，将终端信息携带上来。

### 四、消息序列seq
max_seq.source_id.target_id、max_seq.target_id。source_id,这样存储同步的两份再redis中。

## 后期考虑

- 认证、通讯安全性问题
- 提供消息已读等需求