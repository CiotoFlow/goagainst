package irc

import (
	"errors"
	"strings"
)

type User struct {
	Nickname string
	Username string
	Hostname string
}

type EntityType int

const (
	UNKNOWN EntityType = iota
	SERVER
	USER
	CHANNEL
)

type Entity struct {
	Type EntityType
	Name string
	// Only for users
	User string
	Host string
}

func (e *Entity) String() string {
	if e.Type == SERVER || e.Type == CHANNEL {
		return e.Name
	} else {
		return e.Name + "!" + e.User + "@" + e.Host
	}
}

type Message struct {
	Entity *Entity
	Command string
	Params []string
	Trailing string
}

func (msg *Message) String() string {
	s := ""
	if (msg.Entity != nil && msg.Entity.Type != UNKNOWN) {
		s = ":" + msg.Entity.String() + " "
	}

	s += msg.Command

	for _, p := range(msg.Params) {
		s += " " + p
	}

	if msg.Trailing != "" {
		s += " :" + msg.Trailing
	}
	
	return s
}

func ParseMessage(line string) (*Message, error) {
	var prefixEnd, trailingStart int
	var msg Message

	if strings.HasPrefix(line, ":") {
		prefixEnd = strings.Index(line, " ")
		if (prefixEnd >= 0) {
			msg.Prefix = line[1:prefixEnd]
		} else {
			return nil, errors.New("Invalid message")
		}
	} else {
		msg.Prefix = ""
		prefixEnd = -1
	}

	trailingStart = strings.Index(line, " :")
	if trailingStart >= 0 {
		msg.Trailing = line[trailingStart + 2:]
	} else {
		msg.Trailing = ""
		trailingStart = len(line)
	}

	middle := strings.Split(line[prefixEnd+1:trailingStart], " ")

	if len(middle) < 1 {
		return nil, errors.New("Invalid message")
	}

	msg.Command = middle[0]

	if len(middle) > 1 {
		msg.Params = middle[1:]
	}

	msg.Entity = UNKNOWN
	if len(msg.Prefix) > 0 {
		bangIndex = strings.Index(msg.Prefix, "!")
		if bangIndex >= 0 {
			msg.Entity.Type = USER
			msg.Entity.Name = msg.Prefix[0:bangIndex]
			
			atIndex = strings.Index(msg.Prefix[bangIndex+1:], "@")
			if atIndex >= 0 {
				msg.Entity.User = msg.Prefix[bangIndex+1:atIndex]
				msg.Entity.Host = msg.Prefix[atIndex+1:]
			} else {
				return nil, errors.New("Invalid prefix %s", msg.Prefix)
			}
		} else {
			msg.Entity.Type = SERVER
			msg.Entity.Name = msg.Prefix
		}
	}

	return &msg, nil
}
