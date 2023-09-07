package common

import (
	"errors"
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
	config    ClientConfig
	conn      net.Conn
	shutdown  chan os.Signal
	betLoader *betLoader
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, betsPath string) *Client {
	// shutdown is a channel used to receive the SIGTERM signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM)

	betLoader, err := NewBetLoader(betsPath, config.BatchSize)
	if err != nil {
		log.Errorf("action: open_bets_file | result: fail | error: %v", err)
		return nil
	}

	client := &Client{
		config:    config,
		shutdown:  shutdown,
		betLoader: betLoader,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)

	// Set a deadline for the connection, using the loop lapse as a timeout.
	conn.SetDeadline(time.Now().Add(c.config.LoopLapse))
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
	c.createClientSocket()
	err := begginBatch(c)
	if err != nil {
		log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
		shutdown(c)
		return
	}
	log.Infof("action: todas_las_apuestas_enviadas | result: success | client_id: %v", c.config.ID)

	err = begginDone(c)
	if err != nil {
		log.Errorf("action: send_done | result: fail | client_id: %v | error: %v", c.config.ID, err)
		shutdown(c)
		return
	}
	log.Infof("action: done_enviado | result: success | client_id: %v", c.config.ID)

	winners, err := begginWinners(c)
	if err != nil {
		log.Errorf("action: send_request_winners | result: fail | client_id: %v | error: %v", c.config.ID, err)
		shutdown(c)
		return
	}
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
	shutdown(c)
}

func begginBatch(c *Client) error {
	for c.betLoader.HasNext() {
		select {
		case <-c.shutdown:
			shutdown(c)
		default:
		}

		bets, err := Next(c.betLoader)
		if err != nil {
			log.Errorf("action: read_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
		}
		response, err := sendClientBetsBatch(c.config.ID, bets, c.conn)

		if err != nil || response.MessageType == ErrMessage {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			err = errors.New("Error al enviar el mensaje de apuestas")
			return err
		}
		if response.MessageType == AckMessage {
			log.Debugf("action: batch_de_apuestas_enviada | result: success")
		}
	}
	return nil
}

type error interface {
	Error() string
}

func begginDone(c *Client) error {
	response, err := sendDoneSendingBets(c.config.ID, c.conn)
	if response.MessageType == ErrMessage {
		err = errors.New("Error al enviar el mensaje de done")
	}
	return err
}

func begginWinners(c *Client) ([]ClientBet, error) {
	// default response to start the loop
	response := ResponseWinnersMessage{Message: Message{MessageType: WaitMessage}}
	for response.MessageType != WinnersMessage {
		select {
		case <-c.shutdown:
			shutdown(c)
		default:
		}
		new_response, err := sendRequestWinners(c.config.ID, c.conn)
		response = new_response
		if err != nil || response.MessageType == ErrMessage {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			err = errors.New("Error al enviar el mensaje de winners")
			return nil, err
		}
		if response.MessageType == WinnersMessage {
			break
		}
		if response.MessageType == WaitMessage {
			log.Debugf("action: wait_for_winners | result: success")
			// Wait for the loop period before retrying
			time.Sleep(c.config.LoopPeriod)
		}
	}
	return response.Winners, nil
}

// shutdown Closes the connection and exits the program, closing all the files descriptors.
func shutdown(c *Client) {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
	log.Infof("action: close_bets_file | result: in_progress | client_id: %v", c.config.ID)
	c.betLoader.Close()
	log.Infof("action: close_bets_file | result: success | client_id: %v", c.config.ID)
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
	os.Exit(0)
}
