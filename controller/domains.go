package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/service"
)

func GetDomains(c *gin.Context) {
	var conf = config.Get()
	path := conf.Etcd.PathPrefix

	ds, err := service.Domains(path)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": ds,
	})
}
