package model

import "net/http"

// Receipt represents a receipt with necessary fields
type Receipt struct {
	Retailer     string `json:"retailer" validate:"required,printascii"`
	PurchaseDate string `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string `json:"purchaseTime" validate:"required,datetime=15:04"`
	Items        []Item `json:"items" validate:"required,dive"`
	Total        string `json:"total" validate:"required,price"`
}

// Item represents an item on a receipt
type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,printascii"`
	Price            string `json:"price" validate:"required,price"`
}

// Bind
func (r *Receipt) Bind(req *http.Request) error {
	return nil
}
