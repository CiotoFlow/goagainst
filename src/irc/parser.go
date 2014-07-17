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
	INVALID EntityType = iota
	UNKNOWN
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
	if e.Type == USER {
		return e.Name + "!" + e.User + "@" + e.Host
	} else {
		return e.Name
	}
}

type Message struct {
	Entity *Entity // Must never be null
	Command string
	Params []string
	Trailing string
}

func (msg *Message) String() string {
	s := ""
	if msg.Entity.Type != UNKNOWN {
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
	var prefix string
	var prefixEnd, trailingStart int
	var msg Message

	if strings.HasPrefix(line, ":") {
		prefixEnd = strings.Index(line, " ")
		if (prefixEnd >= 0) {
			prefix = line[1:prefixEnd]
		} else {
			return nil, errors.New("Invalid message")
		}
	} else {
		prefix = ""
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

	msg.Entity = new(Entity)
	msg.Entity.Type = UNKNOWN
	if len(prefix) > 0 {
		bangIndex := strings.Index(prefix, "!")
		if bangIndex >= 0 {
			msg.Entity.Type = USER
			msg.Entity.Name = prefix[0:bangIndex]
			
			atIndex := strings.Index(prefix[bangIndex+1:], "@")
			if atIndex >= 0 {
				msg.Entity.User = prefix[bangIndex+1:bangIndex+atIndex+1]
				msg.Entity.Host = prefix[bangIndex+atIndex+2:]
			} else {
				return nil, errors.New("Invalid prefix: "+prefix)
			}
		} else {
			msg.Entity.Type = SERVER
			msg.Entity.Name = prefix
		}
	}

	return &msg, nil
}
