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
	BatchSize     int
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
	betsPath string
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, betsPath string) *Client {
	// shutdown is a channel used to receive the SIGTERM signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM)

	client := &Client{
		config:   config,
		shutdown: shutdown,
		betsPath: betsPath,
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
	betLoader, err := NewBetLoader(c.betsPath, c.config.BatchSize)
	if err != nil {
		log.Errorf("action: open_bets_file | result: fail | error: %v", err)
		return
	}
	for betLoader.HasNext() {
		select {
		case <-c.shutdown:
			shutdown(c)
		default:
		}
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		bets, err := Next(betLoader)
		if err != nil {
			log.Errorf("action: read_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
		}
		response, err := sendClientBetsBatch(c.config.ID, bets, c.conn)

		c.conn.Close()

		if err != nil || response.ResponseType == ErrMessage {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		if response.ResponseType == AckMessage {
			log.Infof("action: batch_de_apuestas_enviada | result: success")
		}
	}
	log.Infof("action: todas_las_apuestas_enviadas | result: success | client_id: %v", c.config.ID)
}

// shutdown Closes the connection and exits the program
func shutdown(c *Client) {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
	os.Exit(0)
}
