package irc

import (
	"errors"
	"strings"
)

type Message struct {
	Prefix string
	Command string
	Params []string
	Trailing string
}

func ParseMessage(line string) (*Message, error) {
	var prefixEnd, trailingStart int
	var msg Message

	if (strings.HasPrefix(line, ":")) {
		prefixEnd = strings.Index(line, " ")
		if (prefixEnd >= 0) {
			msg.Prefix = line[1:prefixEnd]
		} else {
			return nil, errors.New("Invalid message")
		}
	} else {
		prefixEnd = -1
	}

	trailingStart = strings.Index(line, " :")
	if (trailingStart >= 0) {
		msg.Trailing = line[trailingStart + 2:]
	} else {
		trailingStart = len(line)
	}

	middle := strings.Split(line[prefixEnd+1:trailingStart], " ")

	if (len(middle) < 1) {
		return nil, errors.New("Invalid message")
	}

	msg.Command = middle[0]

	if (len(middle) > 1) {
		msg.Params = middle[1:]
	}

	return &msg, nil
}
