package account

import (
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
)

type RegisterHandler struct {
	db *db.Db
}

func NewRegisterHandler(db *db.Db) *RegisterHandler {
	return &RegisterHandler{
		db: db,
	}
}

func (r *RegisterHandler) Handle(rp session.ReceivedPacket) {

}

type LoginHandler struct {
	db *db.Db
}

func NewLoginHandler(db *db.Db) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func (l *LoginHandler) Handle(rp session.ReceivedPacket) {

}
