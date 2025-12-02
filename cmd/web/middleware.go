package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func commonHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")                         // 允许的方法
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length") // 允许的头部字段
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // 允许的方法
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			host   = r.Host
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
			path   = r.URL.Path
		)

		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			// 部署在反向代理后面时，使用X-Forwarded-For头部获取客户端真实IP地址
			// 获取第一个IP地址，即客户端的真实IP
			clientIP := strings.Split(xff, ",")[0]
			ip = clientIP
		}
		//格式化当前时间
		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		writeLog(fmt.Sprintf("[%s] host:%s remoteIp:%s protocol:%s method:%s uri:%s path:%s \n", formattedTime, host, ip, proto, method, uri, path))
		app.logger.Info("-->", "host", host, "ip", ip, "proto", proto, "method", method, "uri", uri, "path", path)
		next.ServeHTTP(w, r)
	})
}
