package main

import (
	"context"

	pb "github.com/nhite/pb-backend"
)

type backend struct{}

func (b *backend) Store(context.Context, *pb.Element) (*pb.Error, error) {
	return nil, nil
}
func (b *backend) Fetch(context.Context, *pb.ElementID) (*pb.Element, error) {
	return nil, nil
}
func (b *backend) List(context.Context, *pb.Pagination) (*pb.Elements, error) {
	return nil, nil
}
