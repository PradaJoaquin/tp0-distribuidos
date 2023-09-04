package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// ClientBet used by the client
type ClientBet struct {
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     int
}

// Client Entity that encapsulates how
type Client struct {
	config   ClientConfig
	conn     net.Conn
	shutdown chan os.Signal
	bet      ClientBet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet ClientBet) *Client {
	// shutdown is a channel used to receive the SIGTERM signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM)

	client := &Client{
		config:   config,
		shutdown: shutdown,
		bet:      bet,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		case <-c.shutdown:
			shutdown(c)
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()
		msg, err := sendClientBet(c.config.ID, c.bet, c.conn)
		c.conn.Close()
		if err != nil || msg.ResponseType == ErrMessage {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		if msg.ResponseType == AckMessage {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", c.bet.Documento, c.bet.Numero)
			break loop
		}
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

// shutdown Closes the connection and exits the program
func shutdown(c *Client) {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
	os.Exit(0)
}
