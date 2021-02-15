package services

import (
	"log"
	"time"

	"github.com/classAndrew/govalor/apihelper"

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
					name := member.(map[string]interface{})["name"].(string)
					_, online := set[name]
					if online {
						go apihelper.AddActivityMember(name, guildName, timestamp)
					}
				}
				break
			case <-time.After(time.Second * 5):
				log.Println("Timeout get from guild")
			}
		}
		time.Sleep(delay)
	}
}
