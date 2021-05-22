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
				_ = nowtime

				members, _ := res["members"].([]interface{})

				guildMembers := make([]models.GuildMember, len(members))
				contribMap := make(map[string]int64)
				for i := range members {
					m := members[i].(map[string]interface{})
					memberName, _ := m["name"].(string)
					memberUUID, errUUID := m["uuid"].(string)
					if !errUUID {
						log.Fatal("Error casting uuid")
						return
					}
					gxpFloat, _ := m["contributed"].(float64)
					gxpContrib := int64(gxpFloat)
					guildMembers[i] = models.GuildMember{Name: memberName, UUID: memberUUID, Guild: guildName}
					contribMap[memberUUID] = gxpContrib
				}

				apihelper.UpdateMemberListBatch(guildMembers)
				lastXPs := apihelper.GetGuildMembersXP()

				storedCurrentXP := make(map[string]int64)
				uuidToLastXP := make(map[string]int64)
				for _, lastMember := range lastXPs {
					uuidToLastXP[lastMember.UUID] = lastMember.LastXP
					storedCurrentXP[lastMember.UUID] = lastMember.XP
				}

				updatedXPs := []models.UserTotalXP{}
				sliceDelta := []models.MemberRecordXP{}

				for i := range guildMembers {
					m := guildMembers[i]
					memberUUID := m.UUID
					lastXPs, ok := uuidToLastXP[memberUUID]
					// fmt.Printf("%s %d\n", m.Name, contribMap[memberUUID])
					if !ok {
						// new join.
						apihelper.CreateUserTotalXP(m.Guild, m.Name, contribMap[m.Name], m.UUID)
						continue
					}

					gxpContrib := contribMap[memberUUID]
					delta := gxpContrib - lastXPs

					user := models.UserTotalXP{}
					user.XP = storedCurrentXP[memberUUID]
					if delta < 0 {
						user.XP += gxpContrib
						// apihelper.UpdateUserTotalXP(user)
					} else if delta > 0 {
						user.XP += delta
						sliceDelta = append(sliceDelta, models.MemberRecordXP{Guild: guildName, Name: m.Name, UUID: m.UUID, XPGain: uint64(delta), Timestamp: nowtime})
						// apihelper.CreateXPRecord(user.UUID, user.Name, user.Guild, uint64(delta), )
						// apihelper.UpdateUserTotalXP(user)

						// apihelper.CreateUserTotalXP(user.Guild, user.Name, user.XP, user.UUID)
					}
					user.Guild = guildName
					user.LastXP = gxpContrib
					user.Name = m.Name
					user.UUID = m.UUID

					updatedXPs = append(updatedXPs, user)
				}

				apihelper.UpdateUserTotalXPTX(updatedXPs)
				apihelper.InsertXPRecordBatch(sliceDelta)

			case <-time.After(10 * time.Second):
				log.Println("Failed req from timeout")
			}
		}
		time.Sleep(delay)
	}
}
