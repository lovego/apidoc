package router

import (
	"reflect"
	"strings"
)

// Title set router Title.
func (r *R) Title(t string) *R {
	t = strings.TrimSpace(t)
	if strings.ContainsAny(t, `/`) {
		panic(`Title contains '/' : ` + t)
	}
	if t[0] == '+' || t[0] == '-' || t[0] == '.' {
		panic(`Title starts with '+-.' : ` + t)
	}
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
	r.Title(t)
	r.Regex(reg)
	r.Query(query)
	r.Req(req)
	r.Res(res)
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
