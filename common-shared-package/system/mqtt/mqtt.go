package MQTT

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Masher828/MessengerBackend/common-shared-package/conf"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MqttConnect() (mqtt.Client, error) {

	broker := conf.MessengerConfig.Mqtt.Broker
	username := conf.MessengerConfig.Mqtt.Username
	password := conf.MessengerConfig.Mqtt.Password

	var err error = nil
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.TLSConfig, err = loadTLSConfig()
	if err != nil {
		return nil, err
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	ack := token.WaitTimeout(5 * time.Second)

	if token.Error() != nil || !ack {
		err = errors.New(fmt.Sprintf("connect%s mqtt server error: %s", broker, token.Error()))
		return nil, err
	}

	return client, nil
}

func loadTLSConfig() (*tls.Config, error) {
	// load tls config

	serverCaFilePath := conf.MessengerConfig.Mqtt.ServerCaFile
	clientCerti := conf.MessengerConfig.Mqtt.ClientCertificateFile
	clientKey := conf.MessengerConfig.Mqtt.ClientKeyFile

	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	if serverCaFilePath != "" {
		certpool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(serverCaFilePath)
		if err != nil {
			return nil, err
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	if clientCerti != "" && clientKey != "" {
		clientKeyPair, err := tls.LoadX509KeyPair(clientCerti, clientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.ClientAuth = tls.RequestClientCert
		tlsConfig.Certificates = []tls.Certificate{clientKeyPair}
	}
	return &tlsConfig, nil
}
