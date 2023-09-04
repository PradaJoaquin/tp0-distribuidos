package common

import "fmt"

const (
	// ResponseType
	AckMessage = 0
	ErrMessage = 1
)

type BetMessage struct {
	SenderID string
	bet      ClientBet
}

type ResponseMessage struct {
	SenderID     string
	ResponseType int
}

func NewBetMessage(senderID string, bet ClientBet) BetMessage {
	return BetMessage{
		SenderID: senderID,
		bet:      bet,
	}
}

// BetMessageToBytes Converts a bet message to a byte array to be sent through the socket
//
// Protocol: <sender_id>:<nombre>,<apellido>,<documento>,<nacimiento>,<numero>\n
func BetMessageToBytes(betMessage BetMessage) []byte {
	return []byte(
		fmt.Sprintf(
			"%s:%s,%s,%s,%s,%d\n",
			betMessage.SenderID,
			betMessage.bet.Nombre,
			betMessage.bet.Apellido,
			betMessage.bet.Documento,
			betMessage.bet.Nacimiento,
			betMessage.bet.Numero,
		),
	)
}

// ResponseMessageFromBytes Converts a byte array to a response message received through the socket
//
// Protocol: <sender_id>:<response_type>\n
func ResponseMessageFromBytes(bytes []byte) ResponseMessage {
	var senderID string
	var responseType int
	fmt.Sscanf(string(bytes), "%s:%d\n", &senderID, &responseType)
	return ResponseMessage{
		SenderID:     senderID,
		ResponseType: responseType,
	}
}
