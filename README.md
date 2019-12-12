# apidoc
Auto generate api docs from goa routers.

## Usages
### Routers
Replace *goa.RouterGroup with *router.R
```
rootRouter := router.NewRoot(&goa.New().RouterGroup)
```

### Router docs

####  Prepare request or response body structs.
JSON tag "doc" or "comment" will be parsed as fields comments. if it starts with '*', then this field is required.
```
type Sample struct {
	Name string `doc:"*名称"`
	Age  string `json:"age" comment:"年龄"`
}
```
#### Add api docs while writing router.
There are two ways to add docs to a router. 

##### One line docs.

`Doc` add docs in one line for common situation. 
```
router.PostX(`/(\d+)`, func(c *goa.Context) {
    s := helpers.GetSession(c)
    ...
}).Doc(`标题收货`, `orderId:采购订单ID`, `queryArg1:请求Query参数1`, &detail.DetailRes, &detail.DetailRes{})
```
##### More docs.

The other way is to call `Title`,`Desc`,`Regex`,`Query`,`Req`,`Res`, or `ErrRes` individually.
`Req`,`Res` and `ErrRes` will append docs to an array, so it will be parsed one by one.
```
router.
    Title(`订餐`).
    Desc(`路由描述信息`).
    Regex(`id:公司ID`).
    Query(`id:公司ID;id2:公司ID2`).
    Req(`请求体描述信息`, &req{}).
    Res(`返回体描述信息`, &res{}).
    ErrRes(`错误返回体描述信息`, `something-wrong`, `some thing wrong`, &errorRes{}).
```

### Generate docs

Run the code below at anywhere you like.
```
router.ForDoc = true 

... // setup routers.

GenDocs(r, `apidocs`)
```