package common

import (
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
)

const (
	UnauthorizedErrorMsg = "You're not authenticated."
)

func ErrorRes(session *session.Session, msg string) {
	session.Write(protocol.BuildErrorResPacket(msg))
}
