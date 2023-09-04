package common

import "net"

// SendClientBet Sends a client bet to the server and returns the response message
func sendClientBet(senderID string, bet ClientBet, conn net.Conn) (ResponseMessage, error) {
	betMessage := NewBetMessage(senderID, bet)
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
	bytes := make([]byte, 1024)
	bytesReceived := 0
	// Read until a newline is received to prevent short reads
	for bytesReceived == 0 || bytes[bytesReceived-1] == byte('\n') {
		bytes := make([]byte, 1024)
		newBytesReceived, err := conn.Read(bytes)
		if err != nil {
			return ResponseMessage{}, err
		}
		bytesReceived += newBytesReceived
	}
	return ResponseMessageFromBytes(bytes[:bytesReceived]), nil
}
