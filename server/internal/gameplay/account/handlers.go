package account

import (
	"context"
	"errors"

	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
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

func (r *RegisterHandler) sendResponse(session *session.Session, ok bool, msg string) {
	session.Write(
		protocol.BuildRegisterResponsePacket(
			ok,
			msg,
		),
	)
}

func (r *RegisterHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_RegisterRequest)

	username := payload.RegisterRequest.Username
	password := payload.RegisterRequest.Password

	var err error
	err = isUsernameValid(username)
	err = isPasswordValid(username)
	if err != nil {
		r.sendResponse(rp.Session, false, err.Error())
	}

	_, err = r.db.Q.GetAccountByUsername(ctx, username)
	if err == nil {
		r.sendResponse(rp.Session, false, "This username is already taken.")
		return
	}

	bytes := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytes, PasswordSalt)

	_, _ = r.db.Q.CreateAccount(ctx, sources.CreateAccountParams{
		Username: username,
		Password: string(hashedPassword),
	})

	r.sendResponse(rp.Session, true, "Your account has been registered successfully.")
}

type LoginHandler struct {
	db *db.Db
}

func NewLoginHandler(db *db.Db) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func (l *LoginHandler) sendResponse(session *session.Session, ok bool, msg string) {
	session.Write(
		protocol.BuildLoginResponsePacket(
			ok,
			msg,
		),
	)
}

func (l *LoginHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_LoginRequest)

	username := payload.LoginRequest.Username
	password := payload.LoginRequest.Password

	var err error
	err = isUsernameValid(username)
	err = isPasswordValid(username)
	if err != nil {
		l.sendResponse(rp.Session, false, err.Error())
	}

	acc, err := l.db.Q.GetAccountByUsername(ctx, username)
	if err != nil {
		l.sendResponse(rp.Session, false, "Username or password is incorrect.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		l.sendResponse(rp.Session, false, "Username or password is incorrect.")
		return
	}

	l.sendResponse(rp.Session, true, "You have logged in successfully.")
}

func isUsernameValid(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return errors.New("Username length must be between 3 and 50 characters")
	}

	if !UsernameRegex.MatchString(username) {
		return errors.New("Username must not contain illigal characters")
	}

	return nil
}

func isPasswordValid(password string) error {
	if len(password) < 8 || len(password) > 128 {
		return errors.New("Password length must be at least 8 characters.")
	}

	return nil
}
