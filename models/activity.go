package models

// ActivityMember .
type ActivityMember struct {
	Name      string `json:"name"`
	Guild     string `json:"guild"`
	Timestamp int64  `json:"timestamp"`
}

// ActivityGuildInput input for get path parameters
type ActivityGuildInput struct {
	Guild     string `json:"guild"`
	TimeStart int64  `json:"timeStart"`
	TimeEnd   int64  `json:"timeEnd"`
}

// ActivityGuildResult for output returns [timestamp, count] pairs
type ActivityGuildResult struct {
	Guild string  `json:"guild"`
	Pairs []int64 `json:"pairs"`
}

// ActivityMemberResult outputs times spotted online
type ActivityMemberResult struct {
	Guild string  `json:"guild"`
	Times []int64 `json:"times"`
}

// ActivityMemberInput input for activity member
type ActivityMemberInput struct {
	Name      string `json:"name"`
	TimeStart int64  `json:"timeStart"`
	TimeEnd   int64  `json:"timeEnd"`
}
