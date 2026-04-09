package service

import (
	"context"
	"log"
	"rain-im-server/global"
	"rain-im-server/pkg/utils"
	"sync"
	"time"
)

// type ResponseWriter interface {

// type HandlerFunc func(ResponseWriter, *Msg)

// // ServeDNS calls f(w, r).
// func (f HandlerFunc) ServeDNS(w ResponseWriter, r *Msg) {
// 	f(w, r)
// }

type ConnectionManager struct {
	mu          sync.RWMutex
	conns       map[string]*WSClient   // key = clientId+"-"+platformId + "-" + string(len(totalOnline))
	clientConns map[string][]*WSClient // key = clientId
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns:       make(map[string]*WSClient),
		clientConns: make(map[string][]*WSClient),
	}
}

// Add 添加连接，线程安全
func (cm *ConnectionManager) Add(conn *WSClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 如果配置一个账号在一种平台最多一个在线,那么 %s-%s-rand(6)
	preKey := conn.ClientId + "-" + conn.PlatformId.String() + "-"

	var setKey, countKey string
	for {
		countKey = utils.BaseRandomKey(6)
		setKey = preKey + countKey
		// 查询redis,这里要上分布式锁
		if _, ok := cm.conns[setKey]; !ok {
			break
		}
	}

	conn.CountKey = countKey
	cm.conns[setKey] = conn
	cm.clientConns[conn.ClientId] = append(cm.clientConns[conn.ClientId], conn)
	// 存储到redis
	// 构建结构体实例
	detail := global.ConnDetail{
		ClientId:    conn.ClientId,
		PlatformId:  conn.PlatformId.String(),
		GatewayUUID: global.GatewayServerKey,
		CreatedAt:   time.Now().Unix(),
	}
	dataMap, err := utils.StructData2Map(detail)
	if err != nil {
		// 处理错误，建议记录日志并可能返回错误（但 Add 方法目前无返回值）
		log.Printf("StructData2Map error: %v", err)
		return
	}

	ctx := context.Background()

	exist, _ := global.Redis.Exists(ctx, global.ConnDetailHashKey+setKey).Result()
	if exist == 1 {
		panic("key is exist")
	}

	pipe := global.Redis.Pipeline()

	// 使用 map 写入 HSet
	pipe.HSet(ctx, global.ConnDetailHashKey+setKey, dataMap)

	pipe.SAdd(ctx, global.ConnStatusSetKey+conn.ClientId, setKey)

	// 设置 TTL
	pipe.Expire(ctx, global.ConnDetailHashKey+setKey, 24*time.Hour)
	pipe.Expire(ctx, global.ConnStatusSetKey+conn.ClientId, 24*time.Hour)

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Redis pipeline exec error: %v", err)
	}
}

// Remove 移除连接，线程安全
func (cm *ConnectionManager) Remove(conn *WSClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := conn.ClientId + "-" + conn.PlatformId.String() + "-" + conn.CountKey
	delete(cm.conns, key)

	// 从 clientConns 中移除
	conns := cm.clientConns[conn.ClientId]
	for i, c := range conns {
		if c == conn {
			cm.clientConns[conn.ClientId] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(cm.clientConns[conn.ClientId]) == 0 {
		delete(cm.clientConns, conn.ClientId)
	}

	// 同步删除 Redis 中的记录（建议补充）
	ctx := context.Background()
	pipe := global.Redis.Pipeline()
	pipe.HDel(ctx, global.ConnDetailHashKey+key)               // 删除 hash 字段
	pipe.SRem(ctx, global.ConnStatusSetKey+conn.ClientId, key) // 从集合移除
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
}

// Renew 刷新指定连接的 Redis TTL，应在连接活跃时调用（如收到 Ping/Pong 或业务消息）
func (cm *ConnectionManager) Renew(ctx context.Context, conn *WSClient) error {
	if conn == nil || conn.CountKey == "" {
		return nil
	}

	// 设置 TTL，建议与 Add 中保持一致（例如 3 分钟）
	ttl := 3 * time.Minute

	pipe := global.Redis.Pipeline()
	pipe.Expire(ctx, global.ConnDetailHashKey+conn.CountKey, ttl)
	pipe.Expire(ctx, global.ConnStatusSetKey+conn.ClientId, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Renew TTL error for conn %s: %v", conn.CountKey, err)
	}
	return err
}

// 不需要这个方法
// GetByClientAndPlatform 根据用户ID和平台获取连接
// func (cm *ConnectionManager) GetByClientAndPlatform(clientId string, platform gatewayv1.Platform) *WSClient {
// 	cm.mu.RLock()
// 	defer cm.mu.RUnlock()

// 	key := clientId + platform.String()
// 	return cm.conns[key]
// }

// 获取client的所有连接
func (cm *ConnectionManager) GetClientConns(clientId string) []*WSClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conns := cm.clientConns[clientId]
	if len(conns) == 0 {
		return nil
	}

	// 不会修改
	// // 返回副本，避免外部修改
	// result := make([]*WSClient, len(conns))
	// copy(result, conns)

	return conns
}
