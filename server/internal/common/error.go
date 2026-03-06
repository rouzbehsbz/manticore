package common

import (
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
)

func ErrorResponse(session *session.Session, msg string) {
	session.Write(protocol.BuildErrorResponsePacket(msg))
}
