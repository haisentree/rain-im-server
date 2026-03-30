package global

type GatewayRegister struct {
	Key  string `json:"key"`
	Addr string `json:"addr"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
