package service

import (
	"context"
	"fmt"
	psub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/samer955/gomdnsdisco/metrics"
	"github.com/samer955/gomdnsdisco/node"
	"github.com/samer955/gomdnsdisco/pubsub"
	"log"
	"time"
)

const (
	roomSystem = "system"
	roomCpu    = "cpu"
)

var (
	ctx         context.Context
	systemTopic *psub.Topic
	systemSub   *psub.Subscription
)

type Sender struct {
	node   node.Node
	pubsub pubsub.PubSubService
}

func NewSender() *Sender {

	ctx := context.Background()
	node := node.InitializeNode(ctx)
	pubsub := pubsub.NewPubSubService(ctx, node.Host)

	return &Sender{*node, *pubsub}

}

func (s *Sender) subscribeTopics() {
	sysTopic, err := s.pubsub.JoinTopic(roomSystem)
	if err != nil {
		log.Fatal(err)
	}
	systemTopic = sysTopic
	systemSub, _ = s.pubsub.Subscribe(systemTopic)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Sender) sendSystemInfo(system *metrics.System) {

	for {
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}
		system.GetSystemInformation()
		err := s.pubsub.Publish(system, ctx, systemTopic)
		if err != nil {
			log.Println("Error publishing data", err)
		}
		log.Println("Publishing Data: ", system)

		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) receiveSystemData() {
	for {
		m, _ := systemSub.Next(context.TODO())
		if m.ReceivedFrom.ShortString() == s.node.Host.ID().ShortString() {
			continue
		}
		fmt.Printf("Received message from <%s: %s> \n", m.ReceivedFrom.ShortString(), m.Data)
	}
}

func (s *Sender) Start() {

	s.subscribeTopics()
	go s.sendSystemInfo(new(metrics.System))
	go s.receiveSystemData()
}
