package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/classAndrew/govalor/apihelper"
	"github.com/gin-gonic/gin"
)

func GetPlayerIncXP(c *gin.Context) {
	uuid, _t1, _t2 := c.Param("uuid"), c.Param("t1"), c.Param("t2")
	t1, err1 := strconv.Atoi(_t1)
	if err1 != nil {
		log.Fatalln(err1.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
	}
	t2, err2 := strconv.Atoi(_t2)
	if err2 != nil {
		log.Fatalln(err2.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
	}
	res := apihelper.GetPlayerIncXP(uuid, uint64(t1), uint64(t2))
	c.JSON(http.StatusOK, gin.H{"data": res})
}
