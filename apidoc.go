package apidoc

import (
	"fmt"
	"github.com/lovego/apidoc/router"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lovego/apidoc/defaults"
	"github.com/lovego/apidoc/json_doc"
)

type ResBodyTpl struct {
	Code    string      `json:"code" c:"ok 表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"`
}

var BaseRes = ResBodyTpl{Code: "ok", Message: "success"}

type Doc struct {
	Indexes  []string
	Contents []string
}

func GenDocs(r *router.R, dirPath string) {
	l1 := make([]string, 0)
	for i := range r.Nodes {
		o := r.Nodes[i]
		workDir := path.Join(dirPath, o.Info.Path)
		if err := os.MkdirAll(workDir, 0755); err != nil {
			log.Panic(err)
		}
		l1 = append(l1, fmt.Sprintf(`[%s](%s)`, getTitle(r), workDir))

		l2 := make([]string, 0)
		for _, obj := range o.Nodes {
			NewDoc(obj, o.Info.Path).Create(workDir, obj.Info.Path)
			l2 = append(l2, fmt.Sprintf(`[%s](%s)`, getTitle(r), path.Join(workDir, obj.Info.Path+`.md`)))
		}

		buf := []byte(strings.Join(l2, "\n"))
		if err := ioutil.WriteFile(
			filepath.Join(workDir, "README.md"), buf, 0666,
		); err != nil {
			log.Panic(err)
		}
	}

	buf := []byte(strings.Join(l1, "\n"))
	if err := ioutil.WriteFile(
		filepath.Join(dirPath, "README.md"), buf, 0666,
	); err != nil {
		log.Panic(err)
	}
}

func NewDoc(r *router.R, basePath string) *Doc {
	d := &Doc{
		Indexes:  make([]string, 0),
		Contents: make([]string, 0),
	}
	d.Indexes = append(d.Indexes, fmt.Sprintf(`# %s <a name="index"></a>`, getTitle(r)))

	for i := range r.Nodes {
		d.Parse(r.Nodes[i], basePath, 1)
	}
	return d
}

func (d *Doc) Create(dir, name string) {
	docs := append(append(d.Indexes, ``), d.Contents...)
	buf := []byte(strings.Join(docs, "\n"))
	if err := ioutil.WriteFile(
		filepath.Join(dir, name+".md"), buf, 0666,
	); err != nil {
		log.Panic(err)
	}
}
func merge(r *router.R) {
	// merge same path group
	if r == nil || r.Info.IsEntry {
		return
	}
	path2Node := make(map[string]*router.R)
	for i := range r.Nodes {
		n := r.Nodes[i]
		if n.Info.IsEntry {
			continue
		}
		if path2Node[n.Info.Path] == nil {
			path2Node[n.Info.Path] = n
		} else {
			path2Node[n.Info.Path].Nodes = append(path2Node[n.Info.Path].Nodes, n.Nodes...)
			if n.Info.Title != `` {
				path2Node[n.Info.Path].Info.Title = n.Info.Title
			}
			r.Nodes[i] = nil
		}
		merge(n)
	}

	nodes := make([]*router.R, 0)
	for i := range r.Nodes {
		if r.Nodes[i] != nil {
			nodes = append(nodes, r.Nodes[i])
		}
	}
	r.Nodes = nodes
}

func (d *Doc) Parse(r *router.R, basePath string, level int) {
	merge(r)
	basePath += r.Info.Path
	if !r.Info.IsEntry && r.Info.Path != `` && level < 3 {
		anchorName := anchorNameReg.ReplaceAllStringFunc(basePath, func(s string) string {
			res := `-`
			return res
		})

		idx := strings.Repeat("\t", level-1) + `- `
		idx += `[` + getTitle(r) + `](#` + anchorName + `)`
		d.Indexes = append(d.Indexes, idx)
		content := "\n" + strings.Repeat("#", level+1) + ` `
		content += getTitle(r) + `<a name="` + anchorName + `"></a>`
		d.Contents = append(d.Contents, content)
		level += 1
	}
	if r.Info.IsEntry {
		idx, c := parseRouterDoc(r, basePath, level)
		d.Indexes = append(d.Indexes, idx)
		d.Contents = append(d.Contents, c)
	}
	for i := range r.Nodes {
		d.Parse(r.Nodes[i], basePath, level)
	}
}

var anchorNameReg = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func parseRouterDoc(r *router.R, path string, level int) (idx, content string) {
	docs := make([]string, 0)
	anchorName := r.Info.Method + anchorNameReg.ReplaceAllStringFunc(path, func(s string) string {
		res := `-`
		return res
	})

	// title
	title := `#### ` + getTitle(r) + `<a name="` + anchorName + `"></a> [返回目录](#index)`
	docs = append(docs, title)

	// description
	if r.Info.Desc != `` {
		docs = append(docs, r.Info.Desc)
	}

	// URL
	reqUrl := `##### ` + r.Info.Method + ` ` + path
	docs = append(docs, reqUrl)

	// RegComments
	if len(r.Info.RegComments) > 0 {
		docs = append(docs, "\n"+`##### 正则参数说明`)
		for _, o := range r.Info.RegComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	// QueryComments
	if len(r.Info.QueryComments) > 0 {
		docs = append(docs, "\n"+`##### Query 参数说明`)
		for _, o := range r.Info.QueryComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	// req
	if r.Info.Req != nil {
		docs = append(docs, "\n"+`##### 请求体说明(`+r.Info.ReqContentType+`)`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(defaults.Set(r.Info.Req)))
		docs = append(docs, "```")
	}

	// SucRes
	if r.Info.SucRes != nil {
		res := BaseRes
		res.Data = defaults.Set(r.Info.SucRes)
		docs = append(docs, "\n"+`##### 返回体说明`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(&res))
		docs = append(docs, "```")
	}
	// error responses
	if len(r.Info.ErrRes) > 0 {
		docs = append(docs, "\n"+`##### 返回错误说明`)
		for i := range r.Info.ErrRes {
			o := &r.Info.ErrRes[i]
			o.Data = defaults.Set(o.Data)
			docs = append(docs, "\n"+`- 错误码：`+o.Code)
			docs = append(docs, "```json5")
			docs = append(docs, parseJsonDoc(o))
			docs = append(docs, "```")
		}
	}
	// index
	idx = strings.Repeat("\t", level-1) + `- `
	idx += `[` + getTitle(r) + `](#` + anchorName + `)`
	content = strings.Join(docs, "\n")
	return
}

func parseJsonDoc(v interface{}) string {
	data, err := json_doc.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	list := strings.Split(string(data), "\n")

	r := regexp.MustCompile(`@@@([\s\S]*)":`)
	for i := range list {
		res := r.FindAllStringSubmatch(list[i], -1)
		if len(res) > 0 {
			str := r.ReplaceAllString(list[i], `":`)
			str += ` // ` + res[0][1]
			list[i] = str
		}
	}
	return strings.Join(list, "\n")
}

// getTitle return router title, if is empty return "无标题"
func getTitle(r *router.R) string {
	t := r.Info.Title
	if t == `` {
		t = `无标题`
	}
	return t
}
