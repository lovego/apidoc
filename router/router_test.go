package router

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lovego/goa"
)

func ExampleRouterTree() {
	router := goa.New()
	r := NewRoot(&router.RouterGroup)
	g1 := r.Group(`/group1`).Title(`分组1`)
	//g2 := r.Group(`/group2`).Title(`分组2`)
	g1.Get(`/users`, func(c *goa.Context) {}).Doc(`用户`, ``, ``, nil, nil)
	//g2.Post(`/book`, func(c *goa.Context) {})
	g1.Group(`/child`)

	printJson(r)

	// Output:
}

func printJson(v interface{}) {
	data, err := json.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(data))
}
