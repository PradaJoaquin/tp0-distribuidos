package common

import (
	"net"
)

const BUFFER_SIZE = 8192

// SendClientBets Sends a batch of client bets to the server and returns the response
func sendClientBetsBatch(senderID string, bets []ClientBet, conn net.Conn) (ResponseMessage, error) {
	betMessage := NewBetMessage(senderID, bets)
	bytes := BetMessageToBytes(betMessage)

	bytesSent := 0
	for bytesSent < len(bytes) {
		bytesSentAux, err := conn.Write(bytes[bytesSent:])
		if err != nil {
			return ResponseMessage{}, err
		}
		bytesSent += bytesSentAux
	}
	return receiveResponse(conn)
}

// ReceiveResponse Receives a response message from the server
func receiveResponse(conn net.Conn) (ResponseMessage, error) {
	bytes := make([]byte, BUFFER_SIZE)
	bytesReceived := 0
	// Read until a newline is received to prevent short reads
	for bytesReceived == 0 || bytes[bytesReceived-1] != byte('\n') {
		bytesReceivedAux, err := conn.Read(bytes[bytesReceived:])
		if err != nil {
			return ResponseMessage{}, err
		}
		bytesReceived += bytesReceivedAux
	}
	return ResponseMessageFromBytes(bytes[:bytesReceived]), nil
}
