package session

import (
  "time"
  "fmt"
  "sync"
  "io"
  "crypto/rand"
  "encoding/base64"
  "net/http"
  "net/url"
)

type Provider interface {
  SessionInit(sid string) (Session, error)
  SessionRead(sid string) (Session, error)
  SessionDestroy(sid string) error
  SessionGC(maxlifeTime int64)
}

type Session interface {
  Set(key, value interface{}) error
  Get(key interface{}) interface{}
  Delete(key interface{}) error
  SessionID() string
}

type Manager struct {
  cookieName string
  lock sync.Mutex
  provider Provider
  maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
  provider, ok := provides[provideName]
  if !ok {
    return nil, fmt.Errorf("session: unknown provide %q (forgotten import ?)", provideName)
  }
  return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

var provides = make(map[string]Provider)

func Register(name string, provider Provider) {
  if provider == nil {
    panic("session: Register provide is nil")
  }
  if _, dup := provides[name]; dup {
    panic("session: Register called twice for provide " + name)
  }
  provides[name] = provider
}

func (manager *Manager) sessionId() string {
  b := make([]byte, 32)
  if _, err := io.ReadFull(rand.Reader, b); err != nil {
    return ""
  }
  return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
  manager.lock.Lock()
  defer manager.lock.Unlock()
  cookie, err := r.Cookie(manager.cookieName)
  if err != nil || cookie.Value == "" {
    sid := manager.sessionId()
    session, _ = manager.provider.SessionInit(sid)
    cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
    http.SetCookie(w, &cookie)
  } else {
    sid, _ := url.QueryUnescape(cookie.Value)
    session, _ = manager.provider.SessionRead(sid)
  }
  return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie(manager.cookieName)
  if err != nil || cookie.Value == "" {
    return
  } else {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    manager.provider.SessionDestroy(cookie.Value)
    expiration := time.Now()
    cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
    http.SetCookie(w, &cookie)
  }
}

func (manager *Manager) GC() {
  manager.lock.Lock()
  defer manager.lock.Unlock()
  manager.provider.SessionGC(manager.maxlifetime)
  time.AfterFunc(time.Duration(manager.maxlifetime), func() {manager.GC()})
}
