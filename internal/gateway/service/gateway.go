package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"rain-im-server/internal/gateway/global"

	"github.com/gorilla/websocket"
)

type GatewayServer struct {
	Addr       string
	MaxConnNum int
	UpGrader   *websocket.Upgrader
	ClientConn map[string]map[uint8]*WSClient
	MsgH       MessageHandle
}

type WSClient struct {
	*websocket.Conn
	PlatformId uint8
	ClientId   string
}

func NewGatewayServer() *GatewayServer {
	gatewayServer := &GatewayServer{
		Addr:       ":5173",
		MaxConnNum: 200,
		ClientConn: make(map[string]map[uint8]*WSClient),
	}

	gatewayServer.UpGrader = &websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}

	return gatewayServer
}

func (g *GatewayServer) Run() {
	http.HandleFunc("/gateway", g.ConnHandler)
	err := http.ListenAndServe(g.Addr, nil)
	if err != nil {
		panic("websocket listening err:" + err.Error())
	}
}

func (g *GatewayServer) ConnHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:确定client_id存在
	// TODO:确定token有效
	// 检测plantform_id 合法
	// var wsConnReq gatewayv1.WebsocketConnRequest
	// 1.解析参数
	r.ParseForm()
	clientId, ok := r.Form["client_id"]
	if !ok {
		log.Println("clientID is none!")
		return
	}
	platformId, ok := r.Form["platform_id"]
	if !ok {
		log.Println("platformID is none!")
		return
	}
	platform_id_uint64, err := strconv.ParseUint(platformId[0], 10, 64)
	if err != nil {
		log.Println("clientID conv to int fail!")
		return
	}

	// 2.校验参数

	//建立websocket连接
	conn, err := g.UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	newConn := &WSClient{
		Conn:       conn,
		PlatformId: uint8(platform_id_uint64),
		ClientId:   clientId[0],
	}
	g.AddClientConn(newConn)
	// 用户连接成功，对连接的数据进行读取
	go g.ReadMsg(newConn)

}

func (g *GatewayServer) AddClientConn(conn *WSClient) {
	// 设置写入超时时间
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))

	connMap := make(map[uint8]*WSClient)
	connMap[conn.PlatformId] = conn
	g.ClientConn[conn.ClientId] = connMap
	log.Println("add client conn")
	// RedisDB.SetClientStatus(conn.clientID, true)
}

func (g *GatewayServer) DelClientConn(conn *WSClient) {
	err := conn.Conn.Close()
	if err != nil {
		log.Println("del conn err:", err)
	}
	delete(g.ClientConn[conn.ClientId], conn.PlatformId)

	if len(g.ClientConn[conn.ClientId]) == 0 {
		delete(g.ClientConn, conn.ClientId)
	}
	// RedisDB.SetClientStatus(conn.clientID, false)
}

func (g *GatewayServer) WriteMsg(conn *WSClient, msgType int, message []byte) error {
	return conn.WriteMessage(msgType, message)
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

func (g *GatewayServer) ParseMsg(conn *WSClient, binaryMsg []byte) {
	log.Println("ParseMsg")
	msgReq := gatewayv1.SingleMessageRequest{}
	err := json.Unmarshal(binaryMsg, &msgReq)
	if err != nil {
		fmt.Println("json err:", err.Error())
	}
	if err := global.Validate.Struct(&msgReq); err != nil {
		log.Println("validate error:", err)
		return
	}

	fmt.Println(msgReq.SourceId.String())

	switch msgReq.MessageType {
	case gatewayv1.Message_MESSAGE_SINGLE:
		log.Println("single message")
	case gatewayv1.Message_MESSAGE_GROUP:
		log.Println("group message")
	default:
		log.Println("clientType error")
	}
}
