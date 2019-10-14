package router

import (
	"reflect"
	"strings"

	"github.com/lovego/goa"
)

type FieldCommentPair struct {
	Field   string
	Comment string
}

type R struct {
	Path   string
	Method string

	Title         string
	RegComments   []FieldCommentPair
	QueryComments []FieldCommentPair
	ReqBody       interface{}
	ResBody       interface{}
	IsEntry       bool // 是否 api 接口

	RouterGroup *goa.RouterGroup
	Nodes       []*R
}

func New(r *goa.RouterGroup, path string) *R {
	return &R{
		Path:          path,
		RouterGroup:   r,
		Nodes:         make([]*R, 0),
		RegComments:   make([]FieldCommentPair, 0),
		QueryComments: make([]FieldCommentPair, 0),
	}
}

func NewEntry(r *goa.RouterGroup, path string) *R {
	entry := New(r, path)
	entry.IsEntry = true
	return entry
}

func (r *R) Group(path string) *R {
	group := r.RouterGroup.Group(path)
	child := New(group, path)
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Gdoc(t string) *R {
	if r.IsEntry {
		panic(`GroupDoc need router is group.`)
	}
	r.Title = t
	return r
}

func (r *R) Doc(t string, reg, query string, req, res interface{}) *R {
	if (req != nil && reflect.TypeOf(req).Kind() != reflect.Ptr) ||
		(res != nil && reflect.TypeOf(res).Kind() != reflect.Ptr) {
		panic(`Doc need pointer`)
	}
	r.Title = t
	r.RegComments = parseFieldCommentPair(reg)
	r.QueryComments = parseFieldCommentPair(query)
	r.ReqBody = req
	r.ResBody = res
	return r
}

func (r *R) T(t string) *R {
	r.Title = t
	return r
}
func (r *R) Reg(str string) *R {
	r.RegComments = append(r.RegComments, parseFieldCommentPair(str)...)
	return r
}
func (r *R) Query(str string) *R {
	r.QueryComments = append(r.QueryComments, parseFieldCommentPair(str)...)
	return r
}
func (r *R) Req(obj interface{}) *R {
	r.ReqBody = obj
	return r
}
func (r *R) Res(obj interface{}) *R {
	r.ResBody = obj
	return r
}

func parseFieldCommentPair(src string) (list []FieldCommentPair) {
	list = make([]FieldCommentPair, 0)
	if src == `` {
		return
	}
	pairs := strings.Split(src, ";")
	for i := range pairs {
		parts := strings.Split(pairs[i], ":")
		if len(parts) > 0 {
			p := FieldCommentPair{Field: parts[0]}
			if len(parts) > 1 {
				p.Comment = parts[1]
			}
			list = append(list, p)
		}
	}
	return
}

func (r *R) GetX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.GetX(path, handlerFunc), path)
	child.Method = `GET`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Get(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Get(path, handlerFunc), path)
	child.Method = `GET`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) PostX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PostX(path, handlerFunc), path)
	child.Method = `POST`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Post(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Post(path, handlerFunc), path)
	child.Method = `POST`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) PutX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PutX(path, handlerFunc), path)
	child.Method = `PUT`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Put(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Put(path, handlerFunc), path)
	child.Method = `PUT`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) PatchX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PatchX(path, handlerFunc), path)
	child.Method = `PATCH`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Patch(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Patch(path, handlerFunc), path)
	child.Method = `PATCH`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) DeleteX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.DeleteX(path, handlerFunc), path)
	child.Method = `DELETE`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Delete(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Delete(path, handlerFunc), path)
	child.Method = `DELETE`
	r.Nodes = append(r.Nodes, child)
	return child
}
