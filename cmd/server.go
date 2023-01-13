package cmd

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
	"os"

	"student-scope-send/controller"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AllowHeaders = append(config.AllowHeaders, "x-requested-with")
		r.Use(cors.New(config))
		r.GET("api/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		// 上传成绩单表格，并异步进行导出
		r.POST("api/upload", controller.Upload)
		// 查询导出成绩单结果
		r.GET("api/query", controller.Query)
		// 发送成绩单邮件
		r.POST("api/send", controller.Send)
		r.Static("api/export", "files/export")
		if err := r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
			return fmt.Errorf("run server: %w", err)
		}
		return nil
	},
}
