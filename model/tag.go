package model

import "time"

type Tag struct {
	ID               string
	Tag              string
	CreatedAt        int64
	UpdatedAt        int64
	Created, Updated time.Time
}

func (tag *Tag) UseUnixTimeStamp() {
	tag.CreatedAt = tag.Created.Unix()
	tag.UpdatedAt = tag.Updated.Unix()
}
