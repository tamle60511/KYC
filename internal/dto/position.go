package dto

import "time"

type PositionCreate struct {
	Code   string `json:"code"`
	NameVN string `json:"name_vn"`
	NameEN string `json:"name_en"`
	NameTW string `json:"name_tw"`
	Level  int    `json:"level" validate:"min=1"`
	Desc   string `json:"desc"`
}

type PositionUpdate struct {
	NameVN *string `json:"name_vn"`
	NameEN *string `json:"name_en"`
	NameTW *string `json:"name_tw"`
	Level  *int    `json:"level"`
	Desc   *string `json:"desc"`
}

type PositionResponse struct {
	ID        uint64    `json:"id"`
	Code      string    `json:"code"`
	NameVN    string    `json:"name_vn"`
	NameEN    string    `json:"name_en"`
	NameTW    string    `json:"name_tw"`
	Level     int       `json:"level"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PositionID struct {
	ID uint64 `json:"id" uri:"id" binding:"required"`
}
