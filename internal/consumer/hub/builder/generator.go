package builder

import (
	"context"
)

type IBuilder interface {
	ParseJson(ctx context.Context, str string) error // 解析json
	ThirdToPro(ctx context.Context) error            // 消费由【三方平台】向【线上环境】发送的数据
	ProToTest(ctx context.Context) error             // 消费由【线上环境】向【测试环境】发送的数据
	ProToPre(ctx context.Context) error              // 消费由【线上环境】向【预发环境】发送的数据
	ProToTestPass(ctx context.Context) error         // 消费由【管理后台审核通过后】向【线上环境】发送的数据
}

type Generator struct {
	reqFrom string
}

func (g *Generator) GetReqFrom() string {
	return g.reqFrom
}

// GeneratorOption Generator的初始化参数
type GeneratorOption func(*Generator)

// NewGenWithOpt ...
func NewGenWithOpt(opts ...GeneratorOption) *Generator {
	generator := &Generator{}
	for _, opt := range opts {
		opt(generator)
	}
	return generator
}

func WithReqFrom(from string) GeneratorOption {
	return func(pg *Generator) {
		pg.reqFrom = from
	}
}
