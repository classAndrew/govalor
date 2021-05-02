package apihelper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/classAndrew/govalor/models"
	"github.com/classAndrew/govalor/util"
)

// AddActivityMember .
func AddActivityMember(uuid string, name string, guild string, timestamp int64) {
	models.DB.Create(&models.ActivityMember{UUID: uuid, Name: name, Guild: guild, Timestamp: timestamp})
}

// AddActivityMemberBulk adds multiple activity records
func AddActivityMemberBulk(members []models.ActivityMember) {
	tx := models.DB.Begin()
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
			valueStrings = append(valueStrings, "(?, ?, ?, ?)")
			valueArgs = append(valueArgs, members[k].UUID, members[k].Name, members[k].Guild, members[k].Timestamp)
		}

		stmt := fmt.Sprintf("INSERT INTO activity_members (uuid, name, guild, timestamp) VALUES %s", strings.Join(valueStrings, ","))

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

// FindCaptainActivityGuild
func FindCaptainActivityGuild(guild string, timeStart, timeEnd int64) []models.ActivityMember {
	results := FindActivityGuild(guild, timeStart, timeEnd)
	captainPlus := make(map[string]bool)

	// find all the captains+ of the guild
	ch := make(chan map[string]interface{})
	go util.Get(ch, "https://api.wynncraft.com/public_api.php", map[string]string{"action": "guildStats", "command": guild})
	select {
	case res := <-ch:
		members, ok := res["members"]
		if !ok {
			log.Fatalln("Request to get captains failed")
			break
		}
		mem, _ := members.([]interface{})
		for _, _v := range mem {
			v, _ := _v.(map[string]interface{})
			if rank := v["rank"].(string); rank != "RECRUIT" && rank != "RECRUITER" {
				captainPlus[v["name"].(string)] = true
			}
		}
	case <-time.After(10 * time.Second):
		log.Fatalf("Could not fetch captains of %s\n", guild)
	}

	filtered := make([]models.ActivityMember, 0)
	for _, v := range results {
		_, isCaptain := captainPlus[v.Name]
		if isCaptain {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
