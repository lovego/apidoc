package apidoc

import (
	"github.com/lovego/apidoc/router"
	"github.com/lovego/goa"
)

type req struct {
	Name string `json:"name" c:"名称"`
	Req  string `json:"req" c:"请求参数"`
}

type res struct {
	Name string `json:"name" c:"名称"`
	Age  string `json:"age" c:"年龄"`
	Res  string `json:"res" c:"返回信息"`
}

type errorRes struct {
	Name   string `json:"name" c:"名称"`
	Age    string `json:"age" c:"年龄"`
	ErrRes string `json:"errRes" c:"错误返回信息"`
}

func ExampleGenDocs() {

	goaRouter := goa.New()
	r0 := router.New(&goaRouter.RouterGroup, ``)
	r := r0.Group(`/purchases`).Title(`采购`)
	g1 := r.Group(`/order`).Title(`订单`)
	g2 := r.Group(`/arrival`).Title(`到货单`)

	g1.GetX(`/users`, func(c *goa.Context) {}).Doc(`用户`, ``, ``, nil, nil)
	g2.Post(`/book`, func(c *goa.Context) {}).Title(`订餐`).
		Desc(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息`).
		Req(&req{}).
		Res(&res{}).
		AddErrRes(`something-wrong`, `some thing wrong`, &errorRes{}).
		AddErrRes(`something-wrong2`, `some thing wrong2`, &errorRes{})

	g1.Group(`/child`).Title(`1111`)

	GenDocs(r0, `docs`)

	// Output:
}
