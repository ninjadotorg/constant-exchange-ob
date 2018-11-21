package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/ninjadotorg/constant-exchange-ob/config"
	"github.com/ninjadotorg/constant-exchange-ob/messages"
	"github.com/ninjadotorg/constant-exchange-ob/orderbook"
	"github.com/ninjadotorg/constant-exchange-ob/services"
	"github.com/ninjadotorg/constant-exchange-ob/utils"
	"os"
	"sync"
	"time"
	"fmt"
)

const (
	TOPIC_ORDER = "order_stresstest"
	TOPIC_ORDERBOOK = "orderbook_stresstest"
)


var ai = 0
var mutex sync.Mutex

func randomOrder(marketId int, side string, price float64, size float64) *orderbook.Order {
	mutex.Lock()
	p := price
	s := size
	sd := side

	if p == 0 {
		p = utils.RandomFloat64(1,5)
	}

	if s == 0 {
		s = utils.RandomFloat64(1, 5)
	}

	if sd == "" {
		randomSide := utils.RandomInt(0,2)
		fmt.Println("randomSide", randomSide)
		if randomSide == 1 {
			sd = "sell"
		} else {
			sd = "buy"
		}
	}

	ai++
	mutex.Unlock()

	return &orderbook.Order{
		ID: ai,
		Price: p,
		Size: s,
		Side: sd,
		MarketID: marketId,
	}
}

func main() {
	conf := config.GetConfig()

	ctx := context.Background()

	ps := services.InitPubSub(conf.GCProjectID)

	orderTopic := ps.GetOrCreateTopic(TOPIC_ORDER)
	orderbookTopic := ps.GetOrCreateTopic(TOPIC_ORDERBOOK)

	var orderbookSubscribe *pubsub.Subscription
	changeSubscribeName := "order_stresstest"

	symbols := []int{0,1,2,3,4}

	if orderbookTopic != nil {
		orderbookSubscribe = ps.GetOrCreateSubscription(changeSubscribeName, orderbookTopic)
		go orderbookSubscribe.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()

			var orderbookMessage messages.OrderBookMessage
			if err := json.Unmarshal(msg.Data, &orderbookMessage); err != nil {
				fmt.Printf("could not decode message data: %#v", msg)
				return
			}
			if orderbookMessage.Type == "match" {
				fmt.Println("ob match", string(msg.Data))
			}
		})
	}

	for _, symbol := range symbols {
		totalOrder := utils.RandomInt(10000, 50000)
		for i := 0; i < totalOrder; i++ {
			o := randomOrder(symbol, "", 0, 0)
			msg := map[string]interface{}{
				"type":"add",
				"order": o,
			}
			msgJson, _ := json.Marshal(msg)
			uuid, err := orderTopic.Publish(ctx, &pubsub.Message{Data: msgJson}).Get(ctx)
			if err != nil {
				fmt.Println("Publish order error", err.Error())
			} else {
				fmt.Println("Add order", string(msgJson), uuid)
			}

			//time.Sleep(10 * time.Millisecond)
		}
	}

	time.Sleep(30 * time.Second)

	err := orderbookSubscribe.Delete(ctx)
	if err != nil {
		fmt.Println("Delete subscription error: ", err.Error())
	}

	os.Exit(0)
}