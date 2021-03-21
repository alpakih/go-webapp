package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"log"
)

const IDSession = "id"

type FlashMessage struct {
	Type    string
	Message string
	Data    interface{}
}

type UserInfo struct {
	ID         string   `json:"id" form:"id"`
	Name       string   `json:"name" form:"name"`
	Email      string   `json:"email" form:"email"`
	RoleSlug   string   `json:"role_slug" form:"role_slug"`
	Image      string   `json:"image" form:"image"`
	Permission []string `json:"permission" form:"permission"`
}

type Manager struct {
	store    *sessions.CookieStore
	valueKey string
}

func NewSessionManager(store *sessions.CookieStore) *Manager {
	s := new(Manager)
	s.valueKey = "data"
	s.store = store

	return s
}

func (s *Manager) Get(c echo.Context, name string) (interface{}, error) {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if val, ok := session.Values[s.valueKey]; ok {
		return val, nil
	} else {
		return nil, nil
	}
}

func (s *Manager) Set(c echo.Context, name string, value interface{}) error {
	session, _ := s.store.Get(c.Request(), name)
	session.Values[s.valueKey] = value

	err := session.Save(c.Request(), c.Response())
	return err
}

func (s *Manager) Delete(c echo.Context, name string) error {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(c.Request(), c.Response())
}

func (s *Manager) GetWithKeyValues(c echo.Context, name string, keyValue string) (interface{}, error) {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if val, ok := session.Values[keyValue]; ok {
		return val, nil
	} else {
		return nil, nil
	}
}

func (s *Manager) SetFlashMessage(c echo.Context, message string, key string, data interface{}) {
	session, err := s.store.Get(c.Request(), "flash-message")
	if err != nil {
		panic(err)
	}
	mapMessage := FlashMessage{
		Type:    key,
		Message: message,
		Data:    data,
	}
	session.AddFlash(mapMessage)
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		panic(err)
	}
}

func (s *Manager) GetFlashMessage(c echo.Context) FlashMessage {
	session, err := s.store.Get(c.Request(), "flash-message")
	if err != nil {
		return FlashMessage{}
	}
	fm := session.Flashes()
	var flash FlashMessage
	if len(fm) > 0 {
		flash = fm[0].(FlashMessage)
	}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		log.Fatal("ERROR GET FLASH MESSAGE ", err.Error())
	}
	return flash
}

func (s *Manager) SetUserInfo(ID string, name string, email string, roleSlug string, image string, permission []string) UserInfo {
	return UserInfo{
		ID:         ID,
		Name:       name,
		Email:      email,
		RoleSlug:   roleSlug,
		Image:      image,
		Permission: permission,
	}
}
