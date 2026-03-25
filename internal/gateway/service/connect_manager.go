package service

import (
	"fmt"
	"sync"
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
	clientConns map[string][]*WSClient // 第二个key = platformId + "-" + string(len(totalOnline))
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

	// 如果配置一个账号在一种平台最多一个在线,那么 %s-%s-0
	// 查询redis,这里要上分布式锁
	key := conn.ClientId + "-" + conn.PlatformId.String() + "-"
	for i := range 100 {
		tempK := key + fmt.Sprintf("%d", i)
		if _, ok := cm.conns[tempK]; !ok {
			key = key + fmt.Sprintf("%d", i)
			break
		}
		// 如果配置一个账号在一种平台最多一个在线,这里直接退出
	}

	cm.conns[key] = conn
	cm.clientConns[conn.ClientId] = append(cm.clientConns[conn.ClientId], conn)
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
