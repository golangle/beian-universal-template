package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// 定义application结构体，保存全局共享参数和对象
type application struct {
	logger *slog.Logger
}

func main() {

	loadConfigFile(configPath)
	loadData()

	addr := flag.String("addr", ":8901", "服务运行端口")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}))

	app := &application{
		logger: logger,
	}

	// 启动打印日志信息
	logger.Info("starting server", slog.String("addr", *addr))
	logger.Info("reload limited time", slog.Any("refreshInterval", refreshInterval))

	// 使用 http.ListenAndServe() 监听 8901 端口
	err := http.ListenAndServe(*addr, app.routes())

	// 错误输出
	logger.Error(err.Error())
	os.Exit(1)
}
