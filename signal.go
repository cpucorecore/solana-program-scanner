package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"solana-program-scanner/global_stop"
	"solana-program-scanner/log"
)

func watchSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Logger.Warn("Receive system signal, do global stop", zap.Any("signal", sig))
		global_stop.Stop()
	}()
}
