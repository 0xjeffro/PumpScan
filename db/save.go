package db

import (
	"os"
	"time"
)
import supa "github.com/nedpals/supabase-go"

type PumpInsiderEvent struct {
	Tx        string    `json:"tx"`
	BlockTime time.Time `json:"block_time"`
	NBuy      int       `json:"n_buy"`
	NSell     int       `json:"n_sell"`
	BuyAmt    uint64    `json:"buy_amt"`
	SellAmt   uint64    `json:"sell_amt"`
	IsCreate  bool      `json:"is_create"`
}

func InsertInsiderEvent(row PumpInsiderEvent) error {
	supaUrl := os.Getenv("SUPA_URL")
	supaSecretKey := os.Getenv("SUPA_SECRET_KEY")
	supabase := supa.CreateClient(supaUrl, supaSecretKey)
	var results []PumpInsiderEvent
	err := supabase.DB.From("pump_scan").Insert(row).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}
