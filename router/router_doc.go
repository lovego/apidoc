package router

import (
	"reflect"
	"strings"
)

// Title set router Title.
func (r *R) Title(t string) *R {
	r.Info.Title = t
	return r
}

// Desc set router descriptions.
func (r *R) Desc(d string) *R {
	r.Info.Desc = d
	return r
}

// ContentType set request content type.
func (r *R) ContentType(s string) *R {
	r.Info.ReqContentType = s
	return r
}

// Regex set request regex parameters.
func (r *R) Regex(d string) *R {
	r.Info.RegComments = parseFieldCommentPair(d)
	return r
}

// Query set request query parameters.
func (r *R) Query(d string) *R {
	r.Info.QueryComments = parseFieldCommentPair(d)
	return r
}

// Req set request body.
func (r *R) Req(d interface{}) *R {
	if d != nil && reflect.TypeOf(d).Kind() != reflect.Ptr {
		panic(`Req need pointer`)
	}
	r.Info.Req = d
	return r
}

// Res set success response body.
func (r *R) Res(d interface{}) *R {
	if d != nil && reflect.TypeOf(d).Kind() != reflect.Ptr {
		panic(`Res need pointer`)
	}
	r.Info.SucRes = d
	return r
}

// AddErrRes add error response bodies.
func (r *R) AddErrRes(code string, msg string, data interface{}) *R {
	if data != nil && reflect.TypeOf(data).Kind() != reflect.Ptr {
		panic(`AddErrRes need pointer`)
	}
	obj := errRes{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	r.Info.ErrRes = append(r.Info.ErrRes, obj)
	return r
}

// Doc provide quick set common api docs.
func (r *R) Doc(t string, reg, query string, req, res interface{}) *R {
	if (req != nil && reflect.TypeOf(req).Kind() != reflect.Ptr) ||
		(res != nil && reflect.TypeOf(res).Kind() != reflect.Ptr) {
		panic(`Doc need pointer`)
	}
	r.Info.Title = t
	r.Info.RegComments = parseFieldCommentPair(reg)
	r.Info.QueryComments = parseFieldCommentPair(query)
	r.Info.Req = req
	r.Info.SucRes = res
	return r
}

func parseFieldCommentPair(src string) (list []fieldCommentPair) {
	list = make([]fieldCommentPair, 0)
	if src == `` {
		return
	}
	pairs := strings.Split(src, ";")
	for i := range pairs {
		parts := strings.Split(pairs[i], ":")
		if len(parts) > 0 {
			p := fieldCommentPair{Field: parts[0]}
			if len(parts) > 1 {
				p.Comment = parts[1]
			}
			list = append(list, p)
		}
	}
	return
}
