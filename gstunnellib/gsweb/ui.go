package gsweb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func Run(inst gstunnellib.GsStatus) {

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	r.GET("/test2", func(c *gin.Context) {
		html1 := `
		<html>
			<head>
				<h1>hello</h1>
			</head>
			<body></body>
		</html>
		`
		c.Writer.Write([]byte(html1))
	})

	r.GET("/", func(c *gin.Context) {
		html1 := `
		<html>
			<head></head>
			<body>
			<p>GID: %d</p>
			<p></p>
			<p>%s</p>
			</body>
		</html>
		`
		outhtml1 := fmt.Sprintf(html1, inst.GetStatusData().Gid, inst.GetStatusConnList().HTMLString())
		c.Writer.WriteString(outhtml1)
		//c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	r.GET("/json", func(c *gin.Context) {
		c.Writer.WriteString(string(inst.GetJson()))
	})

	r.GET("/js", func(c *gin.Context) {
		v1 := os.Args
		_ = v1
		fd, err := os.Open("../gstunnellib/gsweb/mytest2.html")
		gstunnellib.CheckErrorEx_exit(err, log.Default())
		defer fd.Close()
		gstunnellib.CheckError_panic(err)
		re, err := io.ReadAll(fd)
		gstunnellib.CheckError_panic(err)
		c.Writer.WriteString(string(re))
	})

	r.Run("localhost:8080")
}
