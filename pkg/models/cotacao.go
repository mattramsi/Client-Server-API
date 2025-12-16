package models

import "time"

type BidResponse struct {
	Bid string `json:"bid"`
}

type DolarResponse struct {
	USDBRL DolarInfo `json:"USDBRL"`
}

type DolarInfo struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Codein     string    `json:"codein"`
	Name       string    `json:"name"`
	High       string    `json:"high"`
	Low        string    `json:"low"`
	VarBid     string    `json:"var_bid"`
	PctChange  string    `json:"pct_change"`
	Bid        string    `json:"bid"`
	Ask        string    `json:"ask"`
	Timestamp  string    `json:"timestamp"`
	CreateDate string    `json:"create_date"`
	CreatedAt  time.Time `json:"created_at"`
}
