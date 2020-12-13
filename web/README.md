# web フロントエンド
- あとですべてTypeScriptで書き直す

```
/ # swagger定義書、signup、loginへのリンク
  login/ # login画面 Authが切れてたら飛ばされるページ
  ping/ # serverとの接続を見るping
  private/ # loginしていないと見られないページ
  signup/ # ユーザ登録画面
  user/ # login後に飛ばされる、自分の所属グループ一覧と記事へのリンク一覧が見れるページ
  article/
    [id]/ # article/:id を表示する
  edit/
    [id] # article/:id を編集する 内容だけ編集できる
  new/ # article を新規作成する
  group/ # group を新規作成する
    [name] # groupのメンバー編集 # APIも実装する
```
