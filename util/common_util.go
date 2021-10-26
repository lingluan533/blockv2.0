package util

import (
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"os"
	"time"
)

/**
检查error并抛出panic异常
*/
func CheckErr(err error) {
	if err != nil {
		logger.Info(err)
	}
}
//將datetime格式的时间戳 解析成date格式的
func ParseDate(timeStamp string) string{
	t,err:=time.Parse("2006-01-02 15:04:05",timeStamp)
	if err!=nil{
		panic(err)
	}
	date:=t.Format("2006-01-02")
	return date
}
func ParseDate1(timeStamp string) string{
	t,err:=time.Parse("2006-01-02 15:04:05",timeStamp)
	if err!=nil{
		panic(err)
	}
	date:=t.Format("2006_01_02")
	return date
}
/**
将实例转换为string类型的字符串
*/
func Marshal(v interface{}) string {
	dataJson, err := json.Marshal(v)
	CheckErr(err)
	return fmt.Sprintf("%s\n", dataJson)
}

/**
将时间字符串转换成Long型时间
*/
func TimeConvertFromStringToLong(string string) int64 {
	tm2, err := time.Parse(string, string)
	CheckErr(err)
	return tm2.Unix()
}

//判断文件是否存在
func Exists(path string) (bool, error) {

	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	// 检测是否为路径不存在的错误
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
//时间戳转换成datetime
func TimestampToDate(timestamp int64)(datetime string){
	var format="2006-01-02 15:04:05.000"
	datetime =time.Unix(timestamp,0).Format(format)
	return
}
