package db

import "log"

func Insert(txs Txs) error {
	db, err := GetConn()
	if err != nil {
		log.Panicln("[Insert] GetConn error: ", err)
		return err
	}
	err = db.Create(&txs).Error
	if err != nil {
		log.Panicln("[Insert] Create error: ", err)
		return err
	}
	return nil
}
