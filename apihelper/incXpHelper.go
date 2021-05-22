package apihelper

import (
	"fmt"
	"log"
	"strings"

	"github.com/classAndrew/govalor/models"
)

// CreateXPRecord creates a row of guild name, xp gained, and uuid
func CreateXPRecord(uuid string, name string, guildName string, xpgain uint64, timestamp uint64) {
	models.DB.Create(&models.MemberRecordXP{UUID: uuid,
		Name:      name,
		Guild:     guildName,
		XPGain:    xpgain,
		Timestamp: timestamp})
}

func InsertXPRecordBatch(records []models.MemberRecordXP) {
	tx := models.DB.Begin()
	batchSize := 50
	batches := len(records) / batchSize
	if len(records)%batchSize != 0 {
		batches++
	}
	for i := 0; i < batches; i++ {
		upTo := (i + 1) * batchSize
		if upTo > len(records) {
			upTo = len(records)
		}

		valueStrings := []string{}
		valueArgs := []interface{}{}
		for k := i * batchSize; k < upTo; k++ {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, records[k].UUID, records[k].Name, records[k].Guild, records[k].XPGain, records[k].Timestamp)
		}

		stmt := fmt.Sprintf("INSERT INTO member_record_xps (uuid, name, guild, xp_gain, timestamp) VALUES %s", strings.Join(valueStrings, ","))
		err := tx.Exec(stmt, valueArgs...).Error
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}

	err := tx.Commit().Error
	if err != nil {
		log.Fatalln(err)
	}
}

// GetPlayerIncXP will get a slice of player xp records between 2 timestamps
func GetPlayerIncXP(uuid string, t1 uint64, t2 uint64) map[string]interface{} {
	var res []models.MemberRecordXP
	data := make(map[string]interface{})
	models.DB.Where("uuid = ?", uuid).Having("timestamp >= ?", t1).Having("timestamp <= ?", t2).Find(&res)
	if len(res) == 0 {
		return nil
	}

	data["name"] = res[0].Name
	xpRecordTable := make([][]uint64, len(res))
	for i, v := range res {
		xpRecordTable[i] = []uint64{v.XPGain, v.Timestamp}
	}
	data["values"] = xpRecordTable

	return data
}
