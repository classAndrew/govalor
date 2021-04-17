package services

import (
	"log"
	"time"

	"github.com/classAndrew/govalor/apihelper"

	"github.com/classAndrew/govalor/models"

	"github.com/classAndrew/govalor/util"
)

// UpdateMemberXP Continuously update members for all guilds goroutine
func UpdateMemberXP(guildList []string, delay time.Duration) {
	for {
		reqs := 0
		total := len(guildList)
		// maps are pass by addr (they are reference types)
		ch := make(chan map[string]interface{})
		for _, guild := range guildList {
			if reqs == util.ReqPerMinute {
				// 5 second timeout
				time.Sleep(time.Second * 5)
				reqs = 0
			}
			go util.Get(ch, "https://api.wynncraft.com/public_api.php",
				map[string]string{"action": "guildStats", "command": guild})
			reqs++
		}
		for i := 0; i < total; i++ {
			select {
			case res := <-ch:
				// type assertion
				guildName, ok := res["name"].(string)
				if !ok {
					continue
				}
				reqInfo, errorNowtime := res["request"].(map[string]interface{})
				if !errorNowtime {
					log.Fatalln("Error casting request info to interface{}")
					continue
				}
				_nowtime, _ := reqInfo["timestamp"].(float64)
				nowtime := uint64(_nowtime)

				members, _ := res["members"].([]interface{})
				for i := range members {
					go func(m map[string]interface{}, guildName string) {
						memberName, _ := m["name"].(string)
						memberUUID, errUUID := m["uuid"].(string)
						if !errUUID {
							log.Fatal("Error casting uuid")
							return
						}
						gxpFloat, _ := m["contributed"].(float64)
						gxpContrib := int64(gxpFloat)
						var user models.UserTotalXP
						err := apihelper.FindSpecificUserTotalXP(guildName, memberName, &user)
						if err == nil {
							// the user has left so their current xp count is less than previous
							delta := gxpContrib - user.LastXP
							user.LastXP = gxpContrib
							if delta < 0 {
								user.XP += gxpContrib
								apihelper.UpdateUserTotalXP(user)
							} else if delta > 0 {
								user.XP += delta
								apihelper.CreateXPRecord(memberUUID, memberName, guildName, uint64(delta), nowtime)
								apihelper.UpdateUserTotalXP(user)
							}
							// don't bother sending update query if delta = 0
						} else if err.Error() == "record not found" {
							// user doesn't exist, make one
							if err := apihelper.CreateUserTotalXP(guildName, memberName, gxpContrib); err != nil {
								log.Fatalln(err.Error())
							}
						} else {
							log.Fatalln(err.Error())
						}
					}(members[i].(map[string]interface{}), guildName)
				}
			case <-time.After(10 * time.Second):
				log.Println("Failed req from timeout")
			}
		}
		time.Sleep(delay)
	}
}
