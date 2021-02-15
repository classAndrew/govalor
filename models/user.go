package models

// UserSliceXP .
type UserSliceXP struct {
	Name      string `json:"name"`
	Timestamp uint   `json:"timestamp"`
	Guild     string `json:"guild"`
	XP        int64  `json:"xp"`
}

// UserTotalXP .
type UserTotalXP struct {
	// can't have binding: "required" on any of these since usertotalxpresponse can't act as a 'union'
	Name   string `json:"name"`
	XP     int64  `json:"xp"`
	LastXP int64  `json:"lastxp"`
	Guild  string `json:"guild"`
}

// UserTotalXPResponse .
type UserTotalXPResponse struct {
	Error string      `json:"error"`
	Data  UserTotalXP `json:"data"`
}
