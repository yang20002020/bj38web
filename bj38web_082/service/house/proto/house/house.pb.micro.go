// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/house/house.proto

package house

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	"github.com/asim/go-micro/v3/api"
	client  "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for House service

func NewHouseEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for House service

type HouseService interface {
	PubHouse(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	GetHouseInfo(ctx context.Context, in *GetReq, opts ...client.CallOption) (*GetResp, error)
}

type houseService struct {
	c    client.Client
	name string
}

func NewHouseService(name string, c client.Client) HouseService {
	return &houseService{
		c:    c,
		name: name,
	}
}

func (c *houseService) PubHouse(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "House.PubHouse", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) GetHouseInfo(ctx context.Context, in *GetReq, opts ...client.CallOption) (*GetResp, error) {
	req := c.c.NewRequest(c.name, "House.GetHouseInfo", in)
	out := new(GetResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for House service

type HouseHandler interface {
	PubHouse(context.Context, *Request, *Response) error
	GetHouseInfo(context.Context, *GetReq, *GetResp) error
}

func RegisterHouseHandler(s server.Server, hdlr HouseHandler, opts ...server.HandlerOption) error {
	type house interface {
		PubHouse(ctx context.Context, in *Request, out *Response) error
		GetHouseInfo(ctx context.Context, in *GetReq, out *GetResp) error
	}
	type House struct {
		house
	}
	h := &houseHandler{hdlr}
	return s.Handle(s.NewHandler(&House{h}, opts...))
}

type houseHandler struct {
	HouseHandler
}

func (h *houseHandler) PubHouse(ctx context.Context, in *Request, out *Response) error {
	return h.HouseHandler.PubHouse(ctx, in, out)
}

func (h *houseHandler) GetHouseInfo(ctx context.Context, in *GetReq, out *GetResp) error {
	return h.HouseHandler.GetHouseInfo(ctx, in, out)
}
