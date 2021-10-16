package _interface

import (
	"context"

	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

type ReadWrite interface {
	WriteTag(ctx context.Context, req *pb.Tag) (*pb.Tag, error)
	ModifyTag(ctx context.Context, req *pb.Tag) (*pb.Tag, error)
	RemoveTag(ctx context.Context, req *pb.Select) error
	ReadTags(ctx context.Context) (*pb.Tags, error)

	WriteTopic(ctx context.Context, req *pb.Topic) (*pb.Topic, error)
	ModifyTopic(ctx context.Context, req *pb.Topic) (*pb.Topic, error)
	RemoveTopic(ctx context.Context, req *pb.Select) error
	ReadTopics(ctx context.Context) (*pb.Topics, error)

	WriteNews(ctx context.Context, req *pb.News) (*pb.News, error)
	ModifyNews(ctx context.Context, req *pb.News) (*pb.News, error)
	RemoveNews(ctx context.Context, req *pb.Select) error
	ReadNewses(ctx context.Context) (*pb.Newses, error)
	ReadNewsesByStatus(ctx context.Context, status int32) (*pb.Newses, error)
	ReadNewsesByTopicID(ctx context.Context, topicID string) (*pb.Newses, error)
	ReadNewsesByStatusAndTopicID(ctx context.Context, status int32, topicID string) (*pb.Newses, error)

	RemoveNewsTagsByNewsID(ctx context.Context, req *pb.Select) error
	WriteNewsTags(ctx context.Context, newsID string, tagIDs []string, new bool) error
	ReadNewsTagsTagIDAndTagByNewsID(ctx context.Context, newsID string, all bool) (res []string)
}

type Cache interface {
	SetTag(ctx context.Context, tag *pb.Tag) error
	UnsetTag(ctx context.Context, id string) error
	GetTags(ctx context.Context) (*pb.Tags, error)
	ReloadTags(ctx context.Context, tags *pb.Tags) error

	SetTopic(ctx context.Context, tag *pb.Topic) error
	UnsetTopic(ctx context.Context, id string) error
	GetTopics(ctx context.Context) (*pb.Topics, error)
	ReloadTopics(ctx context.Context, topics *pb.Topics) error

	InvalidateNewses(ctx context.Context) error
	ReloadRequired(ctx context.Context, filter string) bool
	UnsetNews(ctx context.Context, id string) error
	GetNewses(ctx context.Context) (res *pb.Newses, err error)
	ReloadNewses(ctx context.Context, newses *pb.Newses) error

	GetFilter(ctx context.Context) string
	SetFilter(ctx context.Context, filter string) error
}
