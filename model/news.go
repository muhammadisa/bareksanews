package model

import "time"

type News struct {
	ID               string
	TopicID          string
	Title            string
	Content          string
	Status           int32
	CreatedAt        int64
	UpdatedAt        int64
	Created, Updated time.Time
}

type NewsTag struct {
	ID               string
	NewsID           string
	TagID            string
	CreatedAt        int64
	UpdatedAt        int64
	Created, Updated time.Time
}

func (news *News) UseUnixTimeStamp() {
	news.CreatedAt = news.Created.Unix()
	news.UpdatedAt = news.Updated.Unix()
}

func (newsTag *NewsTag) UseUnixTimeStamp() {
	newsTag.CreatedAt = newsTag.Created.Unix()
	newsTag.UpdatedAt = newsTag.Updated.Unix()
}
