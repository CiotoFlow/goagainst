package irc

import "fmt"

type Entity interface {
	isEntity()
	fmt.Stringer
}

/* Unknown entity */
type Unknown struct {
}

func (e *Unknown) String() string {
	return "(unknown)"
}

func (e *Unknown) isEntity() {}

/* Server */

type Server struct {
	Name string
}

func (e *Server) String() string {
	return e.Name
}

func (e *Server) isEntity() {}

/* User */

type User struct {
	Nick string
	Name string
	Host string
	Valid bool
}

func (e *User) String() string {
	return e.Nick + "!" + e.Name + "@" + e.Host
}

func (e *User) isEntity() {}

/* Channel */

type Channel struct {
	Name string
	Users []*User
}

func (e *Channel) String() string {
	return e.Name
}

func (e *Channel) isEntity() {}

/* Utils */

func IsUnknown(e Entity) bool {
	_, unknown := e.(*Unknown)
	return unknown
}
