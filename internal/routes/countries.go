package routes

import (
	"countries-api/internal/database"
	"github.com/gofiber/fiber/v3"
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
