package tooler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func EtcdCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "etcd",
		Short: "use for etcd config",
	}

	cmd.AddCommand(etcdGetCmd())
	cmd.AddCommand(etcdGetPrefixCmd())
	cmd.AddCommand(etcdPutCmd())

	return cmd
}

func etcdGetCmd() *cobra.Command {
	var key, address string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get value from etcd",
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{address}, // etcd 节点地址
				DialTimeout: 5 * time.Second,   // 连接超时时间
			})
			if err != nil {
				log.Fatal(err)
			}
			defer cli.Close()

			ctx := context.Background()

			resp, err := cli.Get(ctx, key)
			if err != nil {
				log.Fatal(err)
			}

			for _, ev := range resp.Kvs {
				fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
			}
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to get from etcd")
	cmd.Flags().StringVarP(&address, "address", "a", "", "address is etcd-server address, example:")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("address")

	return cmd
}

func etcdGetPrefixCmd() *cobra.Command {
	var key, address string

	cmd := &cobra.Command{
		Use:   "getPrefix",
		Short: "Get Prefix value from etcd",
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{address}, // etcd 节点地址
				DialTimeout: 5 * time.Second,   // 连接超时时间
			})
			if err != nil {
				log.Fatal(err)
			}
			defer cli.Close()

			ctx := context.Background()

			resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
			if err != nil {
				log.Fatal(err)
			}

			for _, ev := range resp.Kvs {
				fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
			}
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to get from etcd")
	cmd.Flags().StringVarP(&address, "address", "a", "", "address is etcd-server address, example:")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("address")

	return cmd
}

func etcdPutCmd() *cobra.Command {
	var key, value, address string

	cmd := &cobra.Command{
		Use:   "put",
		Short: "Put key-value pair to etcd",
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{address}, // etcd 节点地址
				DialTimeout: 5 * time.Second,   // 连接超时时间
			})
			if err != nil {
				log.Fatal(err)
			}
			defer cli.Close()

			ctx := context.Background()

			_, err = cli.Put(ctx, key, value)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to put")
	cmd.Flags().StringVarP(&value, "value", "v", "", "Value to put")
	cmd.Flags().StringVarP(&address, "address", "a", "", "address is etcd-server address, example:")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")
	cmd.MarkFlagRequired("address")

	return cmd
}
