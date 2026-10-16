package entity

import "time"

type TimeAttackRankingCfg struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type PlayerGradeCfg struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	SegaID string  `json:"sega_id"`
	Emoji  *string `json:"emoji,omitempty"`
}

type TimeAttackRankingMetadata struct {
	ID           string    `json:"id"`
	CourseID     string    `json:"course_id"`
	RequiredTime time.Time `json:"required_time"`
	RankID       string    `json:"rank_id"`
	RankName     string    `json:"rank_name"`
}
