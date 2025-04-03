package utils

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type PaginationResult[T any] struct {
	Data    []T   `json:"data"`
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	Total   int64 `json:"total"`
	Pages   int   `json:"pages"`
	HasNext bool  `json:"hasNext"`
	HasPrev bool  `json:"hasPrev"`
}

func Paginate[T any](c fiber.Ctx, db *gorm.DB, model *T, where func(tx *gorm.DB) *gorm.DB, allowedSorts ...string) (*PaginationResult[T], error) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	sort := c.Query("sort", "id")
	order := strings.ToUpper(c.Query("order", "ASC"))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	// Проверка допустимых полей сортировки (если переданы)
	if len(allowedSorts) > 0 {
		allowed := false
		for _, field := range allowedSorts {
			if field == sort {
				allowed = true
				break
			}
		}
		if !allowed {
			sort = allowedSorts[0] // fallback на первое поле
		}
	}

	offset := (page - 1) * limit

	query := db.Model(model)
	if where != nil {
		query = where(query)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []T
	if err := query.Order(fmt.Sprintf("%s %s", sort, order)).Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))

	// Добавим Link-заголовки
	addLinkHeader(c, page, pages, limit, sort, order)

	return &PaginationResult[T]{
		Data:    items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   pages,
		HasNext: int64(page*limit) < total,
		HasPrev: page > 1,
	}, nil
}

// addLinkHeader добавляет Link-заголовки в стиле GitHub API
func addLinkHeader(c fiber.Ctx, currentPage, totalPages, limit int, sort, order string) {
	base := c.OriginalURL()
	u, _ := url.Parse(base)
	q := u.Query()

	links := []string{}

	makeLink := func(rel string, page int) {
		q.Set("page", strconv.Itoa(page))
		q.Set("limit", strconv.Itoa(limit))
		q.Set("sort", sort)
		q.Set("order", order)
		u.RawQuery = q.Encode()
		links = append(links, fmt.Sprintf(`<%s>; rel="%s"`, u.String(), rel))
	}

	if currentPage > 1 {
		makeLink("first", 1)
		makeLink("prev", currentPage-1)
	}
	if currentPage < totalPages {
		makeLink("next", currentPage+1)
		makeLink("last", totalPages)
	}

	if len(links) > 0 {
		c.Set("Link", strings.Join(links, ", "))
	}
}
