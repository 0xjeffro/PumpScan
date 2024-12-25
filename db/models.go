package db

import "time"

type Txs struct {
	Tx                   string    `gorm:"primaryKey;type:text"`
	Slot                 int64     `gorm:"column:slot;type:int8"`
	User                 string    `gorm:"column:user;type:text"`
	NUser                int8      `gorm:"column:n_user;type:int8"`
	Mint                 string    `gorm:"column:mint;type:text"`
	NMint                int8      `gorm:"column:n_mint;type:int8"`
	IsCreate             bool      `gorm:"column:is_create"`
	NBuy                 int8      `gorm:"column:n_buy;type:int8"`
	NSell                int8      `gorm:"column:n_sell;type:int8"`
	BuySOLAmt            uint64    `gorm:"column:buy_sol_amt;type:numeric"`
	BuyAmt               uint64    `gorm:"column:buy_amt;type:numeric"`
	SellSOLAmt           uint64    `gorm:"column:sell_sol_amt;type:numeric"`
	SellAmt              uint64    `gorm:"column:sell_amt;type:numeric"`
	JitotipAmt           uint64    `gorm:"column:jitotip_amt;type:numeric"`
	VirtualSolReserves   uint64    `gorm:"column:virtual_sol_reserves;type:numeric"`
	VirtualTokenReserves uint64    `gorm:"column:virtual_token_reserves;type:numeric"`
	BlockTime            time.Time `gorm:"column:block_time;type:timestamptz"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamptz;default:now()"`
}
