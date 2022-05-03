package main

import (
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"flag"
	"github.com/jinzhu/configor"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "configuration file")
	flag.Parse()
	configor.Load(&Config, configPath)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		Config.Username:  Config.Password,
	}))
	
	authorized.GET("/netstat", func(c *gin.Context) {
		stats, err := Stats()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": stats})
	})
	r.Run(Config.Listen)
}