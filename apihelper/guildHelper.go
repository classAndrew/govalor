package apihelper

import (
	"fmt"
	"log"
	"strings"

	"github.com/classAndrew/govalor/models"
)

// UpdateMemberListBatch
func UpdateMemberListBatch(members []models.GuildMember) {
	tx := models.DB.Begin()

	err := tx.Exec("DELETE FROM guild_members").Error
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	batchSize := 50
	batches := len(members) / batchSize
	if len(members)%batchSize != 0 {
		batches++
	}
	for i := 0; i < batches; i++ {
		upTo := (i + 1) * batchSize
		if upTo > len(members) {
			upTo = len(members)
		}

		valueStrings := []string{}
		valueArgs := []interface{}{}
		for k := i * batchSize; k < upTo; k++ {
			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, members[k].UUID, members[k].Name, members[k].Guild)
		}

		stmt := fmt.Sprintf("INSERT INTO guild_members (uuid, name, guild) VALUES %s", strings.Join(valueStrings, ","))

		err = tx.Exec(stmt, valueArgs...).Error
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		log.Fatalln(err)
	}
}
