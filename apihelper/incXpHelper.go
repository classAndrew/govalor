package apihelper

import "github.com/classAndrew/govalor/models"

// CreateXPRecord creates a row of guild name, xp gained, and uuid
func CreateXPRecord(uuid string, name string, guildName string, xpgain uint64, timestamp uint64) {
	models.DB.Create(&models.MemberRecordXP{UUID: uuid,
		Name:      name,
		Guild:     guildName,
		XPGain:    xpgain,
		Timestamp: timestamp})
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
