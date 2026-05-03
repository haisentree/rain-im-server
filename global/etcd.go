package global

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/google/uuid"
)

var (
	GatewayServerKey = "uuid"

	// 基础服务配置
	EtcdRedisConfig    = "/rian-im-server/%s/redis/config"
	EtcdNatsConfig     = "/rian-im-server/%s/nats/config"
	EtcdPostgresConfig = "/rian-im-server/%s/postgres/config"

	// 服务注册
	EtcdServiceRegisterGateway = "/rian-im-server/%s/service/register/gateway/%s"
)

func init() {
	workMode := "prod"
	GatewayServerKey = uuid.New().String()

	env := os.Getenv("RAIN_IM_SERVER")
	if env != "" && slices.Contains([]string{"prod", "dev", "test"}, workMode) {
		workMode = env
	}

	log.Println("rian-im-server current work mode:", workMode)

	SprintfEtcdWorkMode(workMode)
}

func SprintfEtcdWorkMode(mode string) {
	EtcdRedisConfig = fmt.Sprintf(EtcdRedisConfig, mode)
	EtcdServiceRegisterGateway = fmt.Sprintf(EtcdServiceRegisterGateway, mode, GatewayServerKey)
	EtcdNatsConfig = fmt.Sprintf(EtcdNatsConfig, mode)
	EtcdPostgresConfig = fmt.Sprintf(EtcdPostgresConfig, mode)
}
