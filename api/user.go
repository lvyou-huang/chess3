package api

import (
	"chess/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	return router
}

// 注册
func Register(c *gin.Context) {
	data, _ := c.GetRawData()
	var m map[string]string
	_ = json.Unmarshal(data, &m)
	name := m["name"]
	password := m["password"]
	if err := service.Register(name, password); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"msg": "注册失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

func Login(c *gin.Context) {
	//c.HTML(200, "login.html", gin.H{})

	data, _ := c.GetRawData()
	var m map[string]string
	_ = json.Unmarshal(data, &m)
	name := m["name"]
	password := m["password"]
	if err := service.Login(name, password, c); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"msg": "登录失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "登录成功",
	})

}
