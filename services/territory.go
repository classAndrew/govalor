package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/classAndrew/govalor/ws"

	"github.com/classAndrew/govalor/apihelper"
)

// TerritoryTrack tracks territories every delay seconds
func TerritoryTrack(delay time.Duration) {
	for {
		res, err := http.Get("https://api.wynncraft.com/public_api.php?action=territoryList")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer res.Body.Close()
		strData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		var dataTerr map[string]interface{}
		json.Unmarshal(strData, &dataTerr)

		terrs, ok := dataTerr["territories"].(map[string]interface{})
		if !ok {
			log.Fatalln("Error in getting territories.")
		}

		oldTerrs := apihelper.GetTerritories()

		changes := []map[string]string{}

		for _, v := range oldTerrs {
			// track changes
			attacker := terrs[v.Name].(map[string]interface{})["guild"].(string)
			if attacker != v.Guild {
				changes = append(changes, map[string]string{
					"attacker":  attacker,
					"defender":  v.Guild,
					"territory": v.Name,
					"held":      v.Held,
				})
			}
		}

		stringified, err := json.Marshal(changes)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if len(changes) > 0 {
			ws.M.Broadcast(stringified)

			names := make([]string, len(changes))
			guilds := make([]string, len(changes))
			helds := make([]string, len(changes))
			for i, v := range changes {
				names[i] = v["territory"]
				guilds[i] = v["attacker"]
				helds[i] = v["held"]
			}
			apihelper.BatchTerritoryUpdate(names, guilds, helds)
		}

		time.Sleep(delay)
	}
}
