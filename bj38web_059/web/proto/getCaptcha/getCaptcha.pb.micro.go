// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/getCaptcha/getCaptcha.proto

package getCaptcha

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

// Api Endpoints for GetCaptcha service

func NewGetCaptchaEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for GetCaptcha service

type GetCaptchaService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type getCaptchaService struct {
	c    client.Client
	name string
}

func NewGetCaptchaService(name string, c client.Client) GetCaptchaService {
	return &getCaptchaService{
		c:    c,
		name: name,
	}
}

func (c *getCaptchaService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "GetCaptcha.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GetCaptcha service

type GetCaptchaHandler interface {
	Call(context.Context, *Request, *Response) error
}

func RegisterGetCaptchaHandler(s server.Server, hdlr GetCaptchaHandler, opts ...server.HandlerOption) error {
	type getCaptcha interface {
		Call(ctx context.Context, in *Request, out *Response) error
	}
	type GetCaptcha struct {
		getCaptcha
	}
	h := &getCaptchaHandler{hdlr}
	return s.Handle(s.NewHandler(&GetCaptcha{h}, opts...))
}

type getCaptchaHandler struct {
	GetCaptchaHandler
}

func (h *getCaptchaHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.GetCaptchaHandler.Call(ctx, in, out)
}
