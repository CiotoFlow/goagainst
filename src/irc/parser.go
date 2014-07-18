package irc

import (
	"errors"
	"strings"
	"fmt"
)

type Entity interface {
	isEntity()
	fmt.Stringer
}

type Server struct {
	Name string
}

func (e *Server) isEntity() {}

func (e *Server) String() string {
	return e.Name
}

type Unknown struct {
}

func (e *Unknown) String() string {
	return "(unknown)"
}

func (e *Unknown) isEntity() {}

type User struct {
	Nick string
	Name string
	Host string
}

func (e User) isEntity() {}

func (e *User) String() string {
	return e.Nick + "!" + e.Name + "@" + e.Host
}

type Message struct {
	Entity Entity // Must never be null
	Command string
	Params []string
	Trailing string
}

func (msg *Message) String() string {
	s := ""
	_, unknown := msg.Entity.(*Unknown)
	if !unknown {
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

	if len(prefix) > 0 {
		bangIndex := strings.Index(prefix, "!")
		if bangIndex >= 0 {
			atIndex := strings.Index(prefix[bangIndex+1:], "@")
			if atIndex >= 0 {
				msg.Entity = &User{
					prefix[0:bangIndex],
					prefix[bangIndex+1:bangIndex+atIndex+1],
					prefix[bangIndex+atIndex+2:]}
			} else {
				return nil, errors.New("Invalid prefix: "+prefix)
			}
		} else {
			msg.Entity = &Server{prefix}
		}
	}

	return &msg, nil
}
