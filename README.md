# apidoc
Auto generate api docs from goa routers.

## Usage
### Routers
Almost the same usage as `goa`.
```
// Replace *goa.RouterGroup with *router.R
r := router.New(goaRouter.Group(path), path)

orders.Routes(r.Group(`/order`))
```

### Router docs

- Prepare request and response body structs
- Call `r.Doc()`.

```
router.PostX(`/(\d+)`, func(c *goa.Context) {
    s := helpers.GetSession(c)
    ...
}).Doc(`标题收货`, `orderId:采购订单ID`, `queryArg1:请求Query参数1`, &detail.DetailRes, &detail.DetailRes{})
```

### Generate docs

Run the code below at somewhere you like.

```
path := `/purchases`
r := purchases.Routes(goa.New(), path)

apidoc.NewDoc(r).Create(fs.SourceDir(), r.Path)
```