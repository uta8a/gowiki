# gowiki
- golangで素振りする
- CTF Score Serverの前に素振りしておきたかったので

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
```
- 基本は上の２つをやってサーバは ``internal/`` を、フロントは ``web/pages`` を見るとOK
- ``localhost:3000/``がトップページ
