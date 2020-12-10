# Swagger
- https://girigiribauer.com/tech/20190318/
  - これ詳しいし、筋がよさそう
- https://riotz.works/articles/lopburny/2019/08/17/describe-bearer-scheme-in-openapi-3/
-  3.0でのsecurity(Bearer)


# DB
- https://github.com/golang-standards/project-layout/issues/1
  - /configs にmigration(db init)をおいてよさそう。
  - HandleFuncの内側でdbconnectを行う
- https://qiita.com/hiro9/items/e6e41ec822a7077c3568
  - pgweb使うと手軽そう

# DB migration
- https://dev.classmethod.jp/articles/db-migrate-with-golang-migrate/
  - これを見た
```shell
# まずdocker-compose upでdbコンテナを立ち上げておく

# これをしたいがうまくいかない
sudo docker-compose run --rm migration exec migrate create -ext sql -dir /work/migration -seq create_users_table

# 代案
sudo docker-compose -f deployments/docker-compose.dev.yml run --rm migration /bin/bash  # shellに入る
## migrationコンテナでmigration
migrate create -ext sql -dir migration -seq <sql_file_name> # sqlファイルを作成

export POSTGRESQL_URL='postgres://suburi:password@db:5432/suburi_db?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path migrations up # upでmigration、downで切り戻し
```
- volumeでchownしなきゃ問題もあるが動く

# session
- https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/06.0.html
  - sessionを標準パッケージで実装している。
- https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/06.3.html

```go
package memory

import (
    "container/list"
    "github.com/suburi-dev/gowiki/internal/session"
    "sync"
    "time"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
    sid          string                      //session idユニークID
    timeAccessed time.Time                   //最終アクセス時間
    value        map[interface{}]interface{} //sessionに保存される値
}

func (st *SessionStore) Set(key, value interface{}) error {
    st.value[key] = value
    pder.SessionUpdate(st.sid)
    return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
    pder.SessionUpdate(st.sid)
    if v, ok := st.value[key]; ok {
        return v
    } else {
        return nil
    }
    return nil
}

func (st *SessionStore) Delete(key interface{}) error {
    delete(st.value, key)
    pder.SessionUpdate(st.sid)
    return nil
}

func (st *SessionStore) SessionID() string {
    return st.sid
}

type Provider struct {
    lock     sync.Mutex               //ロックに使用します
    sessions map[string]*list.Element //メモリに保存するために使用します
    list     *list.List               //gcを行うために使用します
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    v := make(map[interface{}]interface{}, 0)
    newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
    element := pder.list.PushBack(newsess)
    pder.sessions[sid] = element
    return newsess, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
    if element, ok := pder.sessions[sid]; ok {
        return element.Value.(*SessionStore), nil
    } else {
        sess, err := pder.SessionInit(sid)
        return sess, err
    }
    return nil, nil
}

func (pder *Provider) SessionDestroy(sid string) error {
    if element, ok := pder.sessions[sid]; ok {
        delete(pder.sessions, sid)
        pder.list.Remove(element)
        return nil
    }
    return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
    pder.lock.Lock()
    defer pder.lock.Unlock()

    for {
        element := pder.list.Back()
        if element == nil {
            break
        }
        if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
            pder.list.Remove(element)
            delete(pder.sessions, element.Value.(*SessionStore).sid)
        } else {
            break
        }
    }
}

func (pder *Provider) SessionUpdate(sid string) error {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    if element, ok := pder.sessions[sid]; ok {
        element.Value.(*SessionStore).timeAccessed = time.Now()
        pder.list.MoveToFront(element)
        return nil
    }
    return nil
}

func init() {
    pder.sessions = make(map[string]*list.Element, 0)
    session.Register("memory", pder)
}
```
- sessionストレージというやつらしい。
- メモリ中に保持するので、再起動すると全部消える。

```go
struct SessionStore
SessionStore.Set
SessionStore.Get
SessionStore.Delete
SessionStore.SessionID

struct Provider
Provider.SessionInit
Provider.SessionRead
Provider.SessionDestroy
Provider.SessionGC
Provider.SessionUpdate

init()
```

- SessionStoreはkey-value storeと考えて良い。
- 上のStoreで使われている4つの関数をサポートするためにProviderが存在していて、Init, GC, UpdateのときにMutex Lockをかけている。
- これらはメモリを使う想定だけど、それぞれの関数の中身を差し替えると(たとえばDBセッションなど)動くらしい。
- Cookie

```go
// Cookie使い方
http.SetCookie(w ResponseWriter, cookie *Cookie)

/* 
// Cookieの構造
type Cookie struct {
  Name string
  Value string
  Path string
  Domain string
  Expires time.Time
  RawExpires string

  MaxAge int
  Secure bool
  HttpOnly bool
  Raw string
  Unparsed []string
}
*/

// Set 実際に期限を設定するとき
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie(Name: "username", Value: "uta8a", Expires: expiration)
http.SetCookie(w, &cookie)

// Get cookieからデータ取り出すとき
cookie, _ := r.Cookie("username")
fmt.Fprint(w, cookie)

for _, cookie := range r.Cookies() {
  fmt.Fprint(w, cookie.Name)
}
```

- Session管理

```go
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

var globalSessions *session.Manager

func init() {
  globalSessions, _ = NewManager("memory", "gosessionid", 3600)
}

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

func login(w http.ResponseWriter, r *http.Request) {
  sess := globalSessions.SessionStart(w, r)
  r.ParseForm()
  if r.Method == "GET" {
    t, _ := template.ParseFiles("login.gtpl")
    w.Header().Set("Content-Type", "text/html")
    t.Execute(w, sess.Get("username"))
  } else {
    sess.Set("username". r.Form["username"])
    http.Redirect(w, r, "/", 302)
  }
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
    cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true. Expires: expiration, MaxAge: -1}
    http.SetCookie(w, &cookie)
  }
}

func init() {
  go globalSessions.GC()
}

func (manager *Manager) GC() {
  manager.lock.Lock()
  defer manager.lock.Unlock()
  manager.provider.SessionGC(manager.maxlifetime)
  time.AfterFunc(time.Duration(manager.maxlifetime), func() {manager.GC()})
}
```


# Cookie, Bearer
- Bearerの例。

```yaml
components:
  securitySchemes:
    Bearer:
      type: http
      scheme: bearer
      description: 'credential token for API'
```

- Cookieの例。
- https://swagger.io/docs/specification/authentication/cookie-authentication/

```yaml
components:
  securitySchemes:
    cookieAuth:
      type: apiKey
      in: cookie
      name: SESSIONID
```

# REST
- https://nec-baas.github.io/baas-manual/latest/developer/ja/rest-ref/user/register.html
  - 実例のドキュメントが載っている。
  - signup `/users` POST
  - login `/login` POST
  - logout `/login` DELETE
  - change user info `/users/<userid>` PUT

# Error
- https://tutuz-tech.hatenablog.com/entry/2020/03/26/193519
  - エラーハンドリング

# log
- 2020/12/08
  - https://www.yoheim.net/blog.php?q=20170403
  - net/httpのみで実装してみる
  - swaggerをやって、先にAPIを決める
  - Swaggerを導入し始めた。
  - editorconfigを導入。Makefileに注意
  - API部分はよさそうなのでDBを考える。でもどうせデータはそんな入らないから雑でもいいか。
  - goのディレクトリ構成は https://github.com/golang-standards/project-layout を見ている。
  - DBの組み方、共通化の仕方が分からない。init.sqlで統一してやればいいかも。
  - railsの方でDBの組み方について学びながら直していけばいいかも
  - JWTは後で考えるとして、セッションを素で実装してみることにする。素振りだし失敗したらCTFの問題になるからいいでしょう。
  - testはファイル分割のタイミングで書くとよさそう。
- 2020/12/09
  - http.HandleFuncと、NewServeMuxしてHandleFuncの違いが分かってない。
  - driver localのdocker volumeがホストのどこにファイルが作成されるのか分かってない。
  - ようやくEnvでDB接続URLつくるとこまできた。
  - DB接続、migrationをしないとなあ
  - これ、rails migrateに対応するような、migration toolをAppの外でInitとして使う必要があるのでは？
  - おそらく、compose upしてdbを立ち上げて、そのdbに対してexec migrateをしてmigrationを行い、再度upして立ち上げる感じで後は開発という流れ...？
  - たぶんこの流れでOKっぽい
  - 次はsessionとlogin/logoutを実装する。なんかここ乗り切ればあとはdirectory構造の保持くらいしか鬼門がないし、一気に実装して行ける気がする(フラグ)
  - 通話用に話すネタ帳みたいなのをwikiにおいておくと楽しそう。
  - ここらへんでファイル分割してテスト書きたいな。やります。
  - importは ``github.com/<org>/<repo>/internal`` みたいな絶対パスで書きそう。``_test.go``は同じ階層に置くみたい。
- 2020/12/10
  - テストまだだけどファイル分割は完了した。
  - sessionをやる。
  - 構造体って宣言するとき初期化が部分的でも大丈夫なのかな
