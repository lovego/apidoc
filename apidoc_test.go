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
	rootRouter := router.New(&goaRouter.RouterGroup, ``)
	firstLevel := rootRouter.Group(`/purchases`).Title(`采购`)
	secondLevel := firstLevel.Group(`/arrival`).Title(`到货单`)

	secondLevel.Post(`/book`, func(c *goa.Context) {}).Title(`订餐`).
		Desc(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息`).
		Regex(`id:公司ID`).
		Query(`qid:公司QID;qid2:公司`).
		Req(&req{}).
		Res(&res{}).
		AddErrRes(`something-wrong`, `some thing wrong`, &errorRes{}).
		AddErrRes(`something-wrong2`, `some thing wrong2`, &errorRes{})

	GenDocs(rootRouter, `docs`)

	// Output:
}
