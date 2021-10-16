package model

import "time"

type Topic struct {
	ID               string
	Title            string
	Headline         string
	CreatedAt        int64
	UpdatedAt        int64
	Created, Updated time.Time
}

func (topic *Topic) UseUnixTimeStamp() {
	topic.CreatedAt = topic.Created.Unix()
	topic.UpdatedAt = topic.Updated.Unix()
}
