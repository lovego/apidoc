package apidoc

import (
	"github.com/lovego/goa"
)

type fieldCommentPair struct {
	Field   string
	Comment string
}

type errRes struct {
	Code    string      `json:"code" c:"ok 表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"`
}

type R struct {
	path   string
	method string

	title string
	desc  string // 描述

	regComments    []fieldCommentPair
	queryComments  []fieldCommentPair
	reqContentType string
	req            interface{}
	res            interface{}
	errors         []errRes

	isEntry     bool // 是否 api 接口
	RouterGroup *goa.RouterGroup
	nodes       []*R
}

func New(r *goa.RouterGroup, path string) *R {
	return &R{
		path:        path,
		RouterGroup: r,
		nodes:       make([]*R, 0),
	}
}

func NewEntry(r *goa.RouterGroup, path string) *R {
	entry := New(r, path)
	entry.isEntry = true
	return entry
}

func (r *R) Group(path string) *R {
	group := r.RouterGroup.Group(path)
	child := New(group, path)
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) GetX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.GetX(path, handlerFunc), path)
	child.method = `GET`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) Get(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Get(path, handlerFunc), path)
	child.method = `GET`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) PostX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PostX(path, handlerFunc), path)
	child.method = `POST`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) Post(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Post(path, handlerFunc), path)
	child.method = `POST`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) PutX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PutX(path, handlerFunc), path)
	child.method = `PUT`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) Put(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Put(path, handlerFunc), path)
	child.method = `PUT`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) PatchX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.PatchX(path, handlerFunc), path)
	child.method = `PATCH`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) Patch(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Patch(path, handlerFunc), path)
	child.method = `PATCH`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) DeleteX(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.DeleteX(path, handlerFunc), path)
	child.method = `DELETE`
	r.nodes = append(r.nodes, child)
	return child
}

func (r *R) Delete(path string, handlerFunc goa.HandlerFunc) *R {
	child := NewEntry(r.RouterGroup.Delete(path, handlerFunc), path)
	child.method = `DELETE`
	r.nodes = append(r.nodes, child)
	return child
}
