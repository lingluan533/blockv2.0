package config

import (


	locallog "blockv2.0/util/log"
	"context"
	"database/sql"
	"fmt"
	logger "github.com/sirupsen/logrus"
	_ "github.com/taosdata/driver-go/taosSql"
	"os"
	"strconv"
	"strings"
	"time"
)


type DBTDengine struct {
	DB *sql.DB
}
var defaultDBTDengine *DBTDengine
func NewDBTDengine() *DBTDengine {
	return defaultDBTDengine
}
var DefaultConfig BlockConfig
var CHAINNAME string
//url格式： "root:taosdata@/tcp(" + configPara.hostName + ":" + strconv.Itoa(configPara.serverPort) + ")/"
var ChainConfig GlobalConfig
var SupportLedger =make(map[string]bool)
func init(){
	logger.Info("初始化配置")
	ChainConfig=Initialize()
	var level logger.Level
	switch ChainConfig.Common.LogConfig.LogLevel {
	case "debug":
		level=logger.DebugLevel
	case "info":
		level=logger.InfoLevel
	case "warn":
		level=logger.WarnLevel
	case "error":
		level=logger.ErrorLevel
	default:
		level=logger.InfoLevel
	}
	logger.SetLevel(level)
	CHAINNAME=ChainConfig.ChainName
	DefaultConfig= ChainConfig.Block
	TDengineConfig:=DefaultConfig.TDengineConfig
	locallog.Init(ChainConfig.Common.LogConfig.RootPath)
	if !ChainConfig.Common.LogConfig.OutputFile {
		logger.SetOutput(os.Stdout)
	}
	path := strings.Join([]string{TDengineConfig.User, ":",TDengineConfig.Passwd, "@/tcp(", TDengineConfig.Hostname, ":",strconv.FormatInt(int64(TDengineConfig.Port),10), ")/"}, "")
	db, err := sql.Open(TDengineConfig.Driver, path)
	if err !=nil{
		panic(err)
	}else{
		fmt.Println("成功连接到tdengine数据库上 ")
	}
	db.SetConnMaxLifetime(1 * time.Second)
	db.SetMaxIdleConns(20)   //最大打开的连接数
	db.SetMaxOpenConns(2000) //设置最大闲置个数
	defaultDBTDengine=&DBTDengine{db}
	defaultDBTDengine.CreateDatabase(TDengineConfig.DBName)
	defaultDBTDengine.CreateDataReceiptsTable(TDengineConfig.DBName)
	defaultDBTDengine.CreateTransactionsTable(TDengineConfig.DBName)
	defaultDBTDengine.CreateBlockheadersTable(TDengineConfig.DBName)


}
//func (client *DBTDengine) Close() {
//	client.DB.Close()
//}
func (client *DBTDengine) CreateDatabase(dbName string) {

	sqlStr := "create database if not exists " + dbName + " keep " + strconv.Itoa(DefaultConfig.TDengineConfig.Keep) + " days " + DefaultConfig.TDengineConfig.TableKeepDay+` precision "us" UPDATE 1;`
	_, err := client.DB.Exec(sqlStr)
	CheckErr(err, sqlStr)
}
func (client *DBTDengine) CreateDataReceiptsTable(dbName string) {
	sqlStr := "create table if not exists "+  dbName + ".datareceipts (createtimestamp timestamp ,entityid nchar(255),keyid nchar(255),receiptvalue double(8),version nchar(255),username nchar(255)," +
		"operationtype nchar(255),datatype nchar(255),servicetype nchar(255),filename nchar(255),filesize double(8),filehash nchar(255),uri nchar(255),parentkeyid nchar(255),attachmentfileuris binary(1000)," +
		"attachmenttotalhash nchar(255),blockid nchar(255));"
	_, err := client.DB.Exec(sqlStr)
	CheckErr(err, sqlStr)
}
func (client *DBTDengine) CreateTransactionsTable(dbName string) {
	sqlStr := "create table if not exists "+  dbName + ".transactions (createtimestamp timestamp ,entityid nchar(255),transactionid nchar(255),initiator nchar(255),recipient nchar(255),txamount double(8)," +
		"datatype nchar(255),servicetype nchar(255),remark nchar(255),blockid nchar(255));"
	_, err := client.DB.Exec(sqlStr)
	CheckErr(err, sqlStr)
}
func (client *DBTDengine) CreateBlockheadersTable(dbName string) {
	sqlStr := "create table if not exists "+  dbName + ".blockheaders (createtimestamp timestamp ,keyid nchar(255),blockheight bigint(8),datatype nchar(255),datavalue nchar(255),updatetimestamp timestamp," +
		"datahash nchar(255),blockhash nchar(255),preblockhash nchar(255),nonce int(4),target int(4),currentdatacount bigint(8),currentdatasize bigint(8), version nchar(255),blocktype nchar(255),ledgertype nchar(255),date nchar(255));"
	_, err := client.DB.Exec(sqlStr)
	CheckErr(err, sqlStr)
}



func (client *DBTDengine) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//后期需要扩展支持事务的判断、支持事务
	return client.DB.Exec(query, args...)
}

func (client *DBTDengine) Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	//访问前提前预处理
	return client.DB.Query(query, args...)
}
func CheckErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("%s\n", prompt)
		panic(err)
	}
}