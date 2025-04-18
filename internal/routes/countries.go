package routes

import (
	"countries-api/internal/database"
	"countries-api/internal/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type Country struct {
	ID             uint   `json:"id"`
	CommonName     string `json:"common_name"`
	OfficialName   string `json:"official_name"`
	FlagPNG        string `json:"flag_png"`
	FlagSVG        string `json:"flag_svg"`
	FlagAlt        string `json:"flag_alt"`
	NativeCommon   string `json:"native_common"`
	NativeOfficial string `json:"native_official"`
}

func CreateResponseCountry(countryModel Country) Country {
	return Country{
		ID:             countryModel.ID,
		CommonName:     countryModel.CommonName,
		OfficialName:   countryModel.OfficialName,
		FlagPNG:        countryModel.FlagPNG,
		FlagSVG:        countryModel.FlagSVG,
		FlagAlt:        countryModel.FlagAlt,
		NativeCommon:   countryModel.NativeCommon,
		NativeOfficial: countryModel.NativeOfficial,
	}
}

func GetCountries(c fiber.Ctx) error {
	var countries []Country

	err := database.Database.Db.Find(&countries).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	var responseCountries []Country
	for _, country := range countries {
		responseCountries = append(responseCountries, CreateResponseCountry(country))
	}

	return c.Status(fiber.StatusOK).JSON(responseCountries)
}

func GetPaginatedCountriesPlain(c fiber.Ctx) error {
	// 1. Параметры запроса
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	search := c.Query("search", "")

	// 2. Подготовка запроса
	db := database.Database.Db
	query := db.Model(&Country{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("common_name ILIKE ? OR official_name ILIKE ?", like, like)
	}

	// 3. Подсчёт общего количества
	var total int64
	query.Count(&total)

	// 4. Получение записей
	var countries []Country
	err = query.Offset(offset).Limit(limit).Find(&countries).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 5. Ответ
	return c.JSON(fiber.Map{
		"data":    countries,
		"page":    page,
		"limit":   limit,
		"total":   total,
		"pages":   int(math.Ceil(float64(total) / float64(limit))),
		"hasNext": page*limit < int(total),
		"hasPrev": page > 1,
	})
}

func GetPaginatedCountries(c fiber.Ctx) error {
	db := database.Database.Db

	result, err := utils.Paginate[Country](
		c,
		db,
		&Country{},
		func(tx *gorm.DB) *gorm.DB {
			search := c.Query("search")
			if search != "" {
				pattern := "%" + search + "%"
				//tx = tx.Where("common_name ILIKE ? OR official_name ILIKE ?", pattern, pattern)
				tx = tx.Where("LOWER(common_name) LIKE LOWER(?) OR LOWER(official_name) LIKE LOWER(?)", pattern, pattern)
			}
			return tx
		},
		// допустимые поля сортировки:
		"id", "common_name", "official_name",
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func AutocompleteCountries(c fiber.Ctx) error {
	db := database.Database.Db
	search := strings.ToLower(c.Query("q", ""))
	if search == "" {
		return c.JSON([]string{})
	}

	pattern := "%" + search + "%"
	var results []string

	err := db.Model(&Country{}).
		Where("LOWER(common_name) LIKE ?", pattern).
		Limit(10).
		Pluck("common_name", &results).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// Test ML suggestion
func SuggestCountriesML(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.JSON([]string{})
	}

	cmd := exec.Command("python3", "ml/ml_suggest.py", query)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("call error: Python: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "python call failed"})
	}

	var suggestions []string
	if err := json.Unmarshal(output, &suggestions); err != nil {
		log.Printf("JSON error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "json parse failed"})
	}

	return c.JSON(suggestions)
}
