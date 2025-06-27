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

`Apache`や`Nginx`などがよく聞かれますが、今回は`Apache`を使って配信をしてみましょう。

### Apacheをインストール

まずはインストールから。公式サイトの記述に従っていきます。（Ubuntu想定）

[https://ubuntu.com/tutorials/install-and-configure-apache#1-overview](https://ubuntu.com/tutorials/install-and-configure-apache#1-overview)

```bash
sudo apt update
sudo apt install apache2
```

インストール出来たら、起動・自動起動を設定します。

```bash
sudo systemctl start apache2
sudo systemctl enable apache2
```

正常にできたかどうか確認します。

```bash
sudo systemctl status apache2
```
結果として、
```bash
● apache2.service - The Apache HTTP Server
     Loaded: loaded (/usr/lib/systemd/system/apache2.service; enabled; preset: enabled)
     Active: active (running) since Fri 2025-06-27 03:37:20 UTC; 4min 20s ago
       Docs: https://httpd.apache.org/docs/2.4/
   Main PID: 2801 (apache2)
      Tasks: 55 (limit: 4609)
     Memory: 5.4M (peak: 5.9M)
        CPU: 210ms
     CGroup: /system.slice/apache2.service
             ├─2801 /usr/sbin/apache2 -k start
             ├─2802 /usr/sbin/apache2 -k start
             └─2803 /usr/sbin/apache2 -k start
```
このようなログが出たら成功しています。

ポートフォワーディングでVMの80番ポートを自分の8080番につなげたうえで、
`localhost:8080`を見ると、デフォルトページがみれます。

## 【寄り道】ディストリビューションによってApacheは少し異なる
Ubuntuのような Debian 系では`apache2`という名前ですが、Red Hat 系では`httpd`という名前です。

どちらも同じApacheですが、設定ファイルの配置などが微妙に異なっています。

また、ファイル配信の手順も微妙に異なっているので、これからの教材の内容はDebian系でないディストリビューションでは使えないかもしれないことを頭の隅に置いておいてください。

（公式ドキュメントを見ましょう、という平々凡々な結論に落ち着きます。）

## Apacheを使うと、どこが改善できる？
今の公開の仕方の問題点をまとめておきます。
1. フロントエンドがdevサーバーで配信されている。
2. バックエンドの8080が公開されている。

この２つの問題をApacheを使って改善していきましょう。

### 準備を整えましょう。
まずは、Apacheに必要なポートを開けて、不必要なポートを閉じましょう。
```bash 
sudo ufw allow 80,443/tcp
sudo ufw delete allow 5173
sudo ufw delete allow 8080
```

前者はApacheに必要なポートを開けています。後者はいままでサーバーを立てていたポートを閉じています。（denyなどに設定してもよいです）

Apacheにやってもらうことは以下の２つです。
1. 80番ポートで、フロントエンドの静的ファイル（buildしたやつ）を配信する
2. [サーバーのIP]/api...のようなリクエストをバックエンド向けと判断して、'http://localhost:8080`に転送する

Apacheの設定ファイルを作ることでこれらは実現できます。

今回はNさんが作成してくれた`infra-hands-on.conf`ファイルを雛形にして、設定を進めていきましょう。

ファイルを書き込む前に、`DocumentRoot`などを埋めるために、ビルドファイルの置き場所を作っていきます。

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

次に、雛型の設定ファイルをコピペするか、コピーするかして、
`/etc/apache2/sites-available/infra-hands-on.conf`というような、設定ファイルを配置してください。

できたら、中の必要な部分を埋めていきましょう。

`ServerName`の部分は`localhost`にします。これはポートフォワーディングによってリクエストurlがlocalhostになってしまうからです。

外部サーバーなどにデプロイするときは、この部分を自身のIPアドレスにします。

他の部分は例示に従ってください。

次に、これを適用します。（この適用手順は`httpd`の場合いらないです）

```bash
sudo apache2ctl configtest
```
で、構文チェックをしてから有効化をします。

おそらく、「このコマンドのプラグインつけてないよ」と言うような起こられ方をすると思うので、
```bash
sudo a2enmod rewrite # 行き先変更（SPAのデフォルト設定用）
sudo a2enmod proxy # プロキシ用（基盤）
sudo a2enmod proxy_http # プロキシ用
sudo a2enmod proxy_wstunnel # websocektをあつかうため
sudo a2enmod headers # ヘッダ引継ぎ用
sudo a2enmod remoteip # いつかRemoteIP使うとき用
```
済んだら、設定を有効化しましょう。
```bash
sudo a2ensite infra-hands-on.conf
```
出来たら、再起動もしておきましょう
```bash
sudo systemctl restart apache2
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
