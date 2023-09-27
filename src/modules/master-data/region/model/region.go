package model

type (
	RegProvince struct {
		ID   uint   `json:"id" gorm:"primaryKey"`
		Name string `json:"name" gorm:"not null"`
	}

	FilterProvince struct {
		ID   string
		Name string
	}

	ListProvinceResponse struct {
		Provinces []RegProvince `json:"provinces"`
		TotalPage int64         `json:"total_page"`
		TotalItem int64         `json:"total_item"`
	}

	RegCity struct {
		ID         uint   `json:"id" gorm:"primaryKey"`
		ProvinceID uint   `json:"province_id" gorm:"not null"`
		Name       string `json:"name" gorm:"not null"`
	}

	FilterCity struct {
		ID         string
		ProvinceID string
		Name       string
	}

	ListCityResponse struct {
		Cities    []RegCity `json:"cities"`
		TotalPage int64     `json:"total_page"`
		TotalItem int64     `json:"total_item"`
	}

	RegDistrict struct {
		ID     uint   `json:"id" gorm:"primaryKey"`
		CityID uint   `json:"city_id" gorm:"not null"`
		Name   string `json:"name" gorm:"not null"`
	}

	FilterDistrict struct {
		ID     string
		CityID string
		Name   string
	}

	ListDistrictResponse struct {
		Districts []RegDistrict `json:"districts"`
		TotalPage int64         `json:"total_page"`
		TotalItem int64         `json:"total_item"`
	}

	RegVillage struct {
		ID         uint   `json:"id" gorm:"primaryKey"`
		DistrictID uint   `json:"district_id" gorm:"not null"`
		Name       string `json:"name" gorm:"not null"`
	}

	FilterVillage struct {
		ID         string
		DistrictID string
		Name       string
	}

	ListVillageResponse struct {
		Villages  []RegVillage `json:"villages"`
		TotalPage int64        `json:"total_page"`
		TotalItem int64        `json:"total_item"`
	}
)
