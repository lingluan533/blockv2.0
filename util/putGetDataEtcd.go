package util

import (
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)
var(
	AllKeysCounts = "AllKeysCounts"
)




func GetData(cli *clientv3.Client, key string, requestTimeout time.Duration) (getResponse *clientv3.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()
	return resp
}

//根据key前缀获取数据
func GetDataPrefix(cli *clientv3.Client, key string, requestTimeout time.Duration) (getResponse *clientv3.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()
	return resp
}

//统计所有key数量
func AllKeysCount(cli *clientv3.Client, key string, num int, requestTimeout time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()

	allKeysCountStr := "0"
	for _, ev := range resp.Kvs {
		allKeysCountStr = string(ev.Value)
	}
	allKeysCountInt , _ := strconv.Atoi(allKeysCountStr)

	_, err = cli.Put(context.TODO(), key, strconv.Itoa(allKeysCountInt+num))
	if err != nil {
		log.Fatal(err)
	}
}
