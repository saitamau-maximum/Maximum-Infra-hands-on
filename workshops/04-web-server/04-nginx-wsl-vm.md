## この章でやること
- Webサーバーを知る
- フロントエンドをビルドして配信する
- バックエンドへのリクエストのやり方を変える

## devというだけあって、本番には向かない。
今日は華の金曜日。あなたはNさんと一緒にサイゼリヤで昼食をとっています。

ナイフとフォークで辛味チキンと格闘しているNさんに、以前

> 上司「妙にこのアプリは重いな。修正してくれないか。」

と言われたことについて相談してみました。

> N「このアプリ、そんなに大したことはしていないはずなのにね。今ってどういう構成で動いているんだっけ……」

あなたはNさんに現在の構成を話します。

> N「開発環境をそのまま配信しているんだ……。
> 
> フロントエンドのことはよくわからないけれど、devサーバーはパフォーマンスが良くなかったり、脆弱だったりするという話を聞いたことがあるよ。
> 
> ちょうど明日から休みだし、一緒に本番用の配信方法に移行してみない？　」

---

`npm run dev`で立ち上がる開発用サーバーは、その名の通り開発用で、本番環境には向いていません。

開発用サーバーでは、開発者が使いやすいように様々な工夫がなされています。

例えば、'jsx'などの形式は本来ブラウザでは動きませんが、開発サーバーはそれを **トランスパイル** することによって上手に表示しています。

また、ホットリロード機能（変更を加えるとすぐに反映される）などもあります。

しかし、これらの機能が備わっているゆえに本番環境には向きません。

脆弱性を突かれてコードが漏洩したり、パフォーマンスが悪くなったりしてしまうのです。

## npm run buildはここで使う
devサーバーを公開しないなら、何を公開すればいいのでしょうか。

答えは、buildファイルです。

試しに作ってみましょう。

`~/Maximum-Infra-hands-on/frontend`で、
```bash
npm run build
```
をしてみてください。`frontend`直下に`dist`というディレクトリが生成されていると思います。

これを公開してみましょう。

## 便利なWebサーバーを連れてきたよ。
静的ファイルの配信において便利なのが、Webサーバーです。

`Apache`や`Nginx`などがよく聞かれますが、今回は`Nginx`を使って配信をしてみましょう。

### Nginxをインストール

まずはインストールから。公式サイトの記述に従っていきます。（Ubuntu想定）

[https://nginx.org/en/linux_packages.html#Ubuntu](https://nginx.org/en/linux_packages.html#Ubuntu)

まずは、操作に必要なパッケージをまとめてインストールします。

```bash
sudo apt install curl gnupg2 ca-certificates lsb-release ubuntu-keyring
```

次に、署名鍵を取得します。これによって`sudo apt update`を実行したときに検証がうまく行きます。

```bash
curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor | sudo tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null
```
特にエラーが出なければ、確認をしてみてください。
```bash
gpg --dry-run --quiet --import --import-options import-show /usr/share/keyrings/nginx-archive-keyring.gpg
```
573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62　があればOK

次に、安定版をインストールするためにリポジトリを追加します。

```bash
echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] \
http://nginx.org/packages/ubuntu `lsb_release -cs` nginx" \
    | sudo tee /etc/apt/sources.list.d/nginx.list
```
これをこのまま実行しましょう

できたら、最後にinstall
```bash
sudo apt update
sudo apt install nginx
```

```bash
nginx -v
```
で、バージョンが出力されたらOKです。

ついでに有効化しておきます。
```bash
sudo systemctl start nginx
```

## Nginxを使うと、どこが改善できる？
今の公開の仕方の問題点をまとめておきます。
1. フロントエンドがdevサーバーで配信されている。
2. バックエンドの8080が公開されている。

この２つの問題をNginxを使って改善していきましょう。

### 準備を整えましょう。
まずは、nginxに必要なポートを開けて、不必要なポートを閉じましょう。
```bash 
sudo ufw allow 80,443/tcp
sudo ufw delete allow 5173
sudo ufw delete allow 8080
```

前者はnginxに必要なポートを開けています。後者はいままでサーバーを立てていたポートを閉じています。（denyなどに設定してもよいです）

Nginxにやってもらうことは以下の２つです。
1. 80番ポートで、フロントエンドの静的ファイル（buildしたやつ）を配信する
2. [サーバーのIP]/api...のようなリクエストをバックエンド向けと判断して、'http://localhost:8080`に転送する

Nginxの設定ファイルを作ることでこれらは実現できます。

0から作るには、少しこのアプリについての造形を深める必要があり、難しいかもしれません。

と、いうことで今回はNさんが作成してくれた`InfraHandsOn.conf`ファイルを雛形にして、設定を進めていきましょう。

ファイルの前半がフロントエンド用、後半がバックエンド用になっています。

`server_name`はいつものように埋めてください。

`root`を埋めるために、ファイルの置き場所を作っていきます。

```bash
sudo mkdir -p /var/www/InfraHandsOn
```
`-p`は、親ディレクトリまで含めて一括作成するためのオプションです。

では早速build……の前に。先程あげた、Nginxにやってもらうことその２をするために、`.env`ファイルの中で、IPを指定している部分を削除してください

例
```bash
VITE_API_BASE_URL=http://192.168.123.8:8080
↓
VITE_API_BASE_URL=""
```

変更を反映するために、ビルドをします。
```bash
npm run build
```
できたら、その成果物を先程作った置き場所にコピーします。

```bash
sudo cp -r ~/Maximum-Infra-hands-on/frontend/dist/* /var/www/InfraHandsOn/
```

次に、さっき作った設定ファイルをコピペするか、コピーするかして、
`/etc/nginx/conf.d/InfraHandsOn.conf`というような、設定ファイルを配置してください。

できたら、中の必要な部分を埋めていきましょう。

`root`の部分は、例のままで大丈夫です。

`server_name`の部分は`localhost`にします。これはポートフォワーディングによってリクエストurlがlocalhostになってしまうからです。

外部サーバーなどにデプロイするときは、この部分を自身のIPアドレスにします。


おけたら、
```bash
sudo nginx -t
```
でテスト。

うまくいったら、設定を反映させるために再起動しましょう。

```bash
sudo systemctl restart nginx
```
再起動できたら、あなたのサーバーにアクセスしてみてください。

ポートフォワーディングは変える必要があります。

Nginxは80番のポートでリッスンしているので、VMの80番を8080にフォワーディングするために、

```bash
ssh -L 8080:localhost:80 [いつもsshするときにやってるやつ]
```

をしましょう。そして、

`http://localhost:8080`にアクセスすると……？

フロントエンドが見れましたか？

見れたなら、公開の設定は完了です！！！

お疲れ様でした。

## おまけ〜ログ収集について〜
Nginxと、拡張モジュールを使うと、この章で行ったこと以外にもいろいろなことができます。
- アクセスログ収集
- レスポンスのgzip圧縮
- request/responseのバッファリング
- TLS通信の高速化
- ロードバランス
- コンテンツキャッシュ
- URLルーティング・書き換え
