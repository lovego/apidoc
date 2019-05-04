// - Use reflect to concat json tag and 'comment' or 'c' tag as json name
// - then use json MarshalIndent to get json string.
// - then parse name to name part and comments part. done!

package main

// TODO
// type Sample struct {
// 	Name string `c:"nameeeeeee11你好"`
// 	Age  string `json:"age"`
// 	types.Amount
// 	Name2Id  map[int]types.Amount
// 	List     *types.Basic
// 	PtrList  [2]types.Amount
// 	PtrList2 [2]*types.Amount
// }

func main() {
	// obj := apidoc.BaseRes
	// obj.Data = detail.DetailRes{}
	// fmt.Println(parseJsonDoc(&obj))

	// p := `/purchases`
	// r := purchases.Routes(goa.New(), p)
	//
	// doc := apidoc.NewDoc()
	// doc.Parse(r, ``, 1)
	// doc.Create()
}
