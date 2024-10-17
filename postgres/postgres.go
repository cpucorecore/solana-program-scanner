package postgres

//
//import (
//	"fmt"
//
//	"github.com/lib/pq"
//	"xorm.io/xorm"
//)
//
//const ErrCodeUniqueConstrain = "23505"
//const driverName = "postgres"
//
//var engine *xorm.Engine
//
//func Init() {
//	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
//		config.G.Postgres.User,
//		config.G.Postgres.Passwd,
//		config.G.Postgres.Host,
//		config.G.Postgres.Port,
//		config.G.Postgres.Db)
//
//	var err error
//	engine, err = xorm.NewEngine(driverName, dataSource)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func InsertOperations(opCh chan types.Operation) {
//	for op := range opCh {
//		txCoordinate := types.TxCoordinate(op.BlockHeight, op.TxIndex, op.TxHash)
//		log.Logger.Info(fmt.Sprintf("**insert operation %s begin**", txCoordinate))
//
//		_, err := engine.Insert(op)
//		if err != nil {
//			log.Logger.Error(fmt.Sprintf("insert operation %s-[%s] err: %v", txCoordinate, op.String(), err))
//
//			if pgErr, ok := err.(*pq.Error); ok && string(pgErr.Code) == ErrCodeUniqueConstrain {
//				log.Logger.Info(fmt.Sprintf("insert operation failed for duplicated txId %s", op.TxHash))
//				continue
//			} else {
//				engine.Close()
//				panic(fmt.Sprintf("!!insert operation unknown postgres err %v!!", err))
//			}
//		}
//
//		log.Logger.Info(fmt.Sprintf("**insert operation %s succeed**", txCoordinate))
//	}
//}
//
//func InsertTransferOuts(transferOutCh chan types.TransferOut) {
//	for transferOut := range transferOutCh {
//		txCoordinate := types.TxCoordinate(transferOut.BlockHeight, transferOut.TxIndex, transferOut.TxHash)
//		log.Logger.Info(fmt.Sprintf("**insert TransferOut %s begin**", txCoordinate))
//
//		_, err := engine.Insert(transferOut)
//		if err != nil {
//			log.Logger.Error(fmt.Sprintf("insert TransferOut %s-[%s] err: %v", txCoordinate, transferOut.String(), err))
//
//			if pgErr, ok := err.(*pq.Error); ok && string(pgErr.Code) == ErrCodeUniqueConstrain {
//				log.Logger.Info(fmt.Sprintf("insert TransferOut failed for duplicated txId %s", transferOut.TxHash))
//				continue
//			} else {
//				engine.Close()
//				panic(fmt.Sprintf("!!insert TransferOut unknown postgres err %v!!", err))
//			}
//		}
//
//		log.Logger.Info(fmt.Sprintf("**insert TransferOut %s succeed**", txCoordinate))
//	}
//}
//
//func InsertBatchTransferOuts(transferOutCh chan types.BatchTransferOut) {
//	for transferOut := range transferOutCh {
//		txCoordinate := types.TxCoordinate(transferOut.BlockHeight, transferOut.TxIndex, transferOut.TxHash)
//		log.Logger.Info(fmt.Sprintf("**insert TransferOut %s begin**", txCoordinate))
//
//		_, err := engine.Insert(transferOut)
//		if err != nil {
//			log.Logger.Error(fmt.Sprintf("insert BatchTransferOut %s-[%s] err: %v", txCoordinate, transferOut.String(), err))
//
//			if pgErr, ok := err.(*pq.Error); ok && string(pgErr.Code) == ErrCodeUniqueConstrain {
//				log.Logger.Info(fmt.Sprintf("insert BatchTransferOut failed for duplicated txId %s", transferOut.TxHash))
//				continue
//			} else {
//				engine.Close()
//				panic(fmt.Sprintf("!!insert BatchTransferOut unknown postgres err %v!!", err))
//			}
//		}
//
//		log.Logger.Info(fmt.Sprintf("**insert BatchTransferOut %s succeed**", txCoordinate))
//	}
//}
//
//func Sync() error {
//	err := engine.Sync(types.Operation{})
//	if err != nil {
//		return nil
//	}
//
//	err = engine.Sync(types.TransferOut{})
//	if err != nil {
//		return nil
//	}
//
//	err = engine.Sync(types.BatchTransferOut{})
//	if err != nil {
//		return nil
//	}
//
//	return nil
//}
