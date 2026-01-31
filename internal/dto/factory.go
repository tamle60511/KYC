package dto

type FactoryCreate struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	TaxCode string `json:"tax_code"`
}

type FactoryUpdate struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	TaxCode  *string `json:"tax_code"`
	IsActive *bool   `json:"is_active"`
}

type FactoryResponse struct {
	ID       uint64 `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	TaxCode  string `json:"tax_code"`
	IsActive bool   `json:"is_active"`
}

type FactoryID struct {
	ID uint64 `json:"id" uri:"id" binding:"required"`
}
