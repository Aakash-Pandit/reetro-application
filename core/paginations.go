package core

import (
	"net/http"
	"strconv"
)

const (
	PAGE  = 1
	LIMIT = 10
)

func Pagination(r *http.Request) (int, int) {
	page := PAGE
	limit := LIMIT

	if val, ok := r.URL.Query()["limit"]; ok {
		if l, err := strconv.Atoi(val[0]); err == nil && l > 0 {
			limit = l
		}
	}

	if val, ok := r.URL.Query()["page"]; ok {
		if p, err := strconv.Atoi(val[0]); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit
	return limit, offset
}
