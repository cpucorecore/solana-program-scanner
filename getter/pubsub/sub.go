package pubsub

type SlotBeginSub interface {
	PubSlotBegin(slot uint64)
}

type SlotDoneSub interface {
	PubSlotDone(slot uint64)
}

type SlotErrSub interface {
	PubSlotErr(slot uint64, errCode int)
}

type GetMarketErrSub interface {
	PubGetMarketErr(marketAddr string, errCode int)
}
