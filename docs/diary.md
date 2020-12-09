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
