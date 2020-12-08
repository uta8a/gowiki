# Swagger
- https://girigiribauer.com/tech/20190318/
  - これ詳しいし、筋がよさそう
- https://riotz.works/articles/lopburny/2019/08/17/describe-bearer-scheme-in-openapi-3/
-  3.0でのsecurity(Bearer)


# DB
- https://github.com/golang-standards/project-layout/issues/1
  - /configs にmigration(db init)をおいてよさそう。
  - HandleFuncの内側でdbconnectを行う

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
- 2o2o/12/08
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
