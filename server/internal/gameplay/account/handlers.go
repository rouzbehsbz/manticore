package account

import (
	"context"

	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordSalt = 12
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
	// ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_RegisterRequest)

	// _, err := r.db.Q.GetAccountByUsername(ctx, payload.RegisterRequest.Username)
	// if err == nil {
	// 	rp.Session.Write(
	// 		protocol.BuildRegisterResponsePacket(
	// 			false,
	// 			"This username is already taken.",
	// 		),
	// 	)
	// 	return
	// }

	bytes := []byte(payload.RegisterRequest.Password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytes, PasswordSalt)

	println(string(hashedPassword))

	// _, _ = r.db.Q.CreateAccount(ctx, sources.CreateAccountParams{
	// 	Username: payload.RegisterRequest.Username,
	// 	Password: string(hashedPassword),
	// })

	rp.Session.Write(
		protocol.BuildRegisterResponsePacket(
			true,
			"Your account has been registered successfully.",
		),
	)
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
	_ = context.Background()
	_ = rp.Packet.Payload.(*protocol.Packet_LoginRequest)

	// acc, err := l.db.Q.GetAccountByUsername(ctx, payload.LoginRequest.Username)
	// if err != nil {
	// 	rp.Session.Write(
	// 		protocol.BuildLoginResponsePacket(
	// 			false,
	// 			"Username or password is incorrect.",
	// 		),
	// 	)
	// 	return
	// }

	// if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(payload.LoginRequest.Password)); err != nil {
	// 	rp.Session.Write(
	// 		protocol.BuildLoginResponsePacket(
	// 			false,
	// 			"Username or password is incorrect.",
	// 		),
	// 	)
	// 	return
	// }

	rp.Session.Write(
		protocol.BuildLoginResponsePacket(
			true,
			"You have logged in successfully.",
		),
	)
}
