package service

import (
	"context"
	"fmt"
	"log"
	"rain-im-server/global"
	"rain-im-server/pkg/utils"
	"sync"
	"time"
)

type ConnectionGetter interface {
	GetLocalClientConns(clientId string) []*WSClient
	GetRemoteConnsInfo(clientId string) ([]global.ConnDetail, error)
}

type ConnectionManager struct {
	mu          sync.RWMutex
	conns       map[string]*WSClient   // key = clientId+"-"+platformId + "-" + random(6)
	clientConns map[string][]*WSClient // key = clientId
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns:       make(map[string]*WSClient),
		clientConns: make(map[string][]*WSClient),
	}
}

func (cm *ConnectionManager) Add(conn *WSClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

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

	detail := global.ConnDetail{
		ClientId:    conn.ClientId,
		PlatformId:  conn.PlatformId.String(),
		GatewayUUID: global.GatewayServerKey,
		CreatedAt:   time.Now().Unix(),
	}
	dataMap, err := utils.StructData2Map(detail)
	if err != nil {
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
	pipe.Expire(ctx, global.ConnDetailHashKey+setKey, 5*time.Minute)
	pipe.Expire(ctx, global.ConnStatusSetKey+conn.ClientId, 5*time.Minute)

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

	ctx := context.Background()
	pipe := global.Redis.Pipeline()
	pipe.Del(ctx, global.ConnDetailHashKey+key)                // 删除 hash 字段
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

	ttl := 5 * time.Minute

	pipe := global.Redis.Pipeline()
	pipe.Expire(ctx, global.ConnDetailHashKey+conn.CountKey, ttl)
	pipe.Expire(ctx, global.ConnStatusSetKey+conn.ClientId, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Renew TTL error for conn %s: %v", conn.CountKey, err)
	}
	return err
}

// 获取本地client的所有连接
func (cm *ConnectionManager) GetLocalClientConns(clientId string) []*WSClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conns := cm.clientConns[clientId]
	if len(conns) == 0 {
		return nil
	}

	// 避免外部修改,返回副本(使用上调用者不去修改)
	return conns
}

// 从redis中获取远程client的连接信息，排除掉本地连接
func (cm *ConnectionManager) GetRemoteConnsInfo(clientId string) ([]global.ConnDetail, error) {
	ctx := context.Background()
	setKey := global.ConnStatusSetKey + clientId

	// 1. 获取该 clientId 下所有的 setKey（成员）
	members, err := global.Redis.SMembers(ctx, setKey).Result()
	if err != nil {
		return nil, fmt.Errorf("redis SMembers %s failed: %w", setKey, err)
	}

	var detailList []global.ConnDetail
	for _, member := range members {
		hashKey := global.ConnDetailHashKey + member

		// 2. 获取 Hash 的全部字段
		fields, err := global.Redis.HGetAll(ctx, hashKey).Result()
		if err != nil {
			log.Printf("redis HGetAll %s error: %v", hashKey, err)
			continue // 单条失败不影响其他
		}

		// 3. 转换为 ConnDetail 结构
		var detail global.ConnDetail
		utils.MapData2Struct(fields, &detail)
		if err != nil {
			log.Printf("parse ConnDetail from %s error: %v", hashKey, err)
			continue
		}

		// 4. 排除本地网关的连接
		if detail.GatewayUUID == global.GatewayServerKey {
			continue
		}

		detailList = append(detailList, detail)
	}

	return detailList, nil
}
