package apidoc

import (
	"github.com/lovego/apidoc/router"
	"github.com/lovego/goa"
)

func ExampleGenDocs() {

	goaRouter := goa.New()
	r0 := router.New(&goaRouter.RouterGroup, ``)
	r := r0.Group(`/purchases`).T(`采购`)
	g1 := r.Group(`/order`).Gdoc(`订单`)
	g2 := r.Group(`/arrival`).Gdoc(`到货单`)

	g1.GetX(`/users`, func(c *goa.Context) {}).Doc(`用户`, ``, ``, nil, nil)
	g2.Post(`/book`, func(c *goa.Context) {})

	g1.Group(`/child`)

	GenDocs(r0, `docs`)

	// Output:
}
