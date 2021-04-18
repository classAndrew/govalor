package controllers

import (
	"net/http"
	"strconv"
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
	apihelper.AddActivityMember(member.UUID, member.Name, member.Guild, int64(time.Now().Nanosecond()))
	c.JSON(http.StatusOK, gin.H{"data": member})
}

// FindActivityGuild finds activity of guild / member
func FindActivityGuild(c *gin.Context) {
	// var input models.ActivityGuildInput don't know what I was thinking don't need this... since path param
	guild, timeStartS, timeEndS := c.Param("guild"), c.Param("timeStart"), c.Param("timeEnd")
	var timeStart, timeEnd int64 = 0, int64(time.Now().Nanosecond())
	if timeStartS != "" {
		temp, err := strconv.Atoi(timeStartS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeStart = int64(temp)
	}
	if timeEndS != "" {
		temp, err := strconv.Atoi(timeEndS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeEnd = int64(temp)
	}
	result := apihelper.FindActivityGuild(guild, timeStart, timeEnd)
	// length := len(result)
	count := make(map[int64]int64)
	// pairs := make([][]models.ActivityGuildResult, length)
	for _, v := range result {
		cnt, exists := count[v.Timestamp]
		if !exists {
			count[v.Timestamp] = 0
			cnt = 0
		}
		count[v.Timestamp] = cnt + 1
	}
	c.JSON(http.StatusOK, gin.H{"data": count})
}

// FindActivityMember finds activity of guild / member
func FindActivityMember(c *gin.Context) {
	name, timeStartS, timeEndS := c.Param("name"), c.Param("timeStart"), c.Param("timeEnd")
	var timeStart, timeEnd int64 = 0, int64(time.Now().Nanosecond())
	if timeStartS != "" {
		temp, err := strconv.Atoi(timeStartS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeStart = int64(temp)
	}
	if timeEndS != "" {
		temp, err := strconv.Atoi(timeEndS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeEnd = int64(temp)
	}
	result := apihelper.FindActivityMember(name, timeStart, timeEnd)
	times := make([]int64, len(result))
	for i, v := range result {
		times[i] = v.Timestamp
	}
	c.JSON(http.StatusOK, gin.H{"data": times})
}

// CountGuildCaptain returns {time: guild captain online count} list
func CountGuildCaptain(c *gin.Context) {
	guild, timeStartS, timeEndS := c.Param("guild"), c.Param("timeStart"), c.Param("timeEnd")
	var timeStart, timeEnd int64 = 0, int64(time.Now().Nanosecond())
	if timeStartS != "" {
		temp, err := strconv.Atoi(timeStartS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeStart = int64(temp)
	}
	if timeEndS != "" {
		temp, err := strconv.Atoi(timeEndS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeEnd = int64(temp)
	}
	result := apihelper.FindCaptainActivityGuild(guild, timeStart, timeEnd)
	timeCounts := make(map[int64]int)
	for _, v := range result {
		_, exists := timeCounts[v.Timestamp]
		if !exists {
			timeCounts[v.Timestamp] = 0
		}
		timeCounts[v.Timestamp] += 1
	}
	c.JSON(http.StatusOK, gin.H{"data": timeCounts})
}
