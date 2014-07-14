/*
 * RFCs:
 *   https://tools.ietf.org/html/rfc2810
 *   https://tools.ietf.org/html/rfc2811
 *   https://tools.ietf.org/html/rfc2812
 *   https://tools.ietf.org/html/rfc2813
 */
package irc

import (
	"bufio"
	"net/textproto"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
        IRC_EVENT_NICK = 0
	IRC_EVENT_PRIV_MSG = 1
	/* TODO: define more... */
)

type IRC struct {
	conn bool
	nick string
	server string
	channel string
	useTls bool
	sock net.Conn
}

func NewIRC(nick, server, channel string, useTls bool) *IRC {
	return &IRC {
		nick: nick,
		server: server,
		channel: channel,
		useTls: useTls,
	}
}

func (irc *IRC) AddCallback() error {
	/* TODO */
	return nil
}

func (irc *IRC) Loop() error {
	var cmd string

	/* TODO: try to use textproto.Dial */
	reader := bufio.NewReader(irc.sock)
	tp := textproto.NewReader(reader)

	if (!irc.conn) {
		return errors.New("not connected")
	}

	irc.sock.Write([]byte("NICK " + irc.nick + "\r\n"))
	irc.sock.Write([]byte("USER " + irc.nick + " 0 * :Stocazzo\r\n"))

	for {
		line, _ := tp.ReadLine()
		fmt.Println(line)

		resp := strings.Split(line, " ")

		/* XXX: handle better */
		if (resp[0][0] == ':' && len(resp) > 1) {
			cmd = resp[1]
		} else {
			cmd = resp[0]
		}

		if (cmd == "PING") {
			irc.sock.Write([]byte("PONG " + resp[1] + "\r\n"))
		} else if (cmd == "001") {
			irc.sock.Write([]byte("JOIN " + irc.channel + "\r\n"))
		} else if (cmd == "443") {
			/* Duplicated NICK */
		} else if (cmd == "PRIVMSG") {
			/* TODO */
		}

		/* TODO: add more commands */
	}

	return nil
}

func (irc *IRC) Disconnect() {
	/* XXX TODO: send QUIT? */
	irc.sock.Close()
}

func (irc *IRC) Connect() error {
	var err error

	if (irc.conn) {
		return errors.New("already connected")
	}

	if (irc.useTls) {
		irc.sock, err = tls.Dial("tcp", irc.server, &tls.Config{InsecureSkipVerify: true})
	} else {
		irc.sock, err = net.Dial("tcp", irc.server)
	}

	if (err == nil) {
		irc.conn = true
	}

	return err
}
