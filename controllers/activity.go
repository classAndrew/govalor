package controllers

import (
	"net/http"
	"time"

	"github.com/classAndrew/govalor/apihelper"

	"github.com/classAndrew/govalor/models"
	"github.com/gin-gonic/gin"
)

// AddActivityMember .
func AddActivityMember(c *gin.Context) {
	var member models.ActivityMember
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	apihelper.AddActivityMember(member.Name, member.Guild, int64(time.Now().Nanosecond()))
	c.JSON(http.StatusOK, gin.H{"data": member})
}
