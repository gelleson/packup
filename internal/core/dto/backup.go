package dto

import "github.com/gelleson/packup/internal/core/model"

const (
	DefaultLimit uint = 100
	DefaultSkip  uint = 0
)

type FindSnapshotQuery struct {
	Backup uint   `json:"backup"`
	Agent  uint   `json:"agent"`
	Tag    string `json:"tag"`
	Limit  uint   `json:"limit"`
	Skip   uint   `json:"skip"`
}

func (f *FindSnapshotQuery) Init() {

	if f.Limit == 0 {
		f.Limit = DefaultLimit
	}

	if f.Skip == 0 {
		f.Skip = DefaultSkip
	}
}

type SnapshotWithTotal struct {
	Snapshots []model.Snapshot `json:"snapshots"`
	Total     int64            `json:"total"`
}
