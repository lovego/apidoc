package apidoc

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/apidoc/router"
	"github.com/lovego/jsondoc"
)

var BaseRes = router.ResBodyTpl{Code: "ok", Message: "success"}

func GenDocs(r *router.R, workDir string) {
	if err := os.RemoveAll(workDir); err != nil {
		panic(err)
	}
	//merge(r)
	genDocs(r, ``, workDir)
}

// genDocs Generate doc.
// basePath is base router path.
// dirPath is dictionary path.
func genDocs(r *router.R, basePath, workDir string) {
	basePath = basePath + r.Info.Path
	if err := os.MkdirAll(workDir, 0755); err != nil {
		log.Panic(err)
	}
	indexes := make([]string, 0)
	for i := range r.Nodes {
		child := r.Nodes[i]
		if child.Info.IsEntry {
			if child.Info.Title == `` {
				log.Println(`Warning: Title is required. API: ` + r.Info.Method + ` ` + basePath + child.Info.Path)
				continue
			}
			docStr := parseEntryDoc(child, basePath)
			buf := []byte(docStr)
			fileName := child.Info.Title + `.md`
			fullPath := filepath.Join(workDir, fileName)
			if _, err := os.Stat(fullPath); err == nil {
				panic(`Error: ` + fileName + ` is exist, are you using a existing title ?`)
			}
			if err := ioutil.WriteFile(fullPath, buf, 0666); err != nil {
				log.Panic(err)
			}
			indexes = append(indexes, `### [`+child.Info.Title+`](`+fileName+`)`)
		}

		// If child router is not an entry and title is not empty,
		// then create a sub directory.
		childDir := workDir
		if !child.Info.IsEntry && child.Info.Title != `` {
			childDir = filepath.Join(workDir, child.Info.Title)
			indexes = append(indexes, `### [`+child.Info.Title+`](`+child.Info.Title+`)`)
		}
		genDocs(child, basePath, childDir)
	}
	if len(indexes) > 0 {
		indexesBuf := []byte(strings.Join(indexes, "\n"))
		if err := ioutil.WriteFile(
			filepath.Join(workDir, `README.md`), indexesBuf, 0666,
		); err != nil {
			log.Panic(err)
		}
	}
}

// merge same path group
func merge(r *router.R) {
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

func parseEntryDoc(r *router.R, basePath string) (content string) {
	urlPath := basePath + r.Info.Path
	docs := make([]string, 0)
	// title
	title := `# ` + r.Info.Title
	docs = append(docs, title)

	// description
	if r.Info.Desc != `` {
		docs = append(docs, r.Info.Desc)
	}

	// URL
	reqUrl := `## ` + r.Info.Method + ` ` + urlPath
	docs = append(docs, reqUrl)

	// RegComments
	if len(r.Info.RegComments) > 0 {
		docs = append(docs, "\n"+`## 正则参数说明`)
		for _, o := range r.Info.RegComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	// QueryComments
	if len(r.Info.QueryComments) > 0 {
		docs = append(docs, "\n"+`## Query 参数说明`)
		for _, o := range r.Info.QueryComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}

	var hasResBody = false
	for i := range r.Info.RoundTripBodies {
		o := &r.Info.RoundTripBodies[i]
		switch o.Type {
		case router.TypeReqBody:
			docs = append(docs, "\n"+`## 请求体说明(`+r.Info.ReqContentType+`)`)
			if o.Desc != `` {
				docs = append(docs, "\n"+o.Desc)
			}
			docs = append(docs, "```json5")
			docs = append(docs, makeJsonDoc(o.Body))
			docs = append(docs, "```")
		case router.TypeResBody:
			hasResBody = true
			res := BaseRes
			if o.Body != nil {
				res.Data = o.Body
			}
			docs = append(docs, "\n"+`## 返回体说明`)
			if o.Desc != `` {
				docs = append(docs, "\n"+o.Desc)
			}
			docs = append(docs, "```json5")
			docs = append(docs, makeJsonDoc(&res))
			docs = append(docs, "```")
		case router.TypeErrResBody:
			if obj, ok := o.Body.(router.ResBodyTpl); ok {
				docs = append(docs, "\n"+`## 返回错误说明: 错误码（`+obj.Code+`）`)
				if o.Desc != `` {
					docs = append(docs, "\n"+o.Desc)
				}
				if obj.Data != nil {
					obj.Data = obj.Data
					docs = append(docs, "```json5")
					docs = append(docs, makeJsonDoc(&obj))
					docs = append(docs, "```")
				}
			} else {
				panic(`errResBody type error`)
			}
		}
	}

	if !hasResBody {
		res := BaseRes
		docs = append(docs, "\n"+`## 返回体说明`)
		docs = append(docs, "```json5")
		docs = append(docs, makeJsonDoc(&res))
		docs = append(docs, "```")
	}

	content = strings.Join(docs, "\n")
	return
}

func makeJsonDoc(v interface{}) string {
	data, err := jsondoc.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	return string(data)
}
