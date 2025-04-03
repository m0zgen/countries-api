package utils

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"math"
	"strconv"
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

// Paginate оборачивает GORM-запрос с лимитом, оффсетом и возвратом мета-данных
func Paginate[T any](c fiber.Ctx, db *gorm.DB, model *T, where func(tx *gorm.DB) *gorm.DB) (*PaginationResult[T], error) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
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
	if err := query.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}

	return &PaginationResult[T]{
		Data:    items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   int(math.Ceil(float64(total) / float64(limit))),
		HasNext: int64(page*limit) < total,
		HasPrev: page > 1,
	}, nil
}
