package common

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// ResponseType
	AckMessage = 0
	ErrMessage = 1
)

// BetMessage Message sent by the client to the server
//
// Contains:
// - SenderID: ID of the client that sent the message
// - Bet: List of bets sent by the client
type BetMessage struct {
	SenderID string
	bet      []ClientBet
}

type ResponseMessage struct {
	SenderID     string
	ResponseType int
}

func NewBetMessage(senderID string, bets []ClientBet) BetMessage {
	return BetMessage{
		SenderID: senderID,
		bet:      bets,
	}
}

// BetMessageFromBytes Converts a bet message to a byte array to be sent through the socket
func BetMessageToBytes(betMessage BetMessage) []byte {
	bytes := []byte(betMessage.SenderID + ":")
	for _, bet := range betMessage.bet {
		bytes = append(bytes, _ClientToBytes(bet)...)
	}
	bytes = append(bytes, byte('\r'), byte('\n'))
	return bytes
}

// _ClientToBytes Converts a single client bet to a byte array to be sent through the socket
func _ClientToBytes(bet ClientBet) []byte {
	return []byte(
		fmt.Sprintf(
			"%s,%s,%s,%s,%d\n",
			bet.Nombre,
			bet.Apellido,
			bet.Documento,
			bet.Nacimiento,
			bet.Numero,
		),
	)
}

// ResponseMessageFromBytes Converts a byte array to a response message received through the socket
//
// Protocol: <sender_id>:<response_type>\n
func ResponseMessageFromBytes(bytes []byte) ResponseMessage {
	// Remove the last byte which is a newline
	trimmed := bytes[:len(bytes)-1]
	split := strings.Split(string(trimmed), ":")
	senderID := split[0]
	responseType, err := strconv.Atoi(split[1])
	if err != nil {
		log.Errorf("action: parse_response_type | result: fail | error: %v", err)
	}
	return ResponseMessage{
		SenderID:     senderID,
		ResponseType: responseType,
	}
}
