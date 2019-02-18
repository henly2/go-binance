package binance

type DepthLevelEvent struct {
	UpdateID int
	OrderBook
}

type DepthStreamEvent struct {
	Stream   string
	UpdateID int
	OrderBook
}

type DepthWebsocketRequestLevel struct {
	Symbol string
	Level  int
}

func (b *binance) DepthWebsocketLevel(dwr DepthWebsocketRequestLevel) (chan *DepthLevelEvent, chan struct{}, error) {
	return b.Service.DepthWebsocketLevel(dwr)
}

type DepthWebsocketRequestStream struct {
	Streams []string
}

func (b *binance) DepthWebsocketStream(dwr DepthWebsocketRequestStream) (chan *DepthStreamEvent, chan struct{}, error) {
	return b.Service.DepthWebsocketStream(dwr)
}
