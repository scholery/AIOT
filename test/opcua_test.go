package test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
)

func ExecOpc() {
	endpoint := flag.String("endpoint", "opc.tcp://10.10.10.137:4840", "OPC UA Endpoint URL")
	policy := flag.String("policy", "", "Security policy: None, Basic128Rsa15, Basic256, Basic256Sha256. Default: auto")
	mode := flag.String("mode", "", "Security mode: None, Sign, SignAndEncrypt. Default: auto")
	certFile := flag.String("cert", "", "Path to cert.pem. Required for security mode/policy != None")
	keyFile := flag.String("key", "", "Path to private key.pem. Required for security mode/policy != None")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	ctx := context.Background()

	endpoints, err := opcua.GetEndpoints(*endpoint)
	if err != nil {
		log.Fatal(err)
	}
	ep := opcua.SelectEndpoint(endpoints, *policy, ua.MessageSecurityModeFromString(*mode))
	if ep == nil {
		log.Fatal("Failed to find suitable endpoint")
	}

	fmt.Println("*", ep.SecurityPolicyURI, ep.SecurityMode)

	opts := []opcua.Option{
		opcua.SecurityPolicy(*policy),
		opcua.SecurityModeString(*mode),
		opcua.CertificateFile(*certFile),
		opcua.PrivateKeyFile(*keyFile),
		opcua.AuthAnonymous(),
		opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
	}

	c := opcua.NewClient(ep.EndpointURL, opts...)
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	nodeId, _ := ua.ParseNodeID("ns=2;i=3")
	//v, err := c.Node(ua.NewNumericNodeID(2, 2)).Value()
	v, err := c.Node(nodeId).Value()

	switch {
	case err != nil:
		log.Fatal(err)
	case v == nil:
		log.Print("v == nil")
	default:
		log.Print(v.Int())
	}
	// Subscribe
	sub, err := c.Subscribe(&opcua.SubscriptionParameters{Interval: 500 * time.Millisecond}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()
	log.Printf("Created subscription with id %v", sub.SubscriptionID)

	//id, err := ua.ParseNodeID(*nodeID)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// arbitrary client handle for the monitoring item
	handle := uint32(42)
	id, _ := ua.ParseNodeID("ns=2;i=7")
	miCreateRequest := opcua.NewMonitoredItemCreateRequestWithDefaults(id, ua.AttributeIDValue, handle)
	res, err := sub.Monitor(ua.TimestampsToReturnBoth, miCreateRequest)
	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
		log.Fatal(err)
	}

	go sub.Run(ctx) // start Publish loop

	// read from subscription's notification channel until ctx is cancelled
	for {
		log.Printf("what's this publish result?")
		select {
		case <-ctx.Done():
			return
		case res := <-sub.Notifs:
			if res.Error != nil {
				log.Print(res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					//data := item.Value.Value.Vale
					log.Println("MonitoredItem with client handle ", item.ClientHandle)
				}

			default:
				log.Printf("what's this publish result? %T", res.Value)
			}
		}
	}
}

func TestOpc(t *testing.T) {
	ExecOpc()
}
