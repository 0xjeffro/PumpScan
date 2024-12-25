package parser

import (
	"encoding/json"
	"github.com/0xjeffro/tx-parser/solana"
	"github.com/0xjeffro/tx-parser/solana/programs/pumpfun"
	pumpfunParser "github.com/0xjeffro/tx-parser/solana/programs/pumpfun/parsers"
	systemParser "github.com/0xjeffro/tx-parser/solana/programs/systemProgram/parsers"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/mr-tron/base58"

	"log"
)

func Parser(bytes []byte) ([]types.Action, *types.ParsedResult, error) {
	var txs types.RawTxs
	err := json.Unmarshal(bytes, &txs)
	if err != nil {
		log.Println("[Parser] Transaction unmarshal error: ", err)
		return nil, nil, err
	}

	var parsedResult types.ParsedResult
	parsedResult.RawTx = txs[0]
	parsedResult = *solana.GetAccountList(&parsedResult)

	totalCapacity := len(parsedResult.RawTx.Transaction.Message.Instructions)
	for _, inner := range parsedResult.RawTx.Meta.InnerInstructions {
		totalCapacity += len(inner.Instructions)
	}
	allInstructions := make([]types.Instruction, 0, totalCapacity)
	allInstructions = append(allInstructions, parsedResult.RawTx.Transaction.Message.Instructions...)

	for _, innerInstruction := range parsedResult.RawTx.Meta.InnerInstructions {
		allInstructions = append(allInstructions, innerInstruction.Instructions...)
	}

	actions := make([]types.Action, 0)

	for _, instr := range allInstructions {
		programID := parsedResult.AccountList[instr.ProgramIDIndex]
		switch programID {
		case PumpFunProgramID:
			decodedData, err := base58.Decode(instr.Data)
			if err != nil {
				log.Println("[Parser] Decode error: ", err)
				continue
			}
			discriminator := *(*[16]byte)(decodedData[:16])
			mergedDiscriminator := make([]byte, 0, 16)
			mergedDiscriminator = append(mergedDiscriminator[:], pumpfun.AnchorSelfCPILogDiscriminator[:]...)
			mergedDiscriminator = append(mergedDiscriminator[:], pumpfun.AnchorSelfCPILogSwapDiscriminator[:]...)
			if discriminator == *(*[16]byte)(mergedDiscriminator[:]) {
				action, err := pumpfunParser.AnchorSelfCPILogSwapParser(decodedData)
				if err == nil {
					actions = append(actions, action)
				} else {
					log.Println("[Parser] PumpFun instruction error: ", err)
				}
			}
		case SystemProgramID:
			action, err := systemParser.InstructionRouter(&parsedResult, instr)
			if err == nil {
				actions = append(actions, action)
			} else {
				log.Println("[Parser] System instruction error: ", err)
			}
		}
	}
	return actions, &parsedResult, nil
}
