package global

var (
	ConnDetailHashKey = "conn:"
	ConnStatusSetKey  = "client:conn:"
	// client:conversation:{sourceClientId}-{targetClientId} -> conversationId
	// client:conversation:{groupId}-{groupId} -> conversationId
	ClientConversationStringKey = "client:conversation:"
)

var (
	ConnDetailHashKeyTTL = 5
)

type ConnDetail struct {
	ClientId    string `json:"clientId"`
	PlatformId  string `json:"platformId"`
	CountKey    string `json:"countKey"`
	GatewayUUID string `json:"gatewayUUID"`
	CreatedAt   int64  `json:"createdAt"`
}
