package service

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
	// go ws.readMsg(newConn)

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
	// m := pkgMessage.CommonMsg{}
	// json.Unmarshal(binaryMsg, &m)

	// if err := Validate.Struct(m); err != nil {
	// 	log.Println("validate error:", err)
	// 	return
	// }
	// switch m.MessageType {
	// default:
	// 	log.Println("clientType error")
	// }
}

// func (ws *WServer) readMsg(conn *WSClient) {
// 	for {
// 		msgType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 				// 连接正常关闭或正在关闭
// 				log.Println("连接关闭:", err)
// 				ws.delClientConn(conn)
// 			} else {
// 				// 连接异常关闭
// 				log.Println("连接异常关闭:", err)
// 				ws.delClientConn(conn)
// 			}
// 			log.Println("ws conn error", err)
// 			break
// 		}
// 		// 对接收的消息进行处理
// 		log.Printf("recv: %s", message)
// 		log.Printf("msgType: %d", msgType)
// 		log.Printf("platformID: %d", conn.platformID)
// 		err = conn.WriteMessage(websocket.TextMessage, []byte("send"))
// 		if err != nil {
// 			log.Println("readMsg send error:", err)
// 		}
// 		// 如果解析到platformID=0,直接发送给RecvID

// 		if conn.platformID == 0 {
// 			// 1.解析消息
// 			log.Println("recv platformID==1 message")
// 			var messageToMQ pkgPublic.SingleMsgToMQ
// 			json.Unmarshal(message, &messageToMQ)
// 			// 2.查看RecvID是否连接在当前wss
// 			_, ok := ws.wsClientToConn[uint64(messageToMQ.RecvID)]
// 			if ok {
// 				log.Println("RecvID is online")
// 				// 3.发送消息
// 				for k, v := range ws.wsClientToConn[uint64(messageToMQ.RecvID)] {
// 					err := ws.writeMsg(v, websocket.TextMessage, message)
// 					log.Println("RecvID is sending platform:", v)
// 					if err != nil {
// 						log.Println("Sned RecvID error,platform:", k)
// 					}
// 				}

// 			} else {
// 				// 用户连接再其他的wss,通过redis判断用户是否在线
// 				// 后期再改
// 				log.Println("RecvID offonline")
// 			}
// 			// RecvID不在当前wss，message就丢弃

// 		} else {
// 			ws.msgParse(conn, message)
// 		}

// 		// req := &pbMsgGateway.MessageReq{
// 		// 	Type:    "1",
// 		// 	Message: string(message),
// 		// }
// 	}
// }
