package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

func main() {
	//demo1()
	//demo2()
	//demo3()
	demo4()
}

func demo1() {
	var logger *zap.Logger
	logger, _ = zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	url := "https://www.baidu.com"
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("fetching error", zap.String("url", url), zap.Error(err))
	} else {
		logger.Info("Success..", zap.String("statusCode", resp.Status), zap.String("url", url))
		resp.Body.Close()
	}
}

func demo2() {
	var sugarLogger *zap.SugaredLogger
	logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
	defer sugarLogger.Sync()

	url := "https://www.baidu.com"
	sugarLogger.Debugf("Trying to hit GET request for %s", url)

	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error("fetching error")
		sugarLogger.Errorf("fetching error,url:%s,err:%s", url, err)
	} else {
		sugarLogger.Info("success..")
		sugarLogger.Infof("success...|url:%s", url)

		resp.Body.Close()
	}
}

func demo3() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())

	file, _ := os.Create("/data/go/go-demo/demo/log/test.log")
	defer file.Close()

	//写入配置
	writeSyncer := zapcore.AddSync(file)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	defer sugarLogger.Sync()

	sugarLogger.Debug("你的日志 debug")
	sugarLogger.Info("你的日志 info")
	sugarLogger.Error("你的日志 error")
}

func demo4() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())

	//file, _ := os.Create("./demo.log")
	//defer file.Close()

	file := &lumberjack.Logger{
		Filename:   "/data/go/go-demo/demo/log/test.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		LocalTime:  true,
		Compress:   true, // disabled by default
	}

	//写入配置
	writeSyncer := zapcore.AddSync(file)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	defer sugarLogger.Sync()
	for true {
		sugarLogger.Debug("你的日志 debug")
		sugarLogger.Info("你的日志 info")
		sugarLogger.Error("你的日志 error")
	}

}
