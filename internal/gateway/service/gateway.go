package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"rain-im-server/global"
	"rain-im-server/pkg/utils"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/encoding/protojson"
)

type GatewayServer struct {
	Key        string
	Addr       string
	PublicAddr string

	UpGrader *websocket.Upgrader

	ConnManager *ConnectionManager
	RespWriter  *ResponseWriter
	MsgH        *MessageHandle
}

func NewGatewayServer(addr string) (*GatewayServer, error) {

	gatewayServer := &GatewayServer{
		Key:  global.GatewayServerKey,
		Addr: addr,

		ConnManager: NewConnectionManager(),
		MsgH:        NewMessageHandle(),
	}

	gatewayServer.RespWriter = NewResponseWriter(gatewayServer.ConnManager)

	gatewayServer.UpGrader = &websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}

	return gatewayServer, nil
}

func (g *GatewayServer) Run() {
	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/gateway", g.ConnHandler)

	var listener net.Listener
	var err error

	if g.Addr != "" {
		// 使用指定的地址监听
		listener, err = net.Listen("tcp", g.Addr)
		if err != nil {
			panic("failed to listen on specified address: " + err.Error())
		}

	} else {
		// 随机监听一个端口
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			panic("failed to listen on random port: " + err.Error())
		}
	}

	// 从 listener 中获取实际监听的地址（可能包含 IP 和端口）
	g.Addr = listener.Addr().String()
	_, portStr, err := net.SplitHostPort(g.Addr)
	if err != nil {
		panic("failed to parse address: " + err.Error())
	}
	// 获取公网 IP，拼接为对外暴露的完整地址
	publicIp := utils.GetPublicIP()
	if publicIp == "" {
		g.PublicAddr = g.Addr
	} else {
		g.PublicAddr = publicIp + ":" + portStr
	}

	log.Printf("Gateway server listening on %s, public address: %s", g.Addr, g.PublicAddr)

	// 注册和续期服务
	g.RegisterService()

	// 启动 HTTP 服务（阻塞）
	if err := http.Serve(listener, nil); err != nil {
		panic("websocket serving error: " + err.Error())
	}
}

func (g *GatewayServer) RegisterService() {
	// 仅序列化需要对外暴露的字段，避免不可序列化字段（如 UpGrader、ConnManager）导致错误
	serviceInfo := struct {
		Key        string `json:"key"`
		Addr       string `json:"addr"`
		PublicAddr string `json:"publicAddr"`
	}{
		Key:        g.Key,
		Addr:       g.Addr,
		PublicAddr: g.PublicAddr,
	}
	data, err := json.Marshal(serviceInfo)
	if err != nil {
		panic("failed to marshal gateway service info: " + err.Error())
	}

	// 获取 etcd 客户端（假设已初始化并存储在全局变量中）
	cli := global.Etcd
	if cli == nil {
		panic("etcd client not initialized")
	}

	// 创建租约，TTL 设为 10 秒（可根据实际需求调整）
	lease, err := cli.Grant(context.Background(), 10)
	if err != nil {
		panic("failed to create etcd lease: " + err.Error())
	}

	// 将服务信息写入 etcd，绑定租约
	_, err = cli.Put(context.Background(), global.EtcdServiceRegisterGateway, string(data), clientv3.WithLease(lease.ID))
	if err != nil {
		panic("failed to register gateway service to etcd: " + err.Error())
	}

	// 保持租约活跃（自动续期）
	keepAliveResp, err := cli.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		panic("failed to keep alive lease: " + err.Error())
	}

	// 启动 goroutine 监听续期响应，防止阻塞
	go func() {
		for {
			select {
			case _, ok := <-keepAliveResp:
				if !ok {
					log.Printf("etcd keep alive channel closed for key %s", global.EtcdServiceRegisterGateway)
					return
				}
				// 续期成功，无需额外处理
			}
		}
	}()

	log.Printf("Gateway service registered to etcd with key: %s, lease ID: %x", global.EtcdServiceRegisterGateway, lease.ID)

}

func (g *GatewayServer) ConnHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "解析表单失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	token := r.Form.Get("token")

	var wsConnReq gatewayv1.ConnectRequest
	if err := protojson.Unmarshal([]byte(token), &wsConnReq); err != nil {
		errMsg := fmt.Sprintf("解析 token 失败: %v token内容: %s", err, token)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	if err := g.MsgH.Validater.Struct(&wsConnReq); err != nil {
		errMsg := ""
		for _, e := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("字段%s验证失败: %s\n", e.Field(), e.Tag())
		}
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// TODO:验证token有效
	// TODO:增加时间片校验,防止重放攻击

	//建立websocket连接
	conn, err := g.UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	newConn := &WSClient{
		Conn:       conn,
		PlatformId: wsConnReq.Platform,
		ClientId:   wsConnReq.ClientId,
	}
	g.AddClientConn(newConn)

	go g.ReadMsg(newConn)
}

func (g *GatewayServer) AddClientConn(conn *WSClient) {
	// 设置写入超时时间
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))

	g.ConnManager.Add(conn)

	log.Println("add client conn")
	// 存储到redis中
	// RedisDB.SetClientStatus(conn.clientID, true)
}

func (g *GatewayServer) DelClientConn(conn *WSClient) {
	err := conn.Conn.Close()
	if err != nil {
		log.Println("del conn err:", err)
	}
	g.ConnManager.Remove(conn)

	// TODO:删除缓存
	// delete(g.ClientConn[conn.ClientId], conn.PlatformId)

	// if len(g.ClientConn[conn.ClientId]) == 0 {
	// 	delete(g.ClientConn, conn.ClientId)
	// }
	// RedisDB.SetClientStatus(conn.clientID, false)
}

func (g *GatewayServer) ReadMsg(conn *WSClient) {
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				// 连接正常关闭或正在关闭
				log.Println("连接关闭:", err)
				g.DelClientConn(conn)
			} else {
				// 连接异常关闭
				log.Println("连接异常关闭:", err)
				g.DelClientConn(conn)
			}
			log.Println("ws conn error", err)
			break
		}

		log.Printf("recv: %s", message)
		log.Printf("msgType: %d", msgType)
		log.Printf("platformID: %d", conn.PlatformId)

		err = conn.WriteMessage(websocket.TextMessage, []byte("send"))
		if err != nil {
			log.Println("readMsg send error:", err)
		}

		g.ParseMsg(conn, message)
	}
}

func (g *GatewayServer) ParseMsg(conn *WSClient, b []byte) {
	log.Println("ParseMsg")

	msgReq := gatewayv1.RawMessage{}
	err := protojson.Unmarshal(b, &msgReq)
	if err != nil {
		fmt.Println("json err:", err.Error())
		return
	}
	// if err := global.Validate.Struct(&msgReq); err != nil {
	// 	log.Println("validate error:", err)
	// 	return
	// }

	switch msgReq.Type {
	case gatewayv1.Message_MESSAGE_SINGLE:
		log.Println("single message")

		g.MsgH.SingleMessageHandle(msgReq.Data, conn, g.RespWriter)

	case gatewayv1.Message_MESSAGE_GROUP:
		log.Println("group message")
	default:
		log.Println("clientType error")
	}
}

func (g *GatewayServer) SubMsg() {
	// 订阅主题
	_, err := global.Nats.Subscribe("foo", func(msg *nats.Msg) {
		log.Printf("收到消息: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
}

type WSClient struct {
	*websocket.Conn
	PlatformId gatewayv1.Platform
	ClientId   string
	CountKey   string
}

func (wc *WSClient) WriteToSelf(msg *gatewayv1.RawMessage) error {

	data, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return wc.WriteMessage(websocket.TextMessage, data)
}
