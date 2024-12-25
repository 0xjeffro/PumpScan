package db

import (
	"testing"
	"time"
)

func TestGetConn(t *testing.T) {
	_, err := GetConn()
	if err != nil {
		t.Errorf("GetConn() error: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	db, err := GetConn()
	if err != nil {
		t.Errorf("GetConn() error: %v", err)
	}
	err = db.AutoMigrate(&Txs{})
}

func TestInsert(t *testing.T) {
	txs := Txs{
		Tx:         "test1",
		User:       "Signature",
		NUser:      1,
		BlockTime:  time.Unix(1734944288, 0),
		Slot:       0,
		Mint:       "test",
		NMint:      0,
		IsCreate:   false,
		NBuy:       0,
		NSell:      0,
		BuyAmt:     0,
		BuySOLAmt:  0,
		SellAmt:    0,
		SellSOLAmt: 0,
		JitotipAmt: 0,
	}

	err := Insert(txs)
	if err != nil {
		t.Errorf("InsertTable() error: %v", err)
	}
}
