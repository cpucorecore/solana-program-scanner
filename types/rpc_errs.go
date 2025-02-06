package types

/*
299988448 skipped, err:{"code":-32007,"message":"Slot 299988448 was skipped, or missing due to ledger jump to recent snapshot","data":null}
269154756 skipped, err:{"code":-32009,"message":"Slot 269154756 was skipped, or missing in long-term storage","data":null}
*/
const (
	RpcErr                      = -1
	RpcErrSlotCleanup           = -32001
	RpcErrSlotNotAvailable      = -32004
	RpcErrRpsLimit              = -32005
	RpcErrSlotSkippedLedgerJump = -32007
	RpcErrSlotSkippedLongTerm   = -32009
	RpcErrSlotNotAvailable2     = -32014
)
