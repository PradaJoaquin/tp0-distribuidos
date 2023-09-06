package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type betLoader struct {
	reader    *bufio.Reader
	batchSize int
	file      *os.File
}

// Creates a new bet loader, if an error occurs while opening the file the error is returned
func NewBetLoader(betsPath string, batchSize int) (*betLoader, error) {
	file, err := os.Open(betsPath)
	if err != nil {
		log.Errorf("action: open_bets_file | result: fail | error: %v", err)
		return nil, err
	}
	reader := bufio.NewReader(file)
	return &betLoader{
		reader:    reader,
		batchSize: batchSize,
		file:      file,
	}, nil
}

// Next returns the next batch of bets from the file
func Next(b *betLoader) ([]ClientBet, error) {
	bets := make([]ClientBet, 0)
	for i := 0; i < b.batchSize; i++ {
		line, err := b.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = _StripEndOfLine(line)
		bet, err := _NewClientBetFromCSV(line)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
		if !b.HasNext() {
			break
		}
	}
	return bets, nil
}

// StripEndOfLine removes the end of line characters from a string
func _StripEndOfLine(line string) string {
	str := strings.TrimSuffix(line, "\r\n")
	return strings.TrimSuffix(str, "\n")
}

// NewClientBetFromCSV Creates a new client bet from a CSV line
func _NewClientBetFromCSV(line string) (ClientBet, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 5 {
		return ClientBet{}, fmt.Errorf("invalid number of fields")
	}
	numero, err := strconv.Atoi(fields[4])
	if err != nil {
		return ClientBet{}, fmt.Errorf("invalid number")
	}
	return ClientBet{
		Nombre:     fields[0],
		Apellido:   fields[1],
		Documento:  fields[2],
		Nacimiento: fields[3],
		Numero:     numero,
	}, nil
}

// HasNext returns true if there are more bets to be read from the file
func (b *betLoader) HasNext() bool {
	_, err := b.reader.Peek(1)
	if err != nil {
		return false
	}
	return true
}

// Closes the file, the loader cannot be used after this method is called.
func (b *betLoader) Close() {
	b.file.Close()
}
