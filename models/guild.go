package models

// GuildMember .
type GuildMember struct {
	Name  string `json:"name"`
	Guild string `json:"guild"`
	UUID  string `json:"uuid" gorm:"varchar(36);primaryKey"`
}
