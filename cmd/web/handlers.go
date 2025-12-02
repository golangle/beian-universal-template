package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// 定义一个 home 的处理函数，用来响应一个字符串
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	hostName := strings.Split(r.Host, ":") //分割主机名和端口号

	if firstIp == "" {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			// 部署在反向代理后面时，使用X-Forwarded-For头部获取客户端真实IP地址
			// 获取第一个IP地址，即客户端的真实IP
			firstIp = strings.Split(xff, ",")[0]
		} else {
			firstIp = strings.Split(r.RemoteAddr, ":")[0]
		}

		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		writeLog(fmt.Sprintf("[%s] setting first ip:%s \n", formattedTime, firstIp))
	}
	// 根据主机名获取对应的数据行
	line, ok := HostLineMap[hostName[0]]

	//没有对应的主机名，查找默认的
	if !ok {
		line = HostLineMap["*"]
	}
	// 如果默认的也没有，则三个位置全部显示空白，即：什么也不显示

	// 执行模板，传入数据
	err := templateInstance.Execute(w, map[string]interface{}{
		"number":    line.BeiAn,
		"title":     line.Title,
		"copyRight": line.CopyRight,
	})
	if err != nil {
		panic(err)
	}
}

func (app *application) reload(w http.ResponseWriter, r *http.Request) {

	// 1、检查IP地址，决定是否重新加载数据
	if ip_filter {
		var currentIp string

		if firstIp == "" {
			firstIp = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			xff := r.Header.Get("X-Forwarded-For")
			if xff != "" {
				// 部署在反向代理后面时，使用X-Forwarded-For头部获取客户端真实IP地址
				// 获取第一个IP地址，即客户端的真实IP
				currentIp = strings.Split(xff, ",")[0]
			} else {
				currentIp = strings.Split(r.RemoteAddr, ":")[0]
			}

			if firstIp != currentIp {
				formattedTime := time.Now().Format("2006-01-02 15:04:05")
				writeLog(fmt.Sprintf("[%s] rejected reload ip:%s \n", formattedTime, currentIp))
				w.Write([]byte("only the first ip can reload data"))
				return
			}
		}
		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		writeLog(fmt.Sprintf("[%s] accepted reload ip:%s \n", formattedTime, currentIp))
	}

	// 2、检查时间间隔，决定是否重新加载数据
	if startTime.IsZero() {
		startTime = time.Now()
	} else {
		if time.Since(startTime) < time.Duration(refreshInterval)*time.Second {
			w.Write([]byte(fmt.Sprintf("reload too frequently, please wait for %d seconds", refreshInterval)))
			return
		}
	}
	// 3、检查文件的修改时间，决定是否重新加载数据
	if IsNeedReloadData() {
		loadData()
	}
	startTime = time.Now()
	w.Write([]byte("reload success"))
}
