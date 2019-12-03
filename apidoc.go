package apidoc

import (
	"fmt"
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

func GenDocs(r *R, dirPath string) {
	l1 := make([]string, 0)
	for i := range r.nodes {
		o := r.nodes[i]
		workDir := path.Join(dirPath, o.path)
		if err := os.MkdirAll(workDir, 0755); err != nil {
			log.Panic(err)
		}
		l1 = append(l1, fmt.Sprintf(`[%s](%s)`, o.title, workDir))

		l2 := make([]string, 0)
		for _, obj := range o.nodes {
			NewDoc(obj, o.path).Create(workDir, obj.path)
			l2 = append(l2, fmt.Sprintf(`[%s](%s)`, obj.title, path.Join(workDir, obj.path+`.md`)))
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

func NewDoc(r *R, basePath string) *Doc {
	d := &Doc{
		Indexes:  make([]string, 0),
		Contents: make([]string, 0),
	}
	d.Indexes = append(d.Indexes, `# Base path: `+basePath)
	d.Indexes = append(d.Indexes, `# Index <a name="index"></a>`)

	d.Parse(r, basePath, 1)
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
func merge(r *R) {
	// merge same path group
	if r == nil || r.isEntry {
		return
	}
	path2Node := make(map[string]*R)
	for i := range r.nodes {
		n := r.nodes[i]
		if n.isEntry {
			continue
		}
		if path2Node[n.path] == nil {
			path2Node[n.path] = n
		} else {
			path2Node[n.path].nodes = append(path2Node[n.path].nodes, n.nodes...)
			if n.title != `` {
				path2Node[n.path].title = n.title
			}
			r.nodes[i] = nil
		}
		merge(n)
	}

	nodes := make([]*R, 0)
	for i := range r.nodes {
		if r.nodes[i] != nil {
			nodes = append(nodes, r.nodes[i])
		}
	}
	r.nodes = nodes
}

func (d *Doc) Parse(r *R, basePath string, level int) {
	merge(r)
	basePath += r.path
	if !r.isEntry && r.path != `` && level < 3 {
		idx := strings.Repeat("\t", level-1) + `- `
		idx += `[` + r.title + ` ` + r.path + `](#` + basePath + `)`
		d.Indexes = append(d.Indexes, idx)
		content := "\n" + strings.Repeat("#", level) + ` `
		content += r.title + ` ` + r.path
		d.Contents = append(d.Contents, content)
		level += 1
	}
	if len(r.nodes) == 0 {
		idx, c := parseRouterDoc(r, basePath, level)
		d.Indexes = append(d.Indexes, idx)
		d.Contents = append(d.Contents, c)
	}
	for i := range r.nodes {
		d.Parse(r.nodes[i], basePath, level)
	}
}

var anchorNameReg = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func parseRouterDoc(r *R, path string, level int) (idx, content string) {
	docs := make([]string, 0)
	name := r.method + anchorNameReg.ReplaceAllStringFunc(path, func(s string) string {
		res := `-`
		return res
	})
	title := "\n#### " + r.method + ` ` + path
	if r.title != `` {
		title += ` (` + r.title + `)`
	}
	title += `<a name="` + name + `"></a>`
	title += ` [index](#index)`
	docs = append(docs, title)

	if len(r.regComments) > 0 {
		docs = append(docs, "\n"+`##### 正则参数说明`)
		for _, o := range r.regComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}

	if len(r.queryComments) > 0 {
		docs = append(docs, "\n"+`##### Query 参数说明`)
		for _, o := range r.queryComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	if r.req != nil {
		docs = append(docs, "\n"+`##### Request Body`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(defaults.Set(r.req)))
		docs = append(docs, "```")
	}

	if r.res != nil {

		res := BaseRes
		res.Data = defaults.Set(r.res)
		docs = append(docs, "\n"+`##### Response Body`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(&res))
		docs = append(docs, "```")
	}
	docs = append(docs)
	idx = strings.Repeat("\t", level-1) + `- `
	idx += `[` + r.title + ` ` + r.method + ` ` + r.path + `](#` + name + `)`
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
