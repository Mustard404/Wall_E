package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

func ErrLog(errMsg error) {
	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 控制台输出错误
	log.Printf("%s\n", errMsg)

	// 写入错误到日志文件
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// 写入错误消息到日志文件
	_, err = logFile.WriteString(fmt.Sprintf("[%s] : %s\n", currentTime, errMsg))
	if err != nil {
		log.Fatalf("Failed to write to log file: %v", err)
	}
}
