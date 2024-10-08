package main

import (
	"PumpScan/db"
	"encoding/json"
	"github.com/0xjeffro/tx-parser/solana"
	"github.com/0xjeffro/tx-parser/solana/types"
	"log"
	"time"
)

func WebhookHandler(bytes []byte) {
	var txs types.RawTxs
	err := json.Unmarshal(bytes, &txs)
	if err != nil {
		log.Println("Transaction unmarshal error: ", err)
		return
	}
	parsedData, err := solana.Parser(bytes)
	if err != nil {
		log.Println("Parser error: ", err)
		return
	}
	for _, data := range parsedData {
		blockTime := data.RawTx.BlockTime
		var nBuy, nSell int = 0, 0
		var buyAmount, sellAmount uint64 = 0, 0
		var isCreate bool = false
		var swapTokens = make(map[string]bool)
		for _, action := range data.Actions {
			switch a := action.(type) {
			case *types.PumpFunCreateAction:
				isCreate = true
				swapTokens[a.Mint] = true
			case *types.PumpFunBuyAction:
				nBuy++
				buyAmount += a.ToTokenAmount
				swapTokens[a.ToToken] = true
			case *types.PumpFunSellAction:
				nSell++
				sellAmount += a.FromTokenAmount
				swapTokens[a.FromToken] = true
			}
		}
		var count int = 0
		var mint string
		for k, v := range swapTokens {
			if v {
				mint = k
				count++
			}
		}
		log.Println("Mint: ", mint)
		// if there is more than 1 token in the swap, continue
		if count > 1 {
			log.Println("More than 1 token in the swap")
			continue
		}
		// if there is not a bundle of buy and sell, continue
		if nBuy+nSell <= 1 && isCreate == false {
			continue
		}
		row := db.PumpInsiderEvent{
			Tx:        data.RawTx.Transaction.Signatures[0],
			BlockTime: time.Unix(blockTime, 0),

			Mint:     mint,
			NBuy:     nBuy,
			NSell:    nSell,
			BuyAmt:   buyAmount,
			SellAmt:  sellAmount,
			IsCreate: isCreate,
		}
		err := db.InsertInsiderEvent(row)
		if err != nil {
			log.Println("Error inserting insider event: ", err)
		}
	}
}
