package common

import (
	"net"
)

const BUFFER_SIZE = 8192

// SendClientBets Sends a batch of client bets to the server and returns the response
func sendClientBetsBatch(senderID string, bets []ClientBet, conn net.Conn) (Message, error) {
	betMessage := NewBetMessage(senderID, bets)
	bytes := BetMessageToBytes(betMessage)
	err := _sendToServer(bytes, conn)
	if err != nil {
		return Message{}, err
	}
	return receiveResponse(conn)
}

// SendDoneSendingBets Sends a message to the server indicating that the client has finished sending bets
func sendDoneSendingBets(senderID string, conn net.Conn) (Message, error) {
	message := NewDoneSendingBetsMessage(senderID)
	bytes := DoneSendingBetsMessageToBytes(message)
	err := _sendToServer(bytes, conn)
	if err != nil {
		return Message{}, err
	}
	return receiveResponse(conn)
}

// SendRequestWinners Sends a message to the server requesting the winners
func sendRequestWinners(senderID string, conn net.Conn) (ResponseWinnersMessage, error) {
	message := NewRequestWinnersMessage(senderID)
	bytes := RequestWinnersMessageToBytes(message)
	err := _sendToServer(bytes, conn)
	if err != nil {
		return ResponseWinnersMessage{}, err
	}
	return receiveWinnersResponse(conn)
}

func _sendToServer(bytes []byte, conn net.Conn) error {
	bytesSent := 0
	for bytesSent < len(bytes) {
		bytesSentAux, err := conn.Write(bytes[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += bytesSentAux
	}
	return nil
}

// ReceiveResponse Receives a response message from the server
func receiveResponse(conn net.Conn) (Message, error) {
	bytes, err := readFromServer(conn)
	if err != nil {
		return Message{}, err
	}
	return MessageFromBytes(bytes), nil
}

// ReceiveWinnersResponse Receives a winners response message from the server
func receiveWinnersResponse(conn net.Conn) (ResponseWinnersMessage, error) {
	bytes, err := readFromServer(conn)
	if err != nil {
		return ResponseWinnersMessage{}, err
	}
	return ResponseWinnersMessageFromBytes(bytes)
}

func readFromServer(conn net.Conn) ([]byte, error) {
	bytes := make([]byte, BUFFER_SIZE)
	bytesReceived := 0
	// Read until a \r\n is received to prevent short reads
	for bytesReceived == 0 || bytes[bytesReceived-2] != byte('\r') || bytes[bytesReceived-1] != byte('\n') {
		bytesReceivedAux, err := conn.Read(bytes[bytesReceived:])
		if err != nil {
			return nil, err
		}
		bytesReceived += bytesReceivedAux
	}
	return bytes[:bytesReceived], nil
}
