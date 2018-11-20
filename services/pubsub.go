package services

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
)

type PubSub struct{
	client *pubsub.Client
}

func InitPubSub(projectId string) *PubSub {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		fmt.Println("pubsub Init err", err.Error())
		return nil
	}

	pubsub := &PubSub{
		client: pubsubClient,
	}

	return pubsub
}

func (ps *PubSub) GetOrCreateTopic(topicName string) *pubsub.Topic {
	ctx := context.Background()
	var topic *pubsub.Topic
	var err error
	topic = ps.client.Topic(topicName)
	ok, err := topic.Exists(ctx)

	if err != nil {
		fmt.Println("Topic.Exists error", err.Error())
		return nil
	}

	if ok {
		return topic
	}

	// create topic
	topic, err = ps.client.CreateTopic(context.Background(), topicName)
	if err != nil {
		fmt.Println("pubsub GetOrCreateTopic err", err.Error())
		return nil
	}

	return topic
}

func (ps *PubSub) GetOrCreateSubscription(subscriptionName string, topic *pubsub.Topic) *pubsub.Subscription {
	ctx := context.Background()
	var subscription *pubsub.Subscription
	var err error
	subscription = ps.client.Subscription(subscriptionName)

	ok, err := subscription.Exists(ctx)
	if err != nil {
		fmt.Println("Subscription.Exists error", err.Error())
		return nil
	}

	if ok {
		return subscription
	}

	// create topic
	subscription, err = ps.client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		fmt.Println("pubsub GetOrCreateSubscription err", err.Error())
		return nil
	}

	return subscription
}