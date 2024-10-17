test:
	go build
	./solana-program-scanner > scanner.log 2>&1

clean:
	rm -rf solana-program-scanner scanner.log blocks.json
