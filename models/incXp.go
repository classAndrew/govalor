package models

// MemberRecordXP A row of member xp
type MemberRecordXP struct {
	UUID      string `json:"uuid" gorm:"varchar(36);primaryKey"`
	Name      string `json:"name"`
	Guild     string `json:"guild"`
	XPGain    uint64 `json:"xpgain"`
	Timestamp uint64 `json:"timestamp"`
}
