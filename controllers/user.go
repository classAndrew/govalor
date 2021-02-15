package controllers

import (
	"net/http"

	"github.com/classAndrew/govalor/apihelper"

	"github.com/classAndrew/govalor/models"
	"github.com/gin-gonic/gin"
)

// FindUserTotalXP Fetches user total xp, if no member specified, return all
func FindUserTotalXP(c *gin.Context) {
	if c.Param("guild") == "" {
		var users []models.UserTotalXP
		models.DB.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	} else if c.Param("name") == "" {
		var users []models.UserTotalXP
		models.DB.Where("guild = ?", c.Param("guild")).Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	} else {
		var user models.UserTotalXP
		if err := apihelper.FindSpecificUserTotalXP(c.Param("guild"), c.Param("name"), &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

// CreateUserTotalXP POST to create record in DB
func CreateUserTotalXP(c *gin.Context) {
	var input models.UserTotalXP
	// must be json only. requests.post(json={}) not requests.post(data={})
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if already in table
	// var resjson models.UserTotalXPResponse
	// res, err := http.Get(os.Getenv("SCHEMA") + os.Getenv("HOSTNM") + os.Getenv("PORT") + "/usertotalxp/" + input.Name)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// defer res.Body.Close()

	// if err := json.NewDecoder(res.Body).Decode(&resjson); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// if resjson.Data.Name != "" {
	// 	c.JSON(http.StatusAlreadyReported, gin.H{"error": "user already exists"})
	// 	return
	// }
	if err := apihelper.CreateUserTotalXP(input.Guild, input.Name, input.XP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// UpdateUserTotalXP updates user's total guild xp
func UpdateUserTotalXP(c *gin.Context) {
	var input models.UserTotalXP
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.UserTotalXP
	if err := apihelper.UpdateUserTotalXP(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
