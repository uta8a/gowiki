# gowiki
- ドキュメント置き場
- 継続的なデプロイと開発

# TODO
- error処理ちゃんとやる(サーバがエラーを愚直に返している)
- テストを書く

# 開発手順
- cloneしてくる
- Makefileは``sudo``つけてdocker-compose upとかしています...
```shell
# docker-compose build
$ make dev
# docker-compose up
$ make up # hot-reload server & front

# dbが立ち上がったら
$ make migrate # migrateコンテナに入る
# in docker
## 接続
psql -h db -p 5432 -U suburi -d suburi_db

## migrationコンテナでmigration
migrate create -ext sql -dir migration -seq <sql_file_name> # sqlファイルを作成

## migration up
export POSTGRESQL_URL='postgres://suburi:password@db:5432/suburi_db?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path migration up # upでmigration、downで切り戻し
```
- 基本は上の２つをやってサーバは ``internal/`` を、フロントは ``web/pages`` を見るとOK
- ``localhost:3000/``がトップページ

# document
- ``/api``にOpenAPIのAPI定義書を置いています。
- ``/docs/diary.md``は開発日記です。メモなのでそのうち消します。
