package kafka

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(KafkaPublisher{})
	httpcaddyfile.RegisterHandlerDirective("kafka_publisher", parseCaddyfile)
}

// Middleware implements a handler that receiver messages over HTTP, performs protocol conversion and
// writes to designated Kafka topics
type KafkaPublisher struct {
	BootstrapServers string `json:"bootstrapServers"`
	producer         sarama.SyncProducer
	logger           *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (KafkaPublisher) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.kafka_publisher",
		New: func() caddy.Module { return new(KafkaPublisher) },
	}
}

// Provision implements caddy.Provisioner.
func (kp *KafkaPublisher) Provision(ctx caddy.Context) error {

	brokers := strings.Split(kp.BootstrapServers, ",")
	// kp.logger.Sugar().Infof("Bootstrap servers %+v", brokers)

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	kp.producer = producer
	kp.logger = ctx.Logger(kp)

	return nil
}

// Validate implements caddy.Validator.
func (kp *KafkaPublisher) Validate() error {

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (kp KafkaPublisher) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	// Receive message from service over HTTP

	req := r.Body

	// Read params such as topic name, destination address, auth info and content to post

	topicName := r.Header.Get("X-KAFKA-TOPIC")
	// auth := r.Header.Get("X-AUTH-INFO")
	message, err := ioutil.ReadAll(req)

	if err != nil {
		return err
	}

	// Post to topic

	msg := &sarama.ProducerMessage{
		Topic:     topicName,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	_, _, err = kp.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	// Respond with 200 or error code
	w.Write([]byte("Done"))
	return nil

	// return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (kp *KafkaPublisher) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&kp.BootstrapServers) {
			return d.ArgErr()
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var kp KafkaPublisher
	err := kp.UnmarshalCaddyfile(h.Dispenser)
	return kp, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*KafkaPublisher)(nil)
	_ caddy.Validator             = (*KafkaPublisher)(nil)
	_ caddyhttp.MiddlewareHandler = (*KafkaPublisher)(nil)
	_ caddyfile.Unmarshaler       = (*KafkaPublisher)(nil)
)
