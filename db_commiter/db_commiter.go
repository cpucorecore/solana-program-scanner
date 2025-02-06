package db_commiter

import (
	"solana-program-scanner/types/orms"
	"sync"
)

type DBCommitter struct {
	PairCommiter  *BatchCommitter[*orms.Pair]
	TokenCommiter *BatchCommitter[*orms.Token]
}

func (c *DBCommitter) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	var wg2 sync.WaitGroup
	wg2.Add(2)
	go c.PairCommiter.Run(&wg2)
	go c.TokenCommiter.Run(&wg2)
	wg2.Wait()
}
