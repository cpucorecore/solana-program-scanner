package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blocto/solana-go-sdk/rpc"
)

type BlockProcessorFile struct {
	fd *os.File
}

var _ ObjProcessor[*rpc.GetBlock] = &BlockProcessorFile{}

func (bpf *BlockProcessorFile) id() string {
	return "BlockProcessorFile"
}

func (bpf *BlockProcessorFile) process(block *rpc.GetBlock) error {
	blockData, err := json.Marshal(block)
	if err != nil {
		Logger.Error(fmt.Sprintf(""))
		return err
	}

	_, err = bpf.fd.Write(blockData)
	return err
}

func (bpf *BlockProcessorFile) done() {
	bpf.fd.Close()
}

func NewBlockProcessorFile(name string) (bpf *BlockProcessorFile) {
	fd, err := os.Create(name)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("os.Create err:%s", err.Error()))
	}

	return &BlockProcessorFile{
		fd: fd,
	}
}
