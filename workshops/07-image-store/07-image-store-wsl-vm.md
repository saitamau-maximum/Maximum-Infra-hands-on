## この章でやること
- ストレージサービスについて知る
- MinIOをセットアップしてみる

## アイコンがないと、チャットがつまらん
開発をしているうちに、Nさんは「インフラチャット」の虜になってしまったようです。

「会社用アカウント」「趣味用アカウント」「愚痴用アカウント」の３つを使ってチャットを最大限楽しんでいるようでしたが……

> N「やっぱり、アイコンがないとどのアカウントでログインしているか分かりづらいよ。社用アカウントで愚痴を吐かないためにも、早急に実装の必要があるね」

あなたはNさんの要望に答えるために、ユーザーが好きなアイコンを設定できる構成を考えなくてはならなくなりました。

## データベースに入れる？　ご乱心ですか
あなたは最初に、MySQLの中に画像データを入れることができないか、と考えました。

さっそくNさんに提案しましたが、どうも反応が芳しくありません。

> Nさん「できれば、他のストレージサービスを使ってほしいかも」

どうやらNさんはデータベースに画像のバイナリデータを入れたくないようです。

## MySQLに画像バイナリを入れるのを避けるべき理由
[https://dev.mysql.com/doc/refman/8.0/en/blob.html](https://dev.mysql.com/doc/refman/8.0/en/blob.html)

MySQLの公式サイトでもこの話は書かれています（最後の部分）

MySQLでは、可変長バイナリ・charを扱えるBLOBやTEXTという方法がありますが、これに画像を入れるのは望ましくないとされています。

主な原因はパフォーマンスの低下でしょう。

クエリでBLOBなどを使うと、メモリではなく、ディスクを使用したテーブルが作成されます。

これによって、RDBMSの強みが弱まります。取得クエリが非常に遅くなったりするのです。

よって、代わりにストレージサービスを使うのが良さそうです。

## ストレージサービス３種

以下の２つの記事をかいつまんで紹介します。

深く学びたい方は本文を参照してください。

[AWS　オブジェクトストレージとは何ですか?](https://aws.amazon.com/jp/what-is/object-storage/)

[PURE STORAGE　オブジェクト・ストレージとは](https://www.purestorage.com/jp/knowledge/what-is-object-storage.html)

ストレージには主に3種のサービスがあります。

1. ファイル・ストレージ

かなり一般的なストレージ形式。直接接続型のもの（DAS）、ネットワーク接続型のもの（NAS）で使用されているようです。

普段扱っているようにディレクトリなどの階層構造でファイルを管理する方法です。

パスを使ってアクセスします。

2. ブロック・ストレージ

データを固定長サイズに分割し（それぞれをブロックと呼ぶ）、それぞれにIDをふります。

「固定長サイズ」というのが重要で、これにより、最適な場所に差し込むように保存することが可能になり、ソースの利用率が高まります。

3. オブジェクト・ストレージ

「オブジェクト」と呼ばれる単位でデータを保存します。（データ＋メタデータ＋ID（キー）で構成されます）

「オブジェクト」単位で管理しているがゆえに、クラウド技術と非常に相性が良く、またメタデータによる細かな制限操作もできるため現在注目されています。


今回は、ローカルで使えるオブジェクト・ストレージの`MinIO`を使ってみましょう。

## MinIOサーバー導入
まずはMinIOのサーバー（単一ノードのやり方）を導入しましょう。

[https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd](https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd)

このサイトに従って行います。

まずは、サーバーをインストールします。Debian系（Ubuntuなど）の場合は以下のようにしてインストールできます。
```bash
# インストール実行ファイルをminio.debという名前でURLから取得
wget https://dl.min.io/server/minio/release/linux-amd64/archive/minio_20250422221226.0.0_amd64.deb -O minio.deb
# ファイル実行
sudo dpkg -i minio.deb
```

公式ドキュメントには、バイナリでインストールした人用に`.service`ファイルの作り方が書いてありますが、上記のやり方をした場合には不要です。

一応、正しくサービスが動いているか確認しましょう
```bash
sudo systemctl status minio.service
```
何かしら表示されたらOKです。

次に、環境変数を設定していきます。`/etc/default/minio`というファイルを作ると、中の環境変数を読み取ってくれるようになります。（ファイルのテンプレートはこのファイルと同じディレクトリにあります）

おもにサーバーのルートユーザーとそのパスワードを設定していきます。

`MINIO_ROOT_USER`にはルートユーザー名、`MINIO_ROOT_PASSWORD`にはそのパスワードを入れてください。

基本的に何でもいいですが、必ず例示のモノとは変えてください。

（デフォルトのものは攻撃を受けやすいため）

また、デフォルトで保存に使うディレクトリを`MINIO_VOLUMES="/mnt/data"`として設定しています。このディレクトリを作り、場合によっては権限を付与しておく必要があります

ここまでできたら、MinIOサーバーを起動しましょう。

```bash
sudo systemctl start minio.service
```

正しく起動できたか状態を確認しましょう

```bash
sudo systemctl status minio.service
```

動いていたら、PC起動時に自動で立ち上がるようにしましょう

```bash
sudo systemctl enable minio.service
```
ここまでできたらサーバーの導入は完了です。

## MinIOクライアント導入
次に、MinIOサーバーを使うためのクライアントを導入します。

GUIでの設定方法もありますが、あえてCLIでの方法をとります。GUIが気になる方はサイトを見ながらやってみてください。

まずはMinIOクライアント（mc）をインストールしていきましょう。

[https://min.io/docs/minio/linux/reference/minio-mc.html#minio-client](https://min.io/docs/minio/linux/reference/minio-mc.html#minio-client)

こちらのサイトをもとにやっていきます。

インテルなら、
```bash
curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  --create-dirs \
  -o $HOME/minio-binaries/mc

chmod +x $HOME/minio-binaries/mc
```
これでインストール＋権限付与が完了します。

次にパスを通しましょう。

```bash
export PATH=$PATH:$HOME/minio-binaries/
```
を~/.bashrcなどに書き込んで、再度読み込みをしましょう。

次に、サーバーに司令を送るための準備（ログイン？）を行います

ログインキーが履歴に残ってしまわないように、ちょっと工夫しつつ行います

```bash
bash +o history
mc alias set myminio http://localhost:9000 [あなたが設定したrootユーザー] [あなたが設定したパスワード]
bash -o history
```

これで、サーバーに対する操作ができるようになりました。

次にバケットという、それぞれのオブジェクトを一括管理する単位を作っていきます。

[https://min.io/docs/minio/linux/reference/minio-mc.html#minio-client](https://min.io/docs/minio/linux/reference/minio-mc.html#minio-client)

ここにコマンドのリファレンスがまとまっているので、やりたいことがあるときにはここを覗くとよいでしょう。

新しいバケットを作るには、`mc mb`を使います。

```bash
mc mb myminio/infra-hands-on
```

これで`myminio`の中に`infra-hands-on`というバケットができました。

次に、公開設定をします。

```bash
mc anonymous set public myminio/infra-hands-on
```
これで、認証などをなしに画像にアクセスできるようになりました。

## アプリに導入しましょう（Nさんが一晩でやってくれました）
では、現在ローカルストレージに保存するようにしている仕組みをMinIOに切り替えていきましょう。

今回も、アプリケーション側を触る必要はありません。

環境変数に必要な情報を入れるだけで実装が切り替わるようになっています。

```bash
ICON_STORE_ENDPOINT # エンドポイント。今回はhttp://localhost:9000
ICON_STORE_BUCKET # バケット　さっき作成したもの
ICON_STORE_ACCESS_KEY # ルートユーザーにします
ICON_STORE_SECRET_KEY # パスワードにします
ICON_STORE_BASE_URL #MinIOの場合は<エンドポイントURL>/<バケット> とします
ICON_STORE_PREFIX # 好きなプレフィックスを設定しましょう。
```
この環境変数をGoアプリケーションのサービスファイルに書き込みましょう。(root権限が必要です)

できたら再起動をします。
```bash
sudo systemctl daemon-reload # ファイルの中身を変えたらこれで読み込ませる
sudo systemctl restart <サービス名>
```

参考は03の教材

ここまでできたら、動作確認をしてみましょう。

WSLの場合、9000番をsshトンネルでポートフォワーディングしておく必要があります。

ちゃんとアイコン設定できましたか？

出来たら成功です！！

お疲れ様でした