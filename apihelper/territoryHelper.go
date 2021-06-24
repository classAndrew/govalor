package apihelper

import (
	"fmt"
	"log"

	"github.com/classAndrew/govalor/models"
)

// BatchTerritoryUpdate
func BatchTerritoryUpdate(names []string, guilds []string, helds []string) {
	tx := models.DB.Begin()
	for i := range names {
		q := fmt.Sprintf("UPDATE territories SET guild=\"%s\", held=\"%s\" WHERE name=\"%s\";", guilds[i], helds[i], names[i])
		err := tx.Exec(q).Error
		if err != nil {
			tx.Rollback()
			log.Fatal(err.Error())
			return
		}
	}

	err := tx.Commit().Error
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

// GetTerritories
func GetTerritories() []models.Territory {
	var terrs []models.Territory
	models.DB.Find(&terrs)
	return terrs
}
