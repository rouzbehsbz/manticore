package account

import (
	"context"
	"errors"
	"regexp"

	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordSalt int = 12
)

var (
	UsernameRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,49}$`)
)

type RegisterHandler struct {
	db *db.Db
}

func NewRegisterHandler(db *db.Db) *RegisterHandler {
	return &RegisterHandler{
		db: db,
	}
}

func (r *RegisterHandler) successResponse(session *session.Session) {
	session.Write(protocol.BuildRegisterResponsePacket())
}

func (r *RegisterHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_RegisterRequest)

	username := payload.RegisterRequest.Username
	password := payload.RegisterRequest.Password

	if err := isUsernameValid(username); err != nil {
		common.ErrorResponse(rp.Session, err.Error())
		return
	}
	if err := isPasswordValid(password); err != nil {
		common.ErrorResponse(rp.Session, err.Error())
		return
	}

	_, err := r.db.Q.GetAccountByUsername(ctx, username)
	if err == nil {
		common.ErrorResponse(rp.Session, "This username is already taken.")
		return
	}

	bytes := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytes, PasswordSalt)

	_, _ = r.db.Q.CreateAccount(ctx, sources.CreateAccountParams{
		Username: username,
		Password: string(hashedPassword),
	})

	r.successResponse(rp.Session)
}

type LoginHandler struct {
	db *db.Db
}

func NewLoginHandler(db *db.Db) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func (l *LoginHandler) successResponse(session *session.Session) {
	session.Write(protocol.BuildLoginResponsePacket())
}

func (l *LoginHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_LoginRequest)

	username := payload.LoginRequest.Username
	password := payload.LoginRequest.Password

	if err := isUsernameValid(username); err != nil {
		common.ErrorResponse(rp.Session, err.Error())
		return
	}
	if err := isPasswordValid(password); err != nil {
		common.ErrorResponse(rp.Session, err.Error())
		return
	}

	acc, err := l.db.Q.GetAccountByUsername(ctx, username)
	if err != nil {
		common.ErrorResponse(rp.Session, "Username or password is incorrect.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		common.ErrorResponse(rp.Session, "Username or password is incorrect.")
		return
	}

	rp.Session.AccountId = uint32(acc.ID)

	l.successResponse(rp.Session)
}

type MyCharactersListHandler struct {
	db *db.Db
}

func (m *MyCharactersListHandler) successResponse(session *session.Session, myCharacters []*protocol.MyCharacter) {
	session.Write(protocol.BuildMyCharactersListResponsePacket(myCharacters))
}

func (m *MyCharactersListHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()

	if !rp.Session.IsAuthenticated() {
		common.ErrorResponse(rp.Session, "You're not authenticated.")
		return
	}

	rawCharacters, _ := m.db.Q.GetCharactersByAccountId(ctx, int32(rp.Session.AccountId))
	characters := make([]*protocol.MyCharacter, len(rawCharacters))

	for _, character := range rawCharacters {
		characters = append(characters, models.MapCharacterToMyCharacterPacket(character))
	}

	m.successResponse(rp.Session, characters)
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
