package sbuslistener

import (
	"awesomeProject/internal/utils"
	"context"
	servicebus "github.com/Azure/azure-service-bus-go"
	"os"
)

func MakeListner(coonstring string, c chan *servicebus.Message) {
	ctx := context.Background()
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(coonstring))
	utils.Check(err)
	 topic,err := ns.NewTopic(os.Getenv("TOPIC"))
	 utils.Check(err)
	 sub,err:=topic.NewSubscription(os.Getenv("SUB_NAME"),servicebus.SubscriptionWithPrefetchCount(1000))
	 utils.Check(err)
	err = sub.Receive(ctx, servicebus.HandlerFunc(func(ctx context.Context, message *servicebus.Message) error {
		c <- message
		return nil
	}))
	utils.Check(err)
}