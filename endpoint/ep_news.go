package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	_interface "github.com/muhammadisa/bareksanews/service/interface"
	"google.golang.org/protobuf/types/known/emptypb"
)

func makeAddNewsEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.AddNews(ctx, request.(*pb.News))
		return res, err
	}
}

func (e BareksaNewsEndpoint) AddNews(ctx context.Context, req *pb.News) (*emptypb.Empty, error) {
	_, err := e.AddNewsEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeEditNewsEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.EditNews(ctx, request.(*pb.News))
		return res, err
	}
}

func (e BareksaNewsEndpoint) EditNews(ctx context.Context, req *pb.News) (*emptypb.Empty, error) {
	_, err := e.EditNewsEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeDeleteNewsEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.DeleteNews(ctx, request.(*pb.Select))
		return res, err
	}
}

func (e BareksaNewsEndpoint) DeleteNews(ctx context.Context, req *pb.Select) (*emptypb.Empty, error) {
	_, err := e.DeleteNewsEndpoint(ctx, req)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func makeGetNewsesEndpoint(usecase _interface.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.GetNewses(ctx, request.(*pb.Filters))
		return res, err
	}
}

func (e BareksaNewsEndpoint) GetNewses(ctx context.Context, req *pb.Filters) (*pb.Newses, error) {
	res, err := e.DeleteNewsEndpoint(ctx, req)
	if err != nil {
		return &pb.Newses{}, err
	}
	return res.(*pb.Newses), nil
}
