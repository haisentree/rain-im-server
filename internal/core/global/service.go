package global

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	EtcdClient *clientv3.Client
	DB         *bun.DB
)

func init() {
	ctx := context.Background()
	EtcdClient = initEtcd()

	DB = initDB(ctx)

}

func initEtcd() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"8.148.84.185:2379"}, // etcd 节点地址
		DialTimeout: 5 * time.Second,               // 连接超时时间
	})
	if err != nil {
		log.Fatal(err)
	}
	// defer cli.Close()

	return cli
}

// /core/db/postgres
func initDB(ctx context.Context) *bun.DB {
	var dataSource string

	resp, err := EtcdClient.Get(ctx, "/core/db/postgres")
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
