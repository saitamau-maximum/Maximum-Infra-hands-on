## この節でやること
- nginx のログ出力について知る
- 出力形式を変更してみる
- Webサーバーのログ解析をやってみる

## 適切な運用・改善にはログ解析が必須です。
あなたの運用している「インフラチャット」の社内普及率もかなり上がってきて、そろそろ本格的なサーバー運用が必要になってきました。

かつては時たまアクセス集中でサーバーが落ちる程度でしたが、最近では頻繁にサーバーが落ちたり、異常に重くなるようになってしまいました。

そこで、あなたとNさんはサーバーのログを監視して、どのようなシナリオが負荷をかけているのかを発見しようと考えます。

## Nginxは何のログを取る？
この節では、Webサーバーの一つであるNginxでのログの取り扱いについて学びます。

Webアプリケーションの高速化や、動作不良の原因調査、セキュリティ担保においてログ収集・解析は非常に強い手段です。

Nginxでは、デフォルトでアクセスログ出力、エラーログ出力が有効になっています。

`/var/log/nginx/`の中にログファイルが入っています。

試しにアクセスログを覗いてみましょう。

```bash
less /var/log/nginx/access.log
```
何かしらの文字列が並んでいますね。どうやらuriっぽいものも散見されます。

これがどのようなフォーマットで並んでいるのかは、Nginxの基本設定ファイルを覗くと分かります。

```bash
less /etc/nginx/nginx.conf
```
で基本設定ファイルを覗いてみましょう。

一部抜粋しますが、以下のようなものがみれると思います
```bash
http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;
```
log_format の部分でログの吐き方を定義しています。

access_log の部分でログの場所を指定しています

これで、アクセスの履歴を確認することができるというわけですね。

## 結局jsonが扱い易いってワケ
今後の展望として、ログファイルの一元管理、解析があります。

そのためには、何となく文字列で並んでいるのではなくて、明確に「ラベルー値」のセットになっていて欲しいと思いませんか？

そこで、アクセスログの出力形式をjson形式に変更してみましょう。

まず、Nginxの基本設定ファイルを編集して、新しいフォーマット形式を定義しましょう。

何かしらの編集手段（nanoやvim）で`/etc/nginx/nginx.conf`を開いてください。

先ほど確認した、mainのformatのすぐ下に、以下を追記して、新しいフォーマット形式を定義します。

```conf
    log_format json escape=json '{'
        '"time":"$time_iso8601",'
        '"host":"$remote_addr",'
        '"port":"$remote_port",'
        '"method":"$request_method",'
        '"uri":"$request_uri",'
        '"status":"$status",'
        '"body_bytes":"$body_bytes_sent",'
        '"referer":"$http_referer",'
        '"ua":"$http_user_agent",'
        '"request_time":"$request_time",'
        '"response_time":"$upstream_response_time"'
    '}';
```
これで、`json`という名前の新しいログ整形形式ができました。

次に、この形式を使ってアプリケーションのログを出力していきます。

アプリケーション単位でログを出したいので、`/etc/nginx/conf.d/InfraHandsOn.conf`に出力設定を追記します。

```conf
# jsonフォーマットでアクセスログを出力
access_log /var/log/nginx/infra_hands_on_access.log json;
```

これを、第04章で追加した`/etc/nginx/cond.d/InfraHandsOn.conf`に追記してください。

場所はどこに書いても適用されますが、慣例に倣って、`server_name`の一つ下に書き込みましょう。

出来たら、一応構文テストをかけて
```bash
sudo nginx -t
```
よさそうな設定を反映するために再起動しましょう
```bash
sudo systemctl restart nginx
```

これによって、このアプリケーションに向けたアクセスが、json形式で`infra_hands_on_access.log`というファイルに出力されるようになります。


実際にアプリケーションに訪れてからログを覗いてみてください。

```bash
less /var/log/nginx/infra_hands_on_access.log
```

記録がいくつも残っているかと思います。json形式であることも確認できると思います。

これで、自分たちで定義したフォーマットを使ってログを出力させることができました。

## アクセスログを解析してみる
ログを見れるようになりましたが、これを目視で確認して

「このURIへのアクセスが多いかも！！」

「このアクセスに対する返答が重いかも！！」

とするのは手間がかかって大変です。

これをなすために、様々な便利ツールが発明されていますがここでは`alp`という解析ツールを紹介します。

### alpのインストール
基本的には、バイナリを実行してインストールします。

まずはバイナリをダウンロードしましょう。

コードは<https://github.com/tkuchiki/alp/releases>にあります。AssetsのLinux用のものをwgetとかで持ってきましょう。

```bash
wget https://github.com/tkuchiki/alp/releases/download/v1.0.21/alp_linux_amd64.tar.gz
```
持ってこれたら、解凍します。
```bash
tar zxvf alp_linux_amd64.tar.gz
```
最後に、インストールをしましょう
```bash
sudo install alp /usr/local/bin
```
正常にインストールできたか確認。
```bash
alp --help
```
それらしきものが出たら成功です。

### alpを使う
json形式のログを解析するには、`alp json`を使います。

渡し方には二つあって、
- パイプを使って渡す
```bash
cat /var/log/nginx/infra_hands_on_access.log | alp json
```
- オプションを使う（`--file`）
```bash
alp json --file /var/log/nginx/infra_hands_on_access.log
```
いずれのやり方でも、標準出力に結果が出てきたと思います。

もし記録を残したい場合は` > log.txt`みたいにして、結果をメモしておくとよいでしょう。

### いろいろな切り口で眺める
ただ眺めるだけではなく、例えば時間がかかっている順に表示したり、特定のアクセスをグループとしてみたりできます。

例えば、
```bash
alp json --file /var/log/nginx/infra_hands_on_access.log --sort sum -r
```
のように、`--sort sum -r`オプションを追加すると処理合計時間（SUM）のソートをして表示してくれます。

また、正規表現によってURIをまとめることができます。

このハンズオンのアプリでは、`/api/message/[部屋のID]`によって部屋ごとのメッセージを取得します。

解析時には、俯瞰して「部屋のメッセージ取得には全体に対してどのくらいの時間がかかっているか」を見たいときがあります。

実は、alpは部屋IDの情報はつぶして、`/api/message/`までのくくりで表示するという事が出来ます。（正規表現によって）

先述の例で言うなら、
```bash
alp json --file /var/log/nginx/infra_hands_on_access.log --sort sum -r -m "/api/message/*"
```
のように、`-m`オプションの後に正規表現を書くと、それにマッチするURIをまとめてくれます。

## まとめ
Nginxのログを覗いて、整形方法を変えてから、alpでの解析までやりました。

お疲れ様でした。

