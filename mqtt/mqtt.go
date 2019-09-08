package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Url       string   `json:"url"`
	IstTls    bool     `json:"tls"`
	CaPath    string   `json:"ca"`
	CertPath  string   `json:"cert"`
	KeyPath   string   `json:"key"`
	NumClient int      `json:"numClient"`
	ClientId  []string `json:"clientId"`
}

var (
	clientIDPrefix = "mqtt-carot"
)

type MqttClient interface {
	Connect() error
	Close()
	Subscribe(string, byte, MsgHandler) error
	SubscribeMultiple(map[string]byte, MsgHandler) error
	Unsubscribe(string) error
	AddRoute(string, MsgHandler)
	Publish(string, byte, bool, interface{}) error
}

type subscription struct {
	topic    string
	qos      byte
	callback MsgHandler
}

type Message paho.Message

type MsgHandler func(msg Message)

type mqttClient struct {
	client       paho.Client
	subscription map[string]subscription
}

var (
	defaultCallback = func(msg Message) {
		fmt.Println("Message from " + msg.Topic() + ": " + string(msg.Payload()))
	}

	defaultMsgHandler = func(client paho.Client, msg paho.Message) {
		defaultCallback(msg)
	}
)

func NewClient(uri string, id string) (MqttClient, error) {

	client := &mqttClient{}
	client.subscription = make(map[string]subscription)

	opts := paho.NewClientOptions().AddBroker(uri)
	//	opts.SetClientID(id)
	opts.AutoReconnect = true
	//	opts.KeepAlive = 20
	opts.CleanSession = true
	opts.SetDefaultPublishHandler(defaultMsgHandler)
	opts.SetOnConnectHandler(client.onConnect)

	c := paho.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	client.client = c
	return client, nil
}

func NewClientTLS(uri string, id string, caPath, certPath, keyPath string) (MqttClient, error) {
	log.Println("NewClientTLS", caPath)
	log.Println("NewClientTLS", certPath)
	log.Println("NewClientTLS", keyPath)
	log.Println("NewClientID", id)
	client := &mqttClient{}
	client.subscription = make(map[string]subscription)
	clientId := id
	if len(clientId) <= 0 {
		clientId = fmt.Sprintf("%s-%d", clientIDPrefix, time.Now().UnixNano())
	}

	opts := paho.NewClientOptions().AddBroker(uri)
	opts.CleanSession = true
	opts.AutoReconnect = true
	opts.ClientID = clientId
	//	opts.KeepAlive = 20
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(caPath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Println("NewClientTLS", err.Error())
		return client, err
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Println("NewClientTLS", err.Error())
	}

	cfg := &tls.Config{
		RootCAs: certpool,
		//InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
	}

	opts.TLSConfig = cfg

	opts.CleanSession = true
	opts.SetDefaultPublishHandler(defaultMsgHandler)
	opts.SetOnConnectHandler(client.onConnect)

	c := paho.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	client.client = c
	return client, nil
}

func (m *mqttClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *mqttClient) onConnect(c paho.Client) {
	log.Println("======= onConnect")
	if len(m.subscription) <= 0 {
		return
	}
	for _, v := range m.subscription {
		m.Subscribe(v.topic, v.qos, v.callback)
	}
}

func (m *mqttClient) Close() {
	if m.client.IsConnected() {
		m.client.Disconnect(1000)
	}
}

func (m *mqttClient) Subscribe(topic string, qos byte, callback MsgHandler) error {

	token := m.client.Subscribe(topic, qos, func(client paho.Client, msg paho.Message) {
		callback(msg)
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	sub := subscription{topic, qos, callback}
	m.subscription[topic] = sub

	return nil
}

func (m *mqttClient) SubscribeMultiple(filters map[string]byte, callback MsgHandler) error {

	token := m.client.SubscribeMultiple(filters, func(client paho.Client, msg paho.Message) {
		callback(msg)
	})

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for k, v := range filters {
		sub := subscription{k, v, callback}
		m.subscription[k] = sub
	}

	return nil
}

func (m *mqttClient) Unsubscribe(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	delete(m.subscription, topic)

	return nil
}

func (m *mqttClient) AddRoute(topic string, callback MsgHandler) {
	m.client.AddRoute(topic, func(client paho.Client, message paho.Message) {
		callback(message)
	})
	sub := m.subscription[topic]
	sub.callback = callback
	m.subscription[topic] = sub
}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := m.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
