package cache

import (
	"context"
	"encoding/json"

	"github.com/muhammadisa/bareksanews/constant"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

func (c *cache) ReloadTags(ctx context.Context, tags *pb.Tags) (err error) {
	const funcName = `ReloadTags`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	data := make(map[string]interface{})
	for _, tag := range tags.Tags {
		tagByte, err := json.Marshal(tag)
		if err != nil {
			return err
		}
		data[tag.Id] = string(tagByte)
	}
	return c.redis.HMSet(ctx, constant.Tags, data).Err()
}

func (c *cache) UnsetTag(ctx context.Context, id string) (err error) {
	const funcName = `UnsetTag`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	return c.redis.HDel(ctx, constant.Tags, id).Err()
}

func (c *cache) SetTag(ctx context.Context, tag *pb.Tag) (err error) {
	const funcName = `SetTag`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	tagByte, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	return c.redis.HSetNX(ctx, constant.Tags, tag.Id, string(tagByte)).Err()
}

func (c *cache) GetTags(ctx context.Context) (res *pb.Tags, err error) {
	const funcName = `GetTags`
	_, span := c.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var tags pb.Tags
	tagsMap := c.redis.HGetAll(ctx, constant.Tags).Val()
	for _, v := range tagsMap {
		var tag pb.Tag
		err = json.Unmarshal([]byte(v), &tag)
		if err != nil {
			return res, err
		}
		tags.Tags = append(tags.Tags, &tag)
	}
	return &tags, nil
}
