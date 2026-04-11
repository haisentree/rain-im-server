package global

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// TOOD
// 放在各自的global中,代码冗余;放在公共global中,不会用到所有的。
// 怎么设计一下,而且目前不好close()
var (
	Etcd   *clientv3.Client
	DB     *bun.DB
	Redis  *redisv8.Client
	Nats   *nats.Conn
	NatsJS nats.JetStreamContext
)

func init() {
	var err error

	ctx := context.Background()
	Etcd = NewEtcd()

	DB = NewPG(ctx)
	Redis = NewRedisDB(ctx)
	Nats = NewNats()
	NatsJS, err = Nats.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	_ = initJetStream(NatsJS)
}

func CloseService() {
	Etcd.Close()
	DB.Close()
	Redis.Close()
	Nats.Close()
}

func NewEtcd() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"8.148.84.185:2379"}, // etcd 节点地址
		DialTimeout: 5 * time.Second,               // 连接超时时间
		Username:    os.Getenv("ETCD_USERNAME"),
		Password:    os.Getenv("ETCD_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}
	// defer cli.Close()

	return cli
}

// /rian-im-server/core/db/postgres
func NewPG(ctx context.Context) *bun.DB {
	var dataSource string

	resp, err := Etcd.Get(ctx, "/core/db/postgres")
	if err != nil {
		slog.Error("err",
			slog.String("err", err.Error()),
		)
	}

	for _, val := range resp.Kvs {
		dataSource = string(val.Value)
	}

	sqldb, err := sql.Open("postgres", dataSource)
	if err != nil {
		panic(err)
	}
	// defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

// /rian-im-server/prod/gateway/redis/config
func NewRedisDB(ctx context.Context) *redisv8.Client {
	var redisCfg RedisConfig
	resp, err := Etcd.Get(ctx, EtcdRedisConfig)
	if err != nil {
		fmt.Println("NewRedisDB err:", err.Error())
	}

	if len(resp.Kvs) == 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &redisCfg); err != nil {
			fmt.Println("err:", err.Error())
			panic("NewRedisDB err")
		}
	}

	return redisv8.NewClient(&redisv8.Options{
		Addr: redisCfg.Addr,
		// Password: redisCfg.Password,
		DB: redisCfg.DB,
	})
}

// /rian-im-server/prod/nats/config
func NewNats() *nats.Conn {
	ctx := context.Background()
	var natsCfg NatsConfig
	resp, err := Etcd.Get(ctx, EtcdNatsConfig)
	if err != nil {
		fmt.Println("NewNats err:", err.Error())
	}

	if len(resp.Kvs) == 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &natsCfg); err != nil {
			fmt.Println("err:", err.Error())
			panic("NewNats err")
		}
	} else {
		fmt.Println("NewNats err: nats config not found")
		panic("NewNats err: nats config not found")
	}

	nc, err := nats.Connect(natsCfg.Url)
	if err != nil {
		log.Fatal(err)
	}

	// 可选：设置连接事件回调
	nc.SetReconnectHandler(func(conn *nats.Conn) {
		log.Println("重新连接成功")
	})

	return nc
}
