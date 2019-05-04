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
```
type Sample struct {

    // tags "c" or "comment" will be parsed for fields comments
	Name string `c:"name你好"`
	Age  string `json:"age" comment:"年龄"`
	
	// Map Slice Array Pointer will auto set one element for default.
	Name2Id  map[int]types.Amount
	List     *types.Basic
	PtrList  [2]types.Amount
}
```
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