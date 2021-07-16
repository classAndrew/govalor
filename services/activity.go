package services

import (
	"log"
	"time"

	"github.com/classAndrew/govalor/apihelper"

	"github.com/classAndrew/govalor/models"

	"github.com/classAndrew/govalor/util"
)

// CheckActivity .
func CheckActivity(guildList []string, delay time.Duration) {
	for {
		ch := make(chan map[string]interface{})
		go util.Get(ch, "https://api.wynncraft.com/public_api.php", map[string]string{
			"action": "onlinePlayers",
		})
		var timestamp int64
		set := make(map[string]bool)
		select {
		case resp := <-ch:
			// the first one will always be some timestamp
			timestamp = int64(resp["request"].(map[string]interface{})["timestamp"].(float64))
			for k, people := range resp {
				// get rid of that pesky header
				if k == "request" {
					continue
				}
				for _, player := range people.([]interface{}) {

					set[player.(string)] = true
				}
			}
		case <-time.After(time.Second * 5):
			log.Println("Timeout get online")
		}
		total := len(guildList)
		for i, guild := range guildList {
			if i%util.ReqPerMinute == 0 {
				time.Sleep(60)
			}
			go util.Get(ch, "https://api.wynncraft.com/public_api.php", map[string]string{
				"action": "guildStats", "command": guild,
			})
		}

		currentOnline := []models.ActivityMember{}
		for i := 0; i < total; i++ {
			select {
			case resp := <-ch:
				members, ok := resp["members"].([]interface{})
				if !ok {
					log.Println("Guild not found")
					continue
				}
				// shadow
				guildName := resp["name"].(string)
				for _, member := range members {
					name, ok := member.(map[string]interface{})["name"].(string)
					if !ok {
						log.Println("That one problem")
						break
					}
					uuid := member.(map[string]interface{})["uuid"].(string)
					_, online := set[name]
					if online {
						currentOnline = append(currentOnline, models.ActivityMember{UUID: uuid, Name: name, Guild: guildName, Timestamp: timestamp})
					}
				}
				break
			case <-time.After(time.Second * 5):
				log.Println("Timeout get from guild")
			}
		}
		go apihelper.AddActivityMemberBulk(currentOnline)
		time.Sleep(delay)
	}
}
