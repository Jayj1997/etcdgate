/*
 * @Author       : jayj
 * @Date         : 2021-08-30 18:19:13
 * @Description  : gracefull shutdown
 */

package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func GracefullyShutdown(server *http.Server) {
	// make channel receive closing signal
	done := make(chan os.Signal, 1)

	/**
	os.Interrupt           -> ctrl+c
	syscall.SIGINT|SIGTERM -> kill
	*/
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logrus.Println("closing http server gracefully ...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatalln("closing http server gracefully failed: ", err)
	}

}
