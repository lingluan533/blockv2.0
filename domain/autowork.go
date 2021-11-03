package domain

import (
	"github.com/robfig/cron"
	logger "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (


	//每隔10s 进行一次区块头的扫描
	scanblockheadersSpec = "*/10 * * * * *"

	//每分钟的0s 进行上一个区块头是否已经全部接受的验证
	verifyLastBlockheaderSpec = "0 * * * * ?"

	//10分钟执行一次该时间应该接受的所有的区块头的验证
	verifyAllBlockHeaderSpec = "0 */1 * * * ?"
)

func AutoWorkMain(client *clientv3.Client) {
	//block模块启动时候就应该执行一次主动同步 之后每隔
	go ScanAllBlockHeaders(client)
	logger.Info("定时任务已开启定时扫描区块头中。。。")
	c := cron.New()

	//1、每分钟0s执行一次定时任务，判断是否上上个分钟的区块头已经收到了,并且写入了数据库。
	//2、当启动的时候执行一个把之前所有的区块都检查的步骤 从0号到当前时间应该有的区块都检查是否已经在数据库中存在了。
	c.AddFunc(verifyLastBlockheaderSpec, func( ) {
		go ScanLastBLockHeaders(client)
	})
	//当前时间应当得到的区块头的数量是固定的，所以计算应有的和数据库中已有的比较一下即可
	c.AddFunc(scanblockheadersSpec, func() {
		//扫描比对数据库中的区块头和区块记录是否对应（数量）
		go ScanReceivedBlockHeaders(client)
	})
	c.AddFunc(verifyAllBlockHeaderSpec, func() {

		go ScanAllBlockHeaders(client)
	})
	go c.Start()

	select {}
}
