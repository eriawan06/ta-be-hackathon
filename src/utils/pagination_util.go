package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ParamOrder ...
type ParamOrder struct {
	Field string
	By    string
}

// PaginateQueryOffset ...
type PaginateQueryOffset struct {
	Order  *ParamOrder
	Offset int
	Limit  int
	Page   int
}

// GetParamOrder Parse the url param to get order field & order by value
func GetParamOrder(r *http.Request) (*ParamOrder, error) {
	param := strings.Split(r.URL.Query().Get("order"), ",")
	if len(param) == 2 {
		if param[1] == "" {
			param[1] = "ASC"
		}
	} else if len(param) == 1 {
		param = append(param, "ASC")
	} else if len(param) == 0 {
		param = append(param, "")
		param = append(param, "")
	} else {
		return nil, errors.New("too many order parameters")
	}

	pOrder := &ParamOrder{
		Field: param[0],
		By:    param[1],
	}

	return pOrder, nil
}

// GetPaginateQueryOffset ...
func GetPaginateQueryOffset(r *http.Request) (result *PaginateQueryOffset, err error) {
	result = &PaginateQueryOffset{}

	result.Order, err = GetParamOrder(r)
	if err != nil {
		return
	}

	result.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	result.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))

	if result.Limit < 0 {
		err = errors.New("offset must be a non-negative number")
		return
	}

	result.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if result.Page > 0 {
		result.Offset = (result.Limit * result.Page) - result.Limit
	} else {
		result.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
		if result.Offset < 0 {
			err = errors.New("offset must be a non-negative number")
			return
		}
	}

	return
}

// OffsetLimitSQL generate offset limit query
func OffsetLimitSQL(pg PaginateQueryOffset) string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", pg.Limit, pg.Offset)
}
