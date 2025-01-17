package mqtt

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Options struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type StreamMessage struct {
	Topic string
	Value string
}

type Client struct {
	client paho.Client
	topics TopicMap
	stream chan StreamMessage
}

//  Name of our default topic. TODO: Support other topics
const defaultTopicName = "all"

//  We will subscribe to all MQTT topics. TODO: Support other topics
const defaultTopicMQTT = "#"

func NewClient(o Options) (*Client, error) {
	opts := paho.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", o.Host, o.Port))
	opts.SetClientID(fmt.Sprintf("grafana_%d", rand.Int()))

	if o.Username != "" {
		opts.SetUsername(o.Username)
	}

	if o.Password != "" {
		opts.SetPassword(o.Password)
	}

	opts.SetPingTimeout(60 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c paho.Client, err error) {
		log.DefaultLogger.Error(fmt.Sprintf("MQTT Connection Lost: %s", err.Error()))
	})
	opts.SetReconnectingHandler(func(c paho.Client, options *paho.ClientOptions) {
		log.DefaultLogger.Debug("MQTT Reconnecting")
	})

	log.DefaultLogger.Info("MQTT Connecting")

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker: %s", token.Error())
	}

	return &Client{
		client: client,
		stream: make(chan StreamMessage, 1000),
	}, nil
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnectionOpen()
}

func (c *Client) IsSubscribed(path string) bool {
	_, ok := c.topics.Load(path)
	return ok
}

func (c *Client) Messages(path string) ([]Message, bool) {
	topic, ok := c.topics.Load(path)
	if !ok {
		return nil, ok
	}
	return topic.messages, true
}

func (c *Client) Stream() chan StreamMessage {
	return c.stream
}

func (c *Client) HandleMessage(_ paho.Client, msg paho.Message) {
	log.DefaultLogger.Debug(fmt.Sprintf("Received MQTT Message for topic %s", msg.Topic()))
	//  Accept all topics as "all". TODO: Support other topics.
	//  Previously: topic, ok := c.topics.Load(msg.Topic())
	topic, ok := c.topics.Load(defaultTopicName)
	if !ok {
		log.DefaultLogger.Debug(fmt.Sprintf("Topic not found: %s", defaultTopicName))
		return
	}

	//  Compose message
	message := Message{
		Timestamp: time.Now(),
		Value:     string(msg.Payload()),
	}

	//  TODO: Fix this hack to reject messages without a valid CBOR Base64 Payload.
	//  CBOR Payloads must begin with a CBOR Map: 0xA1 or 0xA2 or 0xA3 or ...
	//  So the Base64 Encoding must begin with "o" or "p" or "q" or ...
	//  We stop at 0xB1 (Base64 "s") because we assume LoRaWAN Payloads will be under 50 bytes.
	//  Join Messages don't have a payload and will also be rejected.
	const frm_payload = "\"frm_payload\":\""
	if !strings.Contains(message.Value, frm_payload+"o") &&
		!strings.Contains(message.Value, frm_payload+"p") &&
		!strings.Contains(message.Value, frm_payload+"q") &&
		!strings.Contains(message.Value, frm_payload+"r") &&
		!strings.Contains(message.Value, frm_payload+"s") {
		log.DefaultLogger.Debug(fmt.Sprintf("Missing or invalid payload: %s", message.Value))
		return
	}

	// store message for query
	topic.messages = append(topic.messages, message)

	// limit the size of the retained messages
	if len(topic.messages) > 1000 {
		topic.messages = topic.messages[1:]
	}

	c.topics.Store(topic)

	//  Stream message to topic "all". TODO: Support other topics.
	//  Previously: streamMessage := StreamMessage{Topic: msg.Topic(), Value: string(msg.Payload())}
	streamMessage := StreamMessage{Topic: defaultTopicName, Value: string(msg.Payload())}

	log.DefaultLogger.Debug(fmt.Sprintf("Stream MQTT Message for topic %s", defaultTopicName))

	select {
	case c.stream <- streamMessage:
	default:
		// don't block if nothing is reading from the channel
	}
}

func (c *Client) Subscribe(t string) {
	if _, ok := c.topics.Load(t); ok {
		return
	}
	//  Subscribe to all topics: "#". TODO: Support other topics.
	//  Previously: log.DefaultLogger.Debug(fmt.Sprintf("Subscribing to MQTT topic: %s", t))
	log.DefaultLogger.Debug(fmt.Sprintf("Subscribing to MQTT topic: %s", defaultTopicMQTT))
	topic := Topic{
		path: t,
	}
	c.topics.Store(&topic)

	//  Subscribe to all topics: "#". TODO: Support other topics.
	//  Previously: c.client.Subscribe(t, 0, c.HandleMessage)
	c.client.Subscribe(defaultTopicMQTT, 0, c.HandleMessage)
}

func (c *Client) Unsubscribe(t string) {
	log.DefaultLogger.Debug(fmt.Sprintf("Unsubscribing from MQTT topic: %s", t))
	c.client.Unsubscribe(t)
	c.topics.Delete(t)
}

func (c *Client) Dispose() {
	log.DefaultLogger.Info("MQTT Disconnecting")
	c.client.Disconnect(250)
}
