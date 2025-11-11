package service

type MessageHandle struct {
}

// func (ws *WServer) parseSingleCommMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.SingleCommMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.SingleMsgReq{
// 		SendID:  msg.SendID,
// 		RecvID:  d.RecvID,
// 		MsgType: uint32(msg.MessageType),
// 		Content: d.Content,
// 	}

// 	resp, err := MsgGatewaySrvClient.ReceiveSingleMsg(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}

// 	log.Println("clientMsg case 10 resp:", resp.Message)
// }

// func (ws *WServer) parseSingleRelayMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.SingleRelayMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 14dd23", err)
// 		return
// 	}

// 	if err := ws.writeMsg(conn, websocket.TextMessage, []byte(d.Content)); err != nil {
// 		log.Println("case1 error:", err)
// 	}
// }

// func (ws *WServer) parseGroupCommMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.GroupCommMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.SingleMsgReq{
// 		SendID:  msg.SendID,
// 		RecvID:  d.RecvID, // 这里就是CollectID
// 		MsgType: uint32(msg.MessageType),
// 		Content: d.Content,
// 	}

// 	resp, err := MsgGatewaySrvClient.ReceiveSingleMsg(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}

// 	log.Println("clientMsg case 10 resp:", resp.Message)

// }

// func (ws *WServer) parseGroupListMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.GroupListMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.ListMsgReq{
// 		SendID:  msg.SendID,
// 		RecvID:  d.RecvIDList,
// 		MsgType: uint32(msg.MessageType),
// 		SeqID:   d.SeqID,
// 		Content: d.Content,
// 	}

// 	resp, err := MsgGatewaySrvClient.ReceiveListMsg(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}

// 	log.Println("clientMsg case 10 resp:", resp.Message)
// }

// func (ws *WServer) parsePullClientMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.PullClientMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.PullClientMsgReq{
// 		OwnerID: d.OwnerID,
// 	}

// 	for _, v := range d.ClientToSeq {
// 		temp := &pbMsgGateway.CommonClientToSeq{ClientID: v.ClientID, SeqID: v.SeqID}
// 		req.ClientToSeq = append(req.ClientToSeq, temp)
// 	}

// 	resp, err := MsgGatewaySrvClient.ControlPullClientMsg(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}
// 	log.Println("clientMsg case 10 resp:", resp.ClientToMsg)
// 	// 将消息返回给conn
// }

// func (ws *WServer) parsePullCollectMsg(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.PullCollectMsgReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.PullCollectMsgReq{}

// 	for _, v := range d.CollectToSeq {
// 		temp := &pbMsgGateway.CommonCollectToSeq{CollectID: v.CollectID, SeqID: v.SeqID}
// 		req.CollectToSeq = append(req.CollectToSeq, temp)
// 	}

// 	resp, err := MsgGatewaySrvClient.ControlPullCollectMsg(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}
// 	log.Println("clientMsg case 10 resp:", resp.CollectToMsg)
// }

// func (ws *WServer) parseGetClientMaxSeq(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.GetClientMaxSeqReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.GetClientMaxSeqReq{
// 		OwnerID:    d.OwnerID,
// 		ClientList: d.ClientList,
// 	}

// 	resp, err := MsgGatewaySrvClient.ControlGetClientMaxSeq(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}
// 	log.Println("clientMsg case 10 resp:", resp.ClientToSeq)
// }

// func (ws *WServer) parseGetCollectMaxSeq(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.GetCollectMaxSeqReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.GetCollectMaxSeqReq{
// 		CollectList: d.CollectList,
// 	}
// 	resp, err := MsgGatewaySrvClient.ControlGetCollectMaxSeq(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}
// 	log.Println("clientMsg case 10 resp:", resp.CollectToSeq)
// }

// func (ws *WServer) parseGetClientStatus(conn *WSClient, msg *pkgMessage.CommonMsg) {
// 	d := &pkgMessage.GetClientStatusReq{}
// 	json.Unmarshal([]byte(msg.Data), &d)
// 	if err := Validate.Struct(d); err != nil {
// 		log.Println("validate error: 1423", err)
// 		return
// 	}

// 	req := &pbMsgGateway.GetClientStatusReq{
// 		ClientIDList: d.ClientList,
// 	}

// 	resp, err := MsgGatewaySrvClient.ControlGetClientStatus(context.Background(), req)
// 	if err != nil {
// 		log.Println("case 10afs error")
// 	}
// 	log.Println("clientMsg case 10 resp:", resp.StatusList)
// }
