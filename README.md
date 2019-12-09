# apidoc
Auto generate api docs from goa routers.

## Usage
### Routers
```
// Replace *goa.RouterGroup with *router.R
rootRouter := router.NewRoot(&goa.New().RouterGroup)
```

### Router docs

- Prepare request or response body structs.
```
type Sample struct {
    // JSON tag "doc" or "comment" will be parsed for fields comments
    // and if starts with '*', this field required.
	Name string `doc:"*名称"`
	Age  string `json:"age" comment:"年龄"`
}
```
- Add api docs while writing router.

```
// One line api doc.
router.PostX(`/(\d+)`, func(c *goa.Context) {
    s := helpers.GetSession(c)
    ...
}).Doc(`标题收货`, `orderId:采购订单ID`, `queryArg1:请求Query参数1`, &detail.DetailRes, &detail.DetailRes{})

// More api doc description
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

... // setup r (*router.R)

GenDocs(r, `apidocs`)
```