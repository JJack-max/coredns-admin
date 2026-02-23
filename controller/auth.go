package controller

import (
	"net/http"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/middleware"

	"github.com/gin-gonic/gin"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

var jwt = middleware.NewJWT()

func Login(c *gin.Context) {
	var conf = config.Get()
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	path := conf.UserEtcdPath

	secret, err := service.EtcdGet(path + user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if !service.VerifySecret(string(secret), user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "username or password is incorrect"})
		return
	}

	token, err := jwt.CreateToken(middleware.CustomClaims{Username: user.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": gin.H{"token": token},
	})
}
