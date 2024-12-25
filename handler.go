package main

import (
	"PumpScan/db"
	"PumpScan/parser"
	"github.com/0xjeffro/tx-parser/solana/types"
	"log"
	"slices"
	"strings"
	"time"
)

var jitoTipAddress = []string{
	"96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5",
	"HFqU5x63VTqvQss8hp11i4wVV8bD44PvwucfZ2bU7gRe",
	"Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY",
	"ADaUMid9yfUytqMBgopwjb2DTLSokTSzL1zt6iGPaS49",
	"DfXygSm4jCyNCybVYYK6DwvWqjKee8pbDmJGcLWNDXjh",
	"ADuUkR4vqLUMWXxW9gh6D6L8pMSawimctcNZ5pGwDcEt",
	"DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL",
	"3AVi9Tg9Uo68tJfuvoKvqKNWKkC5wPdSSdeBnizKZ6jT",
}

func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	return strings.Join(strs, ",")
}

func WebhookHandler(bytes []byte) {
	actions, ctx, err := parser.Parser(bytes)
	if err != nil {
		log.Println("[WebhookHandler] Parser error: ", err)
		return
	}

	if len(actions) == 0 {
		log.Println("[WebhookHandler] No actions found tx: ", ctx.RawTx.Transaction.Signatures[0])
		return
	}

	var txRow db.Txs
	txRow.Tx = ctx.RawTx.Transaction.Signatures[0]
	txRow.BlockTime = time.Unix(ctx.RawTx.BlockTime, 0)
	txRow.Slot = ctx.RawTx.Slot

	var userMap = make(map[string]bool)
	var mintMap = make(map[string]bool)

	for _, action := range actions {

		switch a := action.(type) {
		case *types.PumpFunCreateAction:
			userMap[a.Who] = true
			mintMap[a.Mint] = true
			txRow.IsCreate = true
		case *types.PumpFunAnchorSelfCPILogSwapAction:
			userMap[a.User] = true
			mintMap[a.Mint] = true
			if a.IsBuy {
				txRow.BuyAmt += a.TokenAmount
				txRow.BuySOLAmt += a.SolAmount
				txRow.NBuy++

				// VirtualTokenReserves and VirtualSolReserves may be overwritten by subsequent actions;
				// only applies when there's exactly one swap (NBuy + NSell == 1)
				txRow.VirtualTokenReserves = a.VirtualTokenReserves
				txRow.VirtualSolReserves = a.VirtualSolReserves
			} else {
				txRow.SellAmt += a.TokenAmount
				txRow.SellSOLAmt += a.SolAmount
				txRow.NSell++

				// VirtualTokenReserves and VirtualSolReserves may be overwritten by subsequent actions;
				// only applies when there's exactly one swap (NBuy + NSell == 1)
				txRow.VirtualTokenReserves = a.VirtualTokenReserves
				txRow.VirtualSolReserves = a.VirtualSolReserves
			}
		case *types.SystemProgramTransferAction:
			if slices.Contains(jitoTipAddress, a.To) {
				txRow.JitotipAmt += a.Lamports
			}
		}
	}

	numMint := 0
	var mints []string
	for k, v := range mintMap {
		if v {
			numMint++
			mints = append(mints, k)
		}
	}

	numUser := 0
	var users []string
	for k, v := range userMap {
		if v {
			numUser++
			users = append(users, k)
		}
	}

	txRow.NMint = int8(numMint)
	txRow.Mint = joinStrings(mints)

	txRow.NUser = int8(numUser)
	txRow.User = joinStrings(users)

	if numUser == 0 {
		log.Println("[WebhookHandler] 'Who' is empty, tx: ", ctx.RawTx.Transaction.Signatures[0])
		return
	} else {
		err := db.Insert(txRow)
		if err != nil {
			log.Println("[WebhookHandler] Error inserting tx: ", err)
		}
	}
}
