package account

import (
	"context"
	"errors"
	"regexp"

	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/character"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"github.com/rouzbehsbz/zurvan"
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

func (r *RegisterHandler) successRes(session *session.Session) {
	session.Write(protocol.BuildRegisterResPacket())
}

func (r *RegisterHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_RegisterReq)

	username := payload.RegisterReq.Username
	password := payload.RegisterReq.Password

	if err := isUsernameValid(username); err != nil {
		common.ErrorRes(rp.Session, err.Error())
		return
	}
	if err := isPasswordValid(password); err != nil {
		common.ErrorRes(rp.Session, err.Error())
		return
	}

	_, err := r.db.Q.GetAccountByUsername(ctx, username)
	if err == nil {
		common.ErrorRes(rp.Session, "This username is already taken.")
		return
	}

	bytes := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytes, PasswordSalt)

	_, _ = r.db.Q.CreateAccount(ctx, sources.CreateAccountParams{
		Username: username,
		Password: string(hashedPassword),
	})

	r.successRes(rp.Session)
}

type LoginHandler struct {
	db *db.Db
}

func NewLoginHandler(db *db.Db) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func (l *LoginHandler) successRes(session *session.Session) {
	session.Write(protocol.BuildLoginResPacket())
}

func (l *LoginHandler) Handle(rp session.ReceivedPacket) {
	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_LoginReq)

	username := payload.LoginReq.Username
	password := payload.LoginReq.Password

	if err := isUsernameValid(username); err != nil {
		common.ErrorRes(rp.Session, err.Error())
		return
	}
	if err := isPasswordValid(password); err != nil {
		common.ErrorRes(rp.Session, err.Error())
		return
	}

	acc, err := l.db.Q.GetAccountByUsername(ctx, username)
	if err != nil {
		common.ErrorRes(rp.Session, "Username or password is incorrect.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		common.ErrorRes(rp.Session, "Username or password is incorrect.")
		return
	}

	rp.Session.AccountId = uint32(acc.ID)

	l.successRes(rp.Session)
}

type MyCharactersListHandler struct {
	db *db.Db
}

func NewMyCharactersListHandler(db *db.Db) *MyCharactersListHandler {
	return &MyCharactersListHandler{
		db: db,
	}
}

func (m *MyCharactersListHandler) successRes(session *session.Session, myCharacters []*protocol.MyCharacter) {
	session.Write(protocol.BuildMyCharactersListResPacket(myCharacters))
}

func (m *MyCharactersListHandler) Handle(rp session.ReceivedPacket) {
	if !rp.Session.IsAuthenticated() {
		common.ErrorRes(rp.Session, common.UnauthorizedErrorMsg)
		return
	}

	ctx := context.Background()

	rawCharacters, _ := m.db.Q.GetCharactersByAccountId(ctx, int32(rp.Session.AccountId))
	characters := make([]*protocol.MyCharacter, len(rawCharacters))

	for _, character := range rawCharacters {
		characters = append(characters, models.MapCharacterToMyCharacterPacket(character))
	}

	m.successRes(rp.Session, characters)
}

type CharacterCreateHandler struct {
	db *db.Db
}

func NewCharacterCreateHandler(db *db.Db) *CharacterCreateHandler {
	return &CharacterCreateHandler{
		db: db,
	}
}

func (c *CharacterCreateHandler) successRes(session *session.Session) {
	session.Write(protocol.BuildCharacterCreateResPacket())
}

func (c *CharacterCreateHandler) Handle(rp session.ReceivedPacket) {
	if !rp.Session.IsAuthenticated() || rp.Session.IsCharacterSelected() {
		common.ErrorRes(rp.Session, common.UnauthorizedErrorMsg)
		return
	}

	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_CharacterCreateReq)

	nickname := payload.CharacterCreateReq.Nickname

	_, err := c.db.Q.CreateCharacter(ctx, sources.CreateCharacterParams{
		Nickname:     nickname,
		Vitality:     12,
		Intelligence: 12,
		Willpower:    12,
		Dexterity:    12,
		Spirit:       12,
	})
	if err != nil {
		common.ErrorRes(rp.Session, "This nickname is already taken.")
	}

	c.successRes(rp.Session)
}

type CharacterJoinHandler struct {
	db    *db.Db
	world *zurvan.World
}

func NewCharacterJoinHandler(db *db.Db, world *zurvan.World) *CharacterJoinHandler {
	return &CharacterJoinHandler{
		db:    db,
		world: world,
	}
}

func (c *CharacterJoinHandler) successRes(session *session.Session) {
	session.Write(protocol.BuildCharacterJoinResPacket())
}

func (c *CharacterJoinHandler) Handle(rp session.ReceivedPacket) {
	if !rp.Session.IsAuthenticated() || rp.Session.IsCharacterSelected() {
		common.ErrorRes(rp.Session, common.UnauthorizedErrorMsg)
		return
	}

	ctx := context.Background()
	payload := rp.Packet.Payload.(*protocol.Packet_CharacterJoinReq)

	id := payload.CharacterJoinReq.Id

	char, err := c.db.Q.GetCharacterById(ctx, int32(id))
	if err != nil {
		common.ErrorRes(rp.Session, "Invalid character.")
		return
	}

	rp.Session.CharacterId = uint32(char.ID)

	entity := c.world.Spawn()

	c.world.EmitEvents(
		character.JoinsWorldEvent{
			Character: entity,
			Id:        uint32(char.ID),
		},
	)

	c.successRes(rp.Session)
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
