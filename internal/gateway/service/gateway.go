package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

type GatewayServer struct {
	Addr     string
	UpGrader *websocket.Upgrader

	ConnManager *ConnectionManager
	RespWriter  *ResponseWriter
	MsgH        *MessageHandle
}

func NewGatewayServer(addr string) (*GatewayServer, error) {

	if addr == "" {
		return nil, fmt.Errorf("addr is not empty!")
	}

	gatewayServer := &GatewayServer{
		Addr:        addr,
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
	http.HandleFunc("/gateway", g.ConnHandler)
	err := http.ListenAndServe(g.Addr, nil)
	if err != nil {
		panic("websocket listening err:" + err.Error())
	}
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
