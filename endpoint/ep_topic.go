package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	_interface "github.com/muhammadisa/bareksanews/service/interface"
	"google.golang.org/protobuf/types/known/emptypb"
)

func makeAddTopicEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.AddTopic(ctx, request.(*pb.Topic))
		return res, err
	}
}

func (e BareksaNewsEndpoint) AddTopic(ctx context.Context, req *pb.Topic) (*emptypb.Empty, error) {
	_, err := e.AddTopicEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeEditTopicEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.EditTopic(ctx, request.(*pb.Topic))
		return res, err
	}
}

func (e BareksaNewsEndpoint) EditTopic(ctx context.Context, req *pb.Topic) (*emptypb.Empty, error) {
	_, err := e.EditTopicEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeDeleteTopicEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.DeleteTopic(ctx, request.(*pb.Select))
		return res, err
	}
}

func (e BareksaNewsEndpoint) DeleteTopic(ctx context.Context, req *pb.Select) (*emptypb.Empty, error) {
	_, err := e.DeleteTopicEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeGetTopicsEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.GetTopics(ctx, &emptypb.Empty{})
		return res, err
	}
}

func (e BareksaNewsEndpoint) GetTopics(ctx context.Context, req *emptypb.Empty) (*pb.Topics, error) {
	res, err := e.DeleteTopicEndpoint(ctx, req)
	if err != nil {
		return &pb.Topics{}, err
	}
	return res.(*pb.Topics), nil
}
