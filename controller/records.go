package controller

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
)

//GetRecords gin handle, get all records form etcd database
func GetRecords(c *gin.Context) {
	var conf = config.Get()
	ex, err := service.EtcdGetItems(conf.Etcd.PathPrefix)

	if err != nil {
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	data := []*model.Record{}
	for _, e := range ex {
		r := e.ToRecord()
		if r == nil || r.Path != conf.Etcd.PathPrefix {
			continue
		}
		data = append(data, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}

func PostRecord(c *gin.Context) {
	var rec model.Record
	var conf = config.Get()
	err := c.ShouldBindJSON(&rec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if rec.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Name field"})
		return
	}
	if rec.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Content field"})
		return
	}
	if rec.Type.String() == `""` {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Type field invalid"})
		return
	}
	rec.Content = strings.Trim(rec.Content, " .")
	rec.Path = conf.Etcd.PathPrefix
	etcd, err := rec.ToEtcd()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	kvs, err := service.EtcdGetKvs(etcd.Key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if len(kvs) > 0 {
		bp, abKeys := maxBasicPrefix(kvs)

		if len(abKeys) > 0 {
			abKvs := make(map[string]string)
			for _, k := range abKeys {
				bp = growBasicPrefix(bp)
				abKvs[k+"/"+bp] = string(kvs[k])
			}

			err := service.EtcdPutKvs(abKvs, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			}
		}
		etcd.Key += "/" + growBasicPrefix(bp)
	}
	err = service.EtcdPutItems(etcd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func maxBasicPrefix(kvs map[string][]byte) (bPrefix string, abKeys []string) {
	// bPrefix = "#0"

	for k := range kvs {
		kPart := strings.Split(k, "/")
		if !strings.HasPrefix(kPart[len(kPart)-1], "#") {
			abKeys = append(abKeys, k)
			continue
		}

		switch {
		case len(bPrefix) < len(kPart[len(kPart)-1]):

			bPrefix = kPart[len(kPart)-1]

		case len(bPrefix) == len(kPart[len(kPart)-1]):
			if strings.Compare(bPrefix, kPart[len(kPart)-1]) < 0 {
				bPrefix = kPart[len(kPart)-1]

			}
		}
	}
	return bPrefix, abKeys
}

func growBasicPrefix(bPrefix string) string {
	if !strings.HasPrefix(bPrefix, "#") {
		return "#1"
	}
	if 49 > bPrefix[len(bPrefix)-1] || bPrefix[len(bPrefix)-1] > 56 {
		return bPrefix + "1"
	}
	bbp := []byte(bPrefix)
	bbp[len(bbp)-1]++
	return string(bbp)
}

func DeleteRecord(c *gin.Context) {
	pk := c.Param("key")

	key, _ := base64.RawURLEncoding.DecodeString(pk)

	if string(key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Key field"})
		return
	}
	err := service.EtcdDelete(string(key))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"msg": "success",
	})
}

func PutRecord(c *gin.Context) {
	key, err := base64.RawURLEncoding.DecodeString(c.Param("key"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	var rec model.Record

	err = c.ShouldBindJSON(&rec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if rec.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Name field"})
		return
	}
	if rec.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Empty value for Content field"})
		return
	}
	if rec.Type.String() == `""` {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Type field invalid"})
		return
	}
	rec.Content = strings.Trim(rec.Content, " .")

	etcd, err := rec.ToEtcd()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	etcd.Key = string(key)
	err = service.EtcdPutItems(etcd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
