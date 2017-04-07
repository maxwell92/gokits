package session

import (
	mylog "gokits/log"
	"encoding/json"
	"errors"

	localtime "gokits/time"
	"github.com/pborman/uuid"

	redigo "gokits/redigo"
	redis "github.com/garyburd/redigo/redis"
	"sync"
)

var log = mylog.Log

const (
	DEFAULT_EXPIRATION = "86400" // 24*60*60s
)

type Session struct {
	SessionId  string   `json:"sessionId"`
	UserId     string   `json:"userId"`
	UserName   string   `json:"userName"`
	OrgId      string   `json:"orgId"`
	DcList     []string `json:"dcList"`
	CreatedAt  string   `json:"CreatedAt"`
	Expiration string   `json:"Expiration"` // expiration in seconds
}

func NewSession(userId, userName, orgId string) *Session {
	return &Session{
		SessionId:  uuid.New(),
		UserId:     userId,
		UserName:   userName,
		OrgId:      orgId,
		CreatedAt:  localtime.NewLocalTime().String(),
		Expiration: DEFAULT_EXPIRATION,
	}
}

func (s *Session) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), s)

	if err != nil {
		log.Errorf("Session DecodeJson Error: err=%s", err)
		return err
	}

	return nil
}

func (s *Session) EncodeJson() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		log.Errorf("Session EncodeJson Error: err=%s", err)
		return "", err
	}
	return string(data), nil
}

var instance *SessionStore

var once sync.Once

type SessionStore struct {
	pool *redis.Pool
}

func SessionStoreInstance() *SessionStore {
	return instance
}

func NewSessionStore() *SessionStore {
	once.Do(func() {
		instance = &SessionStore{
			pool: redigo.NewRedisClient(),
		}
	})
	return instance
}

func (ss *SessionStore) ValidateOrgId(sessionIdClient string, OrgIdClient string) (bool, error) {
	session, err := ss.Get(sessionIdClient)
	if err != nil {
		log.Errorf("Get session from sessionIdClient error: sessionIdClient: %s, err=%s", sessionIdClient, err)
		return false, err
	}

	// sessionId invalid
	if session == nil && err == nil {
		return false, errors.New("Validate sessionIdClient failed: invalid sessionIdClient")
	}

	if session.OrgId == OrgIdClient {
		return true, nil
	} else {
		return false, errors.New("Validate sessionId failed: OrgId doesn't match")
	}

}

func (ss *SessionStore) ValidateUserId(sessionIdClient string, UserIdClient string) (bool, error) {

	session, err := ss.Get(sessionIdClient)
	if err != nil {
		log.Errorf("Get session from sessionIdClient error: sessionIdClient: %s, err=%s", sessionIdClient, err)
		return false, err
	}

	// sessionId invalid
	if session == nil && err == nil {
		return false, errors.New("Validate sessionIdClient failed: invalid sessionIdClient")
	}
	if session.UserId == UserIdClient {
		return true, nil
	} else {
		return false, errors.New("Validate sessionId failed: UserId doesn't match")
	}
}

func (ss *SessionStore) Get(sessionId string) (*Session, error) {

	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return nil, errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	// If exists
	exists, err := ss.Exist(sessionId)
	if err != nil {
		log.Fatalf("SessionStore exist error: sessionId=%s, err=%s", sessionId, err)
		return nil, err
	}

	// not exists
	if !exists {
		log.Warnf("The Session not exists: sessionId=%s", sessionId)
		return nil, nil
	}

	// exists
	session := &Session{}

	data, err := redis.Bytes(conn.Do("GET", sessionId))
	if err != nil {
		log.Fatalf("Redis Get error: sessionId=%s, err=%s", sessionId, err)
		return nil, err
	}

	err = json.Unmarshal(data, session)

	if err != nil {
		log.Fatalf("Json unmashal failed: data=%s, err=%s", string(data))
		return nil, err
	}

	return session, nil
}

func (ss *SessionStore) Set(session *Session) error {
	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalf("The Connection is nil: conn := ss.pool.Get()")
		return errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	data, err := json.Marshal(session)

	if err != nil {
		log.Fatalf("Json marshal err: err=%s", err)
		return err
	}

	_, err = conn.Do("SET", session.SessionId, data, "EX", session.Expiration)

	if err != nil {
		log.Errorf("Redis set error: sessionId=%s, err=%s", session.SessionId, err)
		return err
	}

	return nil
}

func (ss *SessionStore) Delete(sessionId string) error {

	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	_, err := conn.Do("DEL", sessionId)
	if err != nil {
		log.Errorf("Redis delete error: sessionId=%s, err=%s", sessionId, err)
		return err
	}

	return nil
}

func (ss *SessionStore) Exist(sessionId string) (bool, error) {
	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return false, errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", sessionId))

	if err != nil {
		log.Fatalf("Redis Bool error: sessionId=%s, err=%s\n", exists, err)
		return false, err
	}

	return exists, nil
}
