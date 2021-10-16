package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	_interface "github.com/muhammadisa/bareksanews/service/interface"
	"google.golang.org/protobuf/types/known/emptypb"
)

func makeAddTagEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.AddTag(ctx, request.(*pb.Tag))
		return res, err
	}
}

func (e BareksaNewsEndpoint) AddTag(ctx context.Context, req *pb.Tag) (*emptypb.Empty, error) {
	_, err := e.AddTagEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeEditTagEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.EditTag(ctx, request.(*pb.Tag))
		return res, err
	}
}

func (e BareksaNewsEndpoint) EditTopics(ctx context.Context, req *pb.Tag) (*emptypb.Empty, error) {
	_, err := e.EditTagEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeDeleteTagEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.DeleteTag(ctx, request.(*pb.Select))
		return res, err
	}
}

func (e BareksaNewsEndpoint) DeleteTag(ctx context.Context, req *pb.Select) (*emptypb.Empty, error) {
	_, err := e.DeleteTagEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeGetTagsEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.GetTags(ctx, &emptypb.Empty{})
		return res, err
	}
}

func (e BareksaNewsEndpoint) GetTags(ctx context.Context, req *emptypb.Empty) (*pb.Tags, error) {
	res, err := e.DeleteTagEndpoint(ctx, req)
	if err != nil {
		return &pb.Tags{}, err
	}
	return res.(*pb.Tags), nil
}
