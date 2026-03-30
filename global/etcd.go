package global

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/google/uuid"
)

var (
	GatewayServerKey           = "uuid"
	EtcdRedisConfig            = "/rian-im-server/%s/redis/config"
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
}
