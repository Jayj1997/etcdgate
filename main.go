/*
 * @Author       : jayj
 * @Date         : 2021-11-15 14:08:34
 * @Description  :
 */
package main

import (
	"confcenter/routes"
	"confcenter/utils"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	g := routes.LoadGin()

	server := &http.Server{
		Addr:           ":8080",
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Errorln("server stop or failed, info: ", err)
		}
	}()

	utils.GracefullyShutdown(server)
}
