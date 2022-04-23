package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/timkippdev/cryptochain/pkg/blockchain"
)

const (
	PubSubChannelBlockchain = "BLOCKCHAIN"
	PubSubChannelTest       = "TEST"
)

var (
	ctx = context.Background()
)

type PubSub struct {
	blockchain       *blockchain.Blockchain
	publisherClient  *redis.Client
	subscriberClient *redis.Client
	subscribers      []*redis.PubSub
}

type PubSubMessage struct {
	Identifier string `json:"identifier"`
	Data       string `json:"data"`
}

func NewPubSub(address string, bc *blockchain.Blockchain) *PubSub {
	opts := &redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	}

	ps := &PubSub{
		blockchain:       bc,
		publisherClient:  redis.NewClient(opts),
		subscriberClient: redis.NewClient(opts),
	}

	ps.AddSubscriber(PubSubChannelBlockchain, func(msg *redis.Message) {
		// TODO: update to decode into PubSubMessage and skip if identifier is same as publisher
		fmt.Printf("(%s) received message :: %+v\n", msg.Channel, msg.Payload)
		var incomingChain []*blockchain.Block
		err := json.Unmarshal([]byte(msg.Payload), &incomingChain)
		if err != nil {
			fmt.Printf("\tprocessing error: %s", err)
			return
		}
		ps.blockchain.ReplaceChain(incomingChain)
	})

	return ps
}

func (ps *PubSub) AddSubscriber(channel string, onMessageReceivedFunc func(*redis.Message)) *redis.PubSub {
	sub := ps.subscriberClient.Subscribe(ctx, channel)
	ps.subscribers = append(ps.subscribers, sub)

	go func() {
		for msg := range sub.Channel() {
			onMessageReceivedFunc(msg)
		}
	}()

	return sub
}

func (ps *PubSub) BroadcastChain() error {
	chainAsJSONBytes, err := json.Marshal(ps.blockchain.GetChain())
	if err != nil {
		return err
	}

	ps.PublishMessage(PubSubChannelBlockchain, PubSubMessage{
		Data: string(chainAsJSONBytes),
	})

	return nil
}

func (ps *PubSub) Close() {
	for _, s := range ps.subscribers {
		s.Close()
	}
}

func (ps *PubSub) PublishMessage(channel string, message PubSubMessage) {
	fmt.Printf("(%s) published message :: %+v\n", channel, message)
	// TODO: publish entire PubSubMessage (need to implement marshalers)
	res := ps.publisherClient.Publish(ctx, channel, message.Data)
	if res.Err() != nil {
		fmt.Println("(ERROR)", res.Err())
	}
}
