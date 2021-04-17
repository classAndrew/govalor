package apihelper

import (
	"github.com/classAndrew/govalor/models"
)

// AddActivityMember .
func AddActivityMember(uuid string, name string, guild string, timestamp int64) {
	models.DB.Create(&models.ActivityMember{UUID: uuid, Name: name, Guild: guild, Timestamp: timestamp})
}

// FindActivityGuild .
func FindActivityGuild(guild string, timeStart, timeEnd int64) []models.ActivityMember {
	var results []models.ActivityMember
	models.DB.Where("guild = ?", guild).Having("timestamp >= ?", timeStart).Having("timestamp <= ?", timeEnd).Find(&results)
	return results
}

// FindActivityMember .
func FindActivityMember(name string, timeStart, timeEnd int64) []models.ActivityMember {
	var results []models.ActivityMember
	models.DB.Where("name = ?", name).Having("timestamp >= ?", timeStart).Having("timestamp <= ?", timeEnd).Find(&results)
	return results
}
