package apihelper

import (
	"github.com/classAndrew/govalor/models"
)

func AddActivityMember(name string, guild string, timestamp int64) {
	models.DB.Create(&models.ActivityMember{Name: name, Guild: guild, Timestamp: timestamp})
}
