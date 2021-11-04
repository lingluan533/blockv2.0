package main
//
import (
	"blockv2.0/config"
	dataStruct "blockv2.0/datastruct"
	"blockv2.0/domain"
	"bufio"
	"context"
	_ "database/sql"
	"encoding/json"
	"os"
	"strconv"

	//"encoding/json"
	"flag"
	"fmt"
	logger "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	//"net"
	//os "os"
	"strings"
	"time"
)

var (
	dialTimout = 5*time.Second
	requestTimeout = 10*time.Second
	endpoints = []string{"127.0.0.1:2379"}
	globalconfig = config.Initialize()
)
const(
	VIDEO="video"
	USER="user_behaviour"
	NODE="node_credible"
	SENSOR="sensor"
	ACCESS="service_access"

	BLOCK_MINUTE="MINUTE"
	BLOCK_TENMINUTE="TENMINUTE"
	BLOCK_DAY="DAY"
)
//func rpcService(port string){
//	//启动RPC服务，将启动
//	ls, err := net.Listen("tcp", ":"+port)
//	if err!=nil{
//		panic(err)
//	}
//	//
//	//创建grpc服务
//	gs := grpc.NewServer()
//	//注册服务
//	//五个类型的账本 都实现了blockService，在查询的时候注册blockService即可
//	//todo：配置文件中没有配置的数据不启动
//	//rpc.RegisterAccessLedgerServiceServer(gs,&service.DefaultAccessLedgerService)
//	//rpc.RegisterNodeLedgerServiceServer(gs,&service.DefaultNodeLedgerService)
//	//rpc.RegisterSensorLedgerServiceServer(gs,&service.DefaultSensorLedgerService)
//	//rpc.RegisterVideoLedgerServiceServer(gs,&service.DefaultVideoLedgerService)
//	//rpc.RegisterUserLedgerServiceServer(gs,&service.DefaultUserService)
//	//rpc.RegisterQueryServiceServer(gs,&service.DefaultBlockService)
//	//reflection.Register(gs)
//	//启动服务
//	//fmt.Println("RPC服务 start listening .....")
//	//gs.Serve(ls)
//}
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: dialTimout,
	})
	if err != nil {
		logger.Fatal(err)
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()
	//启动定时任务用来比对区块头和存放的记录数量是否一致
	go domain.AutoWorkMain(cli)
	//因为天块数据和其他块数据类型的watch key值有差异，需要分开处理
	go func(){
		//启动watch机制，检测分钟块、增强块、创世块
		var now string=time.Now().Format("2006-01-02")
		for dayBlockChan := range cli.Watch(context.Background(), now, clientv3.WithPrefix()){
			for _, ev := range dayBlockChan.Events{
				keyStr := string(ev.Kv.Key)
				logger.Error("watch到变化的key为：",keyStr)
				//logger.Error("watch到变化的value为：",string(ev.Kv.Value))
				if len(ev.Kv.Value)==0{
					logger.Error(string(ev.Kv.Value),"value 是空的")
					break
				}
				strs := strings.Split(keyStr, ":")
				if len(strs) < 2 {
					logger.Error("watch 到数据key 格式不正确")
					break
				}
				if len(strs) == 4{
					//如果key的组成部分大于3部分 则说明是天块/分钟块/增强块
					//当检测到key的变化时去找对应的区块文件，如果文件已经收到了，那么判断文件内的数量与区块头记录的数量是否一致
					switch strs[1]{
					//VIDEO and USER are receipt type
					case VIDEO,USER:
						addDataBlock(keyStr,strs,ev.Kv.Value)
					case NODE,SENSOR,ACCESS:
						addTxBlock(keyStr,strs,ev.Kv.Value)
						default:
					}
				}
			}
			now =time.Now().Format("2006-01-02")
		}
	}()

	//go func() {
	//	//启动watch机制，watch 天块
	//	var preday string=(time.Now().Add(-time.Hour*24)).Format("2006-01-02")
	//	for nonDayblockChan := range cli.Watch(context.Background(), preday, clientv3.WithPrefix()) {
	//		for _, ev := range nonDayblockChan.Events {
	//			keyStr := string(ev.Kv.Key)
	//			if len(ev.Kv.Value)==0{
	//				logger.Error("value 是空的")
	//				break
	//			}
	//			strs := strings.Split(keyStr, ":")
	//			if len(strs) < 2 {
	//				logger.Error("watch key 数据不正确")
	//				break
	//			}
	//		}
	//		preday =(time.Now().Add(-time.Hour*24)).Format("2006-01-02")
	//	}
	//}()
	var hostname string
	flag.StringVar(&hostname, "name", "leader", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	for name,val:=range globalconfig.Common.LedgerName{
		if val.Leader==hostname{
			config.SupportLedger[name]=true
		}else {
			for _,name:=range val.Follower{
				if name==hostname{
					config.SupportLedger[name]=true
				}
			}
		}
	}
	var serviceAddress string
	for {
		if val, ok := globalconfig.Consensus.EtcdGroup[hostname]; !ok {
			fmt.Fprintf(os.Stderr, "输入的主机名有误,请重新输入")
			fmt.Printf("请输入当前节点的主机名: ")
			fmt.Scanf("%s",&hostname)
		} else {
			serviceAddress=val.BlockAddress
			break
		}
	}

	//创建tcp监听
	addr:=strings.Split(serviceAddress,":")
	if len(addr)!=2{
		logger.Error("节点的block_grpc配置项有误")
		os.Exit(1)
	}
	//port:=addr[1]
	//go rpcService(port)
	////go etcdDiscov()
	//go RegisterService(cli)
	select {}
}



func addTxBlock(str string, keyStrs []string, block []byte) {
	//如果key的组成部分有三部分 说明检测到的是一个创世块
	if	len(keyStrs)==3&&(string(keyStrs[2])==BLOCK_TENMINUTE||string(keyStrs[2])==BLOCK_MINUTE||string(keyStrs[2])==BLOCK_DAY) {
		//genesisblock:=&rpc.GenesisBlock{}
		//err:=json.Unmarshal(value,genesisblock)
		//if err!=nil{
		//	logger.Error(err)
		//	return
		//}
		//logger.Info("通过watch收到创世块数据:")
		//logger.Infof("%+v\n",genesisblock)
		//ledgerService.AddGenesisBlock(context.Background(),genesisblock)
	}else if len(keyStrs)>3{
		switch keyStrs[2]{
		case BLOCK_DAY:
			//dayBlockHeader:=&rpc.BlockHeader{}
			//json.Unmarshal(blockHeader,dayBlockHeader)
			//
			//logger.Info("通过watch收到天块数据:")
			//logger.Infof("%+v\n",dayBlock)
			//ledgerService.AddDailyBlock(context.Background(),dayBlock)
		case BLOCK_MINUTE:

			minuteTransactionBlock:=dataStruct.MinuteTransactionBlock{}
			//minuteBlockHeader := &dataStruct.BlockHeader{}
			error := json.Unmarshal(block, &minuteTransactionBlock)
			//把区块头赋值给区块头指针变量
			//minuteBlock.Header = minuteBlockHeader
			if error!=nil{
				fmt.Println("反序列化出错：",error)
				return
			}
			logger.Info("通过watch收到分钟块数据:", minuteTransactionBlock)
			logger.Infof("%+v\n",minuteTransactionBlock.Header)
			//如果分钟快中没有记录的话就不做处理+

			//1.首先存入区块头
			domain.AddBlockHeaderToTdengine(keyStrs[0],*minuteTransactionBlock.Header)
			if minuteTransactionBlock.Header.CurrentDataCount==0{
				return
			}else if minuteTransactionBlock.Header.CurrentDataCount>0{
				//1.找对应的区块文件
				//1.1 拼接出目标区块文件名字
				const TIME_LAYOUT = "2006-01-02 15:04:05"
				time,err:=time.Parse(TIME_LAYOUT,minuteTransactionBlock.Header.CreateTimestamp)
				if err!=nil{
					panic(err)
				}
				blockPath :=  "/root/go/hraft/scope"+"/"+strconv.FormatInt(int64(time.Year()),10)+"-"+fmt.Sprintf("%02d",int64(time.Month()))+"-"+fmt.Sprintf("%02d",int64(time.Day()))+"/"+keyStrs[1]+"/"+
					keyStrs[2]+"/"+keyStrs[3]
				//blockPath :=  "/root/"+config.Initialize().ChainName+"/2021-10-26/node_credible/MINUTE/1025"
				logger.Info("目标交易区块文件目录为：",blockPath)

				//1.2打开区块文件
				var file *os.File
				var rawBlockdata []byte
				rawBlockdata = make([]byte,10240)
				var MinuteTransactionBlockFromfile dataStruct.MinuteTransactionBlock
				file, err = os.OpenFile(blockPath,os.O_RDONLY,0777)
				defer file.Close()
				if err != nil{
					logger.Error("打开文件失败,err=",err)
					return
				}
				//2.读入区块文件到结构体数组
				reader := bufio.NewReader(file)
				n,err := reader.Read(rawBlockdata)
				if err!=nil{
					panic(err)
				}
				fmt.Println("n:=",n)
				err = json.Unmarshal(rawBlockdata[:n], &MinuteTransactionBlockFromfile)
				if err != nil{
					logger.Error("从区块文件反序列化失败：",err)
					return
				}
				logger.Info("从区块文件反序列化成功：",MinuteTransactionBlockFromfile)
				//插入tdengine数据库

				domain.AddTransactionBlockToTdengine(MinuteTransactionBlockFromfile)
			}
		}
	}
}

//处理添加存证记录的逻辑
func addDataBlock(keyStr string,keyStrs []string, block[]byte) {
	//如果key的组成部分有三部分 说明检测到的是一个创世块
	if	len(keyStrs)==3&&(string(keyStrs[2])==BLOCK_TENMINUTE||string(keyStrs[2])==BLOCK_MINUTE||string(keyStrs[2])==BLOCK_DAY) {
		//genesisblock:=&rpc.GenesisBlock{}
		//err:=json.Unmarshal(value,genesisblock)
		//if err!=nil{
		//	logger.Error(err)
		//	return
		//}
		//logger.Info("通过watch收到创世块数据:")
		//logger.Infof("%+v\n",genesisblock)
		//ledgerService.AddGenesisBlock(context.Background(),genesisblock)
	}else if len(keyStrs)>3{
		switch keyStrs[2]{
		case BLOCK_DAY:

			//logger.Infof("%+v\n",dayBlock)
			//ledgerService.AddDailyBlock(context.Background(),dayBlock)
		case BLOCK_MINUTE:
			//因为暂时hraft是存储去块体的，所以传过来的数据其实是区块头jia qukuaiti
			minuteBlock:=dataStruct.MinuteDataBlock{}
			//minuteBlockHeader := &dataStruct.BlockHeader{}
			error := json.Unmarshal(block, &minuteBlock)
			//把区块头赋值给区块头指针变量
			//minuteBlock.Header = minuteBlockHeader
			if error!=nil{
				logger.Error("反序列化出错：",error)
				return
			}
			logger.Infof("通过watch收到分钟块数据:%v %v", *minuteBlock.Header,&minuteBlock.DataReceipts)
			//logger.Infof("%+v\n",minuteBlock)
			//如果分钟快中没有记录的话就不做处理

			//1.先存入区块头到tdengine
			domain.AddBlockHeaderToTdengine(keyStrs[0],*minuteBlock.Header)

			if minuteBlock.Header.CurrentDataCount==0{
				return
			}else if minuteBlock.Header.CurrentDataCount>0{
				//1.找对应的区块文件
				//1.1 拼接出目标区块文件名字
				const TIME_LAYOUT = "2006-01-02 15:04:05"
				time,err:=time.Parse(TIME_LAYOUT,minuteBlock.Header.CreateTimestamp)
				if err!=nil{
					panic(err)
				}
				blockPath :=  "/root/go/hraft/scope"+"/"+strconv.FormatInt(int64(time.Year()),10)+"-"+fmt.Sprintf("%02d",int64(time.Month()))+"-"+fmt.Sprintf("%02d",int64(time.Day()))+"/"+keyStrs[1]+"/"+
					keyStrs[2]+"/"+keyStrs[3]
				//blockPath :=  "/root/"+config.Initialize().ChainName+"/2021-10-26/video/MINUTE/915"
				logger.Info("目标存证区块文件目录为：",blockPath)

				//1.2打开区块文件
				var file *os.File
				var rawBlockdata []byte
				rawBlockdata = make([]byte,99999999)
				var minuteBlockFromfile dataStruct.MinuteDataBlock
				file, err = os.OpenFile(blockPath,os.O_RDONLY,0777)
				defer file.Close()
			if err != nil{
				logger.Error("打开文件失败,err=",err)
				return
			}
			//2.读入区块文件到结构体数组
			reader := bufio.NewReader(file)
			n,err := reader.Read(rawBlockdata)
				if err!=nil{
					panic(err)
				}
				fmt.Println("n:=",n)
			err = json.Unmarshal(rawBlockdata[:n], &minuteBlockFromfile)
			if err != nil{
				logger.Error("从区块文件反序列化失败：",err)
				return
			}
			logger.Info("从区块文件反序列化成功：",minuteBlockFromfile)
			//插入tdengine数据库

				domain.AddDataBlockToTdengine(minuteBlockFromfile)
			}
		}
	}
}
