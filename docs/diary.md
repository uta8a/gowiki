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
