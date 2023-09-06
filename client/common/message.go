package common

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// MessageType
	AckMessage      = 0
	ErrMessage      = 1
	WaitMessage     = 2
	WinnersMessage  = 3
	BetMessageType  = 4
	DoneSendingBets = 5
	RequestWinners  = 6
)

type Message struct {
	SenderID    string
	MessageType int
}

// BetMessage Message sent by the client to the server
//
// Contains:
// - SenderID: ID of the client that sent the message
// - Bet: List of bets sent by the client
type BetMessage struct {
	bet []ClientBet
	Message
}

type ResponseWinnersMessage struct {
	Message
	Winners []ClientBet
}

func NewBetMessage(senderID string, bets []ClientBet) BetMessage {
	return BetMessage{
		bet:     bets,
		Message: Message{SenderID: senderID, MessageType: BetMessageType},
	}
}

func NewDoneSendingBetsMessage(senderID string) Message {
	return Message{SenderID: senderID, MessageType: DoneSendingBets}
}

func NewRequestWinnersMessage(senderID string) Message {
	return Message{SenderID: senderID, MessageType: RequestWinners}
}

// BetMessageFromBytes Converts a bet message to a byte array to be sent through the socket
//
// Protocol: <sender_id>:<message_type>:<client_bet_1>\n<client_bet_2>\n...\n<client_bet_n>\n\r\n
func BetMessageToBytes(betMessage BetMessage) []byte {
	bytes := []byte(betMessage.SenderID + ":" + strconv.Itoa(BetMessageType) + ":")
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

// DoneSendingBetsMessageToBytes Converts a DoneSendingBets message to a byte array to be sent through the socket.
//
// It is used to notify the server that the client has finished sending bets.
//
// Protocol: <sender_id>:<message_type>\r\n
func DoneSendingBetsMessageToBytes(message Message) []byte {
	return []byte(message.SenderID + ":" + strconv.Itoa(DoneSendingBets) + "\r\n")
}

// RequestWinnersMessageToBytes Converts a RequestWinners message to a byte array to be sent through the socket.
//
// It is used to notify the server that the client wants to receive the winners, it should be sent after the done sending bets message.
//
// Protocol: <sender_id>:<message_type>\r\n
func RequestWinnersMessageToBytes(message Message) []byte {
	return []byte(message.SenderID + ":" + strconv.Itoa(RequestWinners) + "\r\n")
}

// MessageFromBytes Converts a byte array to a response message received through the socket
//
// Protocol: <sender_id>:<message_type>\n
func MessageFromBytes(bytes []byte) Message {
	// Remove the last 2 bytes which is the end of message
	trimmed := bytes[:len(bytes)-2]
	split := strings.Split(string(trimmed), ":")
	senderID := split[0]
	MessageType, err := strconv.Atoi(split[1])
	if err != nil {
		log.Errorf("action: parse_response_type | result: fail | error: %v", err)
	}
	return Message{
		SenderID:    senderID,
		MessageType: MessageType,
	}
}

// ResponseWinnersMessageFromBytes Converts a response a byte array to a winners message to be sent through the socket
//
// Protocol: <sender_id>:<message_type>:<client_bet_1>\n<client_bet_2>\n...\n<client_bet_n>\n\r\n
func ResponseWinnersMessageFromBytes(bytes []byte) (ResponseWinnersMessage, error) {
	// Remove the last 2 bytes which is the end of message
	trimmed := bytes[:len(bytes)-2]
	splited := strings.Split(string(trimmed), ":")
	senderID := splited[0]
	MessageType, err := strconv.Atoi(splited[1])
	if err != nil {
		log.Errorf("action: parse_response_type | result: fail | error: %v", err)
		return ResponseWinnersMessage{}, err
	}
	if MessageType == WaitMessage {
		return ResponseWinnersMessage{
			Message: Message{
				SenderID:    senderID,
				MessageType: MessageType,
			},
		}, nil
	}
	winners := splited[2]
	winnersList := strings.Split(winners, "\n")

	// Create a list of ClientBet from the winners list
	var winnersBet []ClientBet
	for _, winner := range winnersList {
		if winner != "" {
			winnerBet := strings.Split(winner, ",")
			numero, err := strconv.Atoi(winnerBet[5])
			if err != nil {
				log.Errorf("action: parse_response_type | result: fail | error: %v", err)
				return ResponseWinnersMessage{}, err
			}
			winnersBet = append(winnersBet, ClientBet{
				Nombre:     winnerBet[1],
				Apellido:   winnerBet[2],
				Documento:  winnerBet[3],
				Nacimiento: winnerBet[4],
				Numero:     numero,
			})
		}
	}
	return ResponseWinnersMessage{
		Message: Message{
			SenderID:    senderID,
			MessageType: MessageType,
		},
		Winners: winnersBet,
	}, nil
}
