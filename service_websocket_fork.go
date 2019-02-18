package binance

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	"strings"
)

func (as *apiService) DepthWebsocketLevel(dwr DepthWebsocketRequestLevel) (chan *DepthLevelEvent, chan struct{}, error) {
	if dwr.Level == 0 {
		dwr.Level = 5
	}
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth%d", strings.ToLower(dwr.Symbol), dwr.Level)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {

		return nil, nil, err
	}

	//{"lastUpdateId":176319407,
	// "bids":[["0.00258680","21.03000000",[]],["0.00258650","2.29000000",[]],["0.00258640","367.74000000",[]],["0.00258630","310.58000000",[]],["0.00258520","103.37000000",[]]],
	// "asks":[["0.00258800","9.68000000",[]],["0.00258860","4.08000000",[]],["0.00258900","50.57000000",[]],["0.00258920","0.47000000",[]],["0.00258940","573.51000000",[]]]}

	done := make(chan struct{})
	dech := make(chan *DepthLevelEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawDepth := struct {
					//Type          string          `json:"-"`
					//Time          float64         `json:"-"`
					//Symbol        string          `json:"-"`
					UpdateID      int             `json:"lastUpdateId"`
					BidDepthDelta [][]interface{} `json:"bids"`
					AskDepthDelta [][]interface{} `json:"asks"`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

					return
				}
				//t, err := timeFromUnixTimestampFloat(rawDepth.Time)
				//if err != nil {
				//	level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
				//
				//	return
				//}
				de := &DepthLevelEvent{
					//WSEvent: WSEvent{
					//	Type:   rawDepth.Type,
					//	Time:   t,
					//	Symbol: rawDepth.Symbol,
					//},
					UpdateID: rawDepth.UpdateID,
					OrderBook: OrderBook{
						LastUpdateID: rawDepth.UpdateID,
					},
				}
				for _, b := range rawDepth.BidDepthDelta {
					p, err := floatFromString(b[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					q, err := floatFromString(b[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					de.Bids = append(de.Bids, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				for _, a := range rawDepth.AskDepthDelta {
					p, err := floatFromString(a[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					q, err := floatFromString(a[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					de.Asks = append(de.Asks, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				dech <- de
			}
		}
	}()

	go as.exitHandler(c, done)
	return dech, done, nil
}

func (as *apiService) DepthWebsocketStream(dwr DepthWebsocketRequestStream) (chan *DepthStreamEvent, chan struct{}, error) {
	if len(dwr.Streams) == 0 {
		return nil, nil, fmt.Errorf("streams is empty")
	}

	streams := ""
	for _, stream := range dwr.Streams {
		if streams != "" {
			streams += "/"
		}
		streams += strings.ToLower(stream)
	}
	//url := fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=ethbtc@depth5/bnbbtc@depth5")
	url := fmt.Sprintf("wss://stream.binance.com:9443/stream?streams=%s", streams)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {

		return nil, nil, err
	}

	//{"stream":"bnbbtc@depth5",
	// "data":{"lastUpdateId":177304217,
	// "bids":[["0.00247730","0.44000000",[]],["0.00247610","104.68000000",[]],["0.00247600","334.55000000",[]],["0.00247560","119.41000000",[]],["0.00247500","816.14000000",[]]],"asks":[["0.00247800","380.30000000",[]],["0.00247840","10.17000000",[]],["0.00247930","293.62000000",[]],["0.00247940","200.00000000",[]],["0.00247950","4.09000000",[]]]}}

	done := make(chan struct{})
	dech := make(chan *DepthStreamEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawDepth := struct {
					//Type          string          `json:"-"`
					//Time          float64         `json:"-"`
					//Symbol        string          `json:"-"`
					Stream string `json:"stream"`
					Data   struct {
						UpdateID      int             `json:"lastUpdateId"`
						BidDepthDelta [][]interface{} `json:"bids"`
						AskDepthDelta [][]interface{} `json:"asks"`
					} `json:data`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

					return
				}
				//t, err := timeFromUnixTimestampFloat(rawDepth.Time)
				//if err != nil {
				//	level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
				//
				//	return
				//}
				de := &DepthStreamEvent{
					//WSEvent: WSEvent{
					//	Type:   rawDepth.Type,
					//	Time:   t,
					//	Symbol: rawDepth.Symbol,
					//},
					Stream:   rawDepth.Stream,
					UpdateID: rawDepth.Data.UpdateID,
					OrderBook: OrderBook{
						LastUpdateID: rawDepth.Data.UpdateID,
					},
				}
				for _, b := range rawDepth.Data.BidDepthDelta {
					p, err := floatFromString(b[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					q, err := floatFromString(b[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					de.Bids = append(de.Bids, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				for _, a := range rawDepth.Data.AskDepthDelta {
					p, err := floatFromString(a[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					q, err := floatFromString(a[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))

						return
					}
					de.Asks = append(de.Asks, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				dech <- de
			}
		}
	}()

	go as.exitHandler(c, done)
	return dech, done, nil
}
