package _interface

import pb "github.com/muhammadisa/bareksanews/protoc/api/v1"

type Service interface {
	pb.BareksaNewsServiceServer
}
