/*
 * @Author       : jayj
 * @Date         : 2021-11-16 14:09:39
 * @Description  :
 */

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addFrontend(g *gin.Engine) {

	g.LoadHTMLGlob("./ui/dist/index.html")

	g.StaticFS("/css", http.Dir("./ui/dist/css"))
	g.StaticFS("/js", http.Dir("./ui/dist/js"))
	g.StaticFS("/img", http.Dir("./ui/dist/img"))
	g.StaticFS("/fonts", http.Dir("./ui/dist/fonts"))

	g.Handle("GET", "/ui/*filepath", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
