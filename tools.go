package main

import (
	"flag"
	"fmt"
	"net/http"

	"chatgpt-tools/internal/config"
	"chatgpt-tools/internal/handler"
	"chatgpt-tools/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/tools.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors(), rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		if r.URL.Path == "/users" || r.URL.Path == "users/notify/unread" || r.URL.Path == "/users/history/chat" {
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
