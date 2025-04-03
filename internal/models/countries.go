package models

type Country struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	CommonName     string `json:"common_name" gorm:"not null"`
	OfficialName   string `json:"official_name"`
	FlagPNG        string `json:"flag_png" gorm:"not null"`
	FlagSVG        string `json:"flag_svg"`
	FlagAlt        string `json:"flag_alt"`
	NativeCommon   string `json:"native_common"`
	NativeOfficial string `json:"native_official"`
}
