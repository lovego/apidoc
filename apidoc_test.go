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
	Name string `json:"name" doc:"名称"`
	Age  string `json:"age" doc:"年龄"`
	Res  string `json:"res" doc:"返回信息"`
}

type errorRes struct {
	Name   string `json:"name" c:"*名称"`
	Age    string `json:"age" c:"*年龄"`
	ErrRes string `json:"errRes" c:"错误返回信息"`
}

func ExampleGenDocs() {

	goaRouter := goa.New()
	rootRouter := router.New(&goaRouter.RouterGroup, `/root`)
	purchaseRouter := rootRouter.Group(`/purchases`).Title(`采购`)
	arlRouter := purchaseRouter.Group(`/arrival`).Title(`到货单`)

	arlRouter.Post(`/book`, func(c *goa.Context) {}).Title(`订餐`).
		Desc(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息`).
		Regex(`id:公司ID`).
		Query(`qid:公司QID;qid2:公司`).
		Req(&req{}).
		Res(&res{}).
		AddErrRes(`something-wrong`, `some thing wrong`, &errorRes{}).
		AddErrRes(`something-wrong2`, `some thing wrong2`, &errorRes{})

	saleRouter := rootRouter.Group(`/sales`).Title(`销售`)
	saleOrderRouter := saleRouter.Group(`/order`).Title(`订单`)
	saleOrderRouter.GetX(`/detail/(\d+)`, func(c *goa.Context) {

	}).Doc(`订单详情`, `ID:订单ID`, `name:用户名`, &req{}, &res{})

	GenDocs(rootRouter, `docs`)

	// Output:
}
