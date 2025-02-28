package pagination

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	pageDefault     = 1
	pageSizeDefault = 10
	maxPageSize     = 100
)

func GetPageInfo(c echo.Context) (page int, pageSize int) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = pageDefault
	}

	pageSize, err = strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = pageSizeDefault
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}
