package opcClient

import (
	"context"
	"fmt"

	"koudai-box/iot/gateway/model"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/sirupsen/logrus"
)

//httpserver
var cache_servers map[string]*opcua.Client = make(map[string]*opcua.Client)

func ConnectOPC(gateway *model.GatewayConfig) error {
	url := fmt.Sprintf("opc.tcp://%s:%d", gateway.Ip, gateway.Port)
	// endpoint := flag.String("endpoint", url, "OPC UA Endpoint URL")
	endpoint := url
	// policy := flag.String("policy", "", "Security policy: None, Basic128Rsa15, Basic256, Basic256Sha256. Default: auto")
	policy := ""
	// mode := flag.String("mode", "", "Security mode: None, Sign, SignAndEncrypt. Default: auto")
	mode := ""
	//flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	//flag.Parse()

	ctx := context.Background()

	//endpoints, err := opcua.GetEndpoints(*endpoint)
	endpoints, err := opcua.GetEndpoints(endpoint)
	if err != nil {
		logrus.Errorf("ConnectOPC gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
		return fmt.Errorf("ConnectOPC gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
	}
	ep := opcua.SelectEndpoint(endpoints, policy, ua.MessageSecurityModeFromString(mode))
	if ep == nil {
		logrus.Error("ConnectOPC gateway[%s] start error.url[%s].Failed to find suitable endpoint", gateway.Key, url)
		return fmt.Errorf("ConnectOPC gateway[%s] start error.url[%s].Failed to find suitable endpoint", gateway.Key, url)
	}

	//fmt.Println("*", ep.SecurityPolicyURI, ep.SecurityMode)
	//certFile := flag.String("cert", "", "Path to cert.pem. Required for security mode/policy != None")
	//keyFile := flag.String("key", "", "Path to private key.pem. Required for security mode/policy != None")
	certFile := ""
	keyFile := ""
	opts := []opcua.Option{
		opcua.SecurityPolicy(policy),
		opcua.SecurityModeString(mode),
		opcua.CertificateFile(certFile),
		opcua.PrivateKeyFile(keyFile),
		opcua.AuthAnonymous(),
		opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
	}

	client := opcua.NewClient(ep.EndpointURL, opts...)
	if err := client.Connect(ctx); err != nil {
		logrus.Errorf("ConnectOPC gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
		return fmt.Errorf("ConnectOPC gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
	}
	cache_servers[gateway.Key] = client
	logrus.Infof("++++++++++++++开启OPC客户端[%s][%s]++++++++++++++", gateway.Key, url)
	return nil
}

func CloseOPC(gateway *model.GatewayConfig) error {
	client, ok := cache_servers[gateway.Key]
	if !ok {
		logrus.Errorf("opc[%s]'s client is not exist.", gateway.Key)
		return fmt.Errorf("opc[%s]'s client is not exist", gateway.Key)
	}
	client.Close()
	logrus.Infof("++++++++++++++关闭OPC客户端[%s]++++++++++++++", gateway.Key)
	return nil
}

func QueryValue(gateway *model.GatewayConfig, device *model.Device, item model.ItemConfig) (interface{}, error) {
	client, ok := cache_servers[gateway.Key]
	if !ok {
		err := ConnectOPC(gateway)
		client, ok = cache_servers[gateway.Key]
		if err != nil || !ok {
			logrus.Errorf("opc[%s]'s client is not exist.", gateway.Key)
			return nil, fmt.Errorf("opc[%s]'s client is not exist", gateway.Key)
		}
	}
	nodeIdStr := fmt.Sprintf("ns=%s;i=%s", device.SourceId, item.NodeId)
	nodeId, err := ua.ParseNodeID(nodeIdStr)
	if err != nil {
		logrus.Errorf("ParseNodeID opc[%s]'s node[%s] err,%+v", gateway.Key, nodeIdStr, err)
		return nil, fmt.Errorf("ParseNodeID opc[%s]'s node[%s] err,%+v", gateway.Key, nodeIdStr, err)
	}
	value, err := client.Node(nodeId).Value()
	if err != nil || nil == value {
		logrus.Errorf("opc[%s]'s node[%s] err,%+v", gateway.Key, nodeIdStr, err)
		return nil, fmt.Errorf("opc[%s]'s node[%s] err,%+v", gateway.Key, nodeIdStr, err)
	}
	logrus.Debugf("opc[%s]'s node[%s] value[%+v]", gateway.Key, nodeIdStr, value.Value())
	return value.Value(), nil
}
