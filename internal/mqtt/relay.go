package mqtt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/Ege-Guler/gochat/internal/config"
)

func StartSession() error {

	log.Println("Starting session...")

	s := &Session{}

	if err := s.GenerateSessionID(); err != nil {
		return fmt.Errorf("failed to generate session ID: %w", err)
	}
	if err := s.GenerateOwnKeys(); err != nil {
		return fmt.Errorf("failed to generate private and public keys: %w", err)
	}

	startRelay(s)

	return nil
}

func startRelay(s *Session) error {

	var conf config.Config
	if err := conf.LoadConf(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tls://%s:%d", conf.MQTT.Broker, conf.MQTT.Port))
	opts.SetClientID(conf.MQTT.ClientID)
	opts.SetUsername(conf.MQTT.Username)
	opts.SetPassword(conf.MQTT.Password)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	defer client.Disconnect(250)

	if err := subscribe(client, &conf, s); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	if err := publish(client, &conf, s, InitExchangeFunc); err != nil {
		return fmt.Errorf("failed to publish messages: %w", err)
	}

	return nil
}

func subscribe(client mqtt.Client, cfg *config.Config, s *Session) error {
	topic := cfg.MQTT.Topic

	token := client.Subscribe(topic, byte(cfg.MQTT.QoS), callback(s))
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", cfg.MQTT.Topic, token.Error())
	}

	fmt.Printf("Subscribed to topic: %s\n", topic)
	return nil
}

func callback(s *Session) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {

		var g generic

		json.Unmarshal(msg.Payload(), &g)

		switch g.Type {
		case IM:
			imHandler(msg, s)
		case KM:
			kmHandler(msg, s)
		default:
			log.Printf("Unknown message type: %s", g.Type)
		}

	}
}

func imHandler(msg mqtt.Message, s *Session) {

	var im InitExchangeMessage
	if err := im.Decode(msg.Payload()); err != nil {
		log.Printf("Failed to decode exchange message: %v", err)
		return
	}
	//correct way
	log.Println(base64.StdEncoding.EncodeToString(im.PublicKey), string(im.Type), "zaa")

	// ignore sent messages for the same session
	if im.SessionID != s.SessionID {

		log.Printf("Received message for different session %s, ignoring \n", im.SessionID)

		s.PeerPubKey = im.PublicKey

		if err := s.GenerateKEK(); err != nil {
			log.Printf("Failed to generate KEK: %v", err)
			return
		}
		if err := s.SetAES(); err != nil {
			log.Printf("Failed to set AES: %v", err)
			return
		}

	}
}

func kmHandler(msg mqtt.Message, s *Session) {

}

func publish(client mqtt.Client, cfg *config.Config, s *Session, messageFunc func(*Session) ([]byte, error)) error {

	log.Println("Publishing messages to topic:", cfg.MQTT.Topic)
	payload, err := messageFunc(s)
	if err != nil {
		return fmt.Errorf("failed to encode exchange message: %w", err)
	}

	token := client.Publish(cfg.MQTT.Topic, byte(cfg.MQTT.QoS), cfg.MQTT.Retain, payload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}

	fmt.Printf("Published message: %s\n", payload)

	return nil

}

// func publish(client mqtt.Client, cfg *config.Config, s *Session) error {

// 	log.Println("Publishing messages to topic:", cfg.MQTT.Topic)

// 	im := s.ToInitExchange()
// 	payload, err := im.Encode()
// 	if err != nil {
// 		return fmt.Errorf("failed to encode exchange message: %w", err)
// 	}

// 	token := client.Publish(cfg.MQTT.Topic, byte(cfg.MQTT.QoS), cfg.MQTT.Retain, payload)
// 	token.Wait()

// 	if token.Error() != nil {
// 		return fmt.Errorf("failed to publish message: %w", token.Error())
// 	}

// 	fmt.Printf("Published message: %s\n", payload)
// 	time.Sleep(1 * time.Second)

// 	return nil
// }
