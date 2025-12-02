package main

import "net/http"

// routes() 方法是 application的成员方法，返回一个servemux的实例，包含整个应用的路由设置。
func (app *application) routes() http.Handler {
	// 初始化一个新的 ServeMux，绑定home 到 / 路径下。
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.home)         // 没有防止子树路径模式 {$}，可以匹配所有以 / 开头的路径
	mux.HandleFunc("GET /reload", app.reload) // 重新加载数据

	return app.logRequest(commonHeaders(mux))
}
