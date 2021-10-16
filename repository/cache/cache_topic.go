package cache

import (
	"context"
	"encoding/json"

	"github.com/muhammadisa/bareksanews/constant"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

func (c *cache) ReloadTopics(ctx context.Context, topics *pb.Topics) (err error) {
	data := make(map[string]interface{})
	for _, topic := range topics.Topics {
		topicByte, err := json.Marshal(topic)
		if err != nil {
			return err
		}
		data[topic.Id] = string(topicByte)
	}
	return c.redis.HMSet(ctx, constant.Topics, data).Err()
}

func (c *cache) GetTopics(ctx context.Context) (res *pb.Topics, err error) {
	var topics pb.Topics
	topicsMap := c.redis.HGetAll(ctx, constant.Topics).Val()
	for _, v := range topicsMap {
		var topic pb.Topic
		err = json.Unmarshal([]byte(v), &topic)
		if err != nil {
			return res, err
		}
		topics.Topics = append(topics.Topics, &topic)
	}
	return &topics, nil
}

func (c *cache) UnsetTopic(ctx context.Context, id string) (err error) {
	return c.redis.HDel(ctx, constant.Topics, id).Err()
}

func (c *cache) SetTopic(ctx context.Context, topic *pb.Topic) (err error) {
	topicByte, err := json.Marshal(topic)
	if err != nil {
		return err
	}
	return c.redis.HSetNX(ctx, constant.Topics, topic.Id, string(topicByte)).Err()
}
