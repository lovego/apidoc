package router

import (
	"reflect"
	"strings"
)

// T set router Title.
func (r *R) T(t string) *R {
	r.Info.Title = t
	return r
}

// D set router descriptions.
func (r *R) D(d string) *R {
	r.Info.Desc = d
	return r
}

func (r *R) ContentType(s string) *R {
	r.Info.ReqContentType = s
	return r
}

func (r *R) AddReqRes(d string) *R {
	r.Info.Desc = d
	return r
}

func (r *R) Doc(t string, reg, query string, req, res interface{}) *R {
	if (req != nil && reflect.TypeOf(req).Kind() != reflect.Ptr) ||
		(res != nil && reflect.TypeOf(res).Kind() != reflect.Ptr) {
		panic(`Doc need pointer`)
	}
	r.Info.Title = t
	//r.RegComments = parseFieldCommentPair(reg)
	//r.QueryComments = parseFieldCommentPair(query)
	//r.ReqBody = Req
	//r.ResBody = Res
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
