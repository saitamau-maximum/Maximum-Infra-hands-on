## この節でやること

* Apache のログ出力について知る
* 出力形式を変更してみる
* Webサーバーのログ解析をやってみる

## 適切な運用・改善にはログ解析が必須です。

あなたの運用している「インフラチャット」の社内普及率もかなり上がってきて、そろそろ本格的なサーバー運用が必要になってきました。

かつては時たまアクセス集中でサーバーが落ちる程度でしたが、最近では頻繁にサーバーが落ちたり、異常に重くなるようになってしまいました。

そこで、あなたとNさんはサーバーのログを監視して、どのようなシナリオが負荷をかけているのかを発見しようと考えます。

## Apacheは何のログを取る？

この節では、Webサーバーの一つであるApacheでのログの取り扱いについて学びます。

Webアプリケーションの高速化や、動作不良の原因調査、セキュリティ担保においてログ収集・解析は非常に強い手段です。

Apacheでは、デフォルトでアクセスログ出力、エラーログ出力が有効になっています。

`/var/log/apache2/`の中にログファイルが入っています。

試しにアクセスログを覗いてみましょう。

```bash
less /var/log/apache2/access.log
```

何かしらの文字列が並んでいますね。どうやらuriっぽいものも散見されます。

これがどのようなフォーマットで並んでいるのかは、Apacheの設定ファイルを覗くと分かります。

```bash
less /etc/apache2/apache2.conf
```

一部抜粋しますが、以下のようなものがみれると思います。

```apache
LogFormat "%h %l %u %t \"%r\" %>s %b" combined
```

LogFormatでログの出力形式を定義することができます。

次に、これを使って実際にログを出力している部分を見てみましょう。

```bash
less /etc/apache2/sites-available/000-default.conf
```
CustomLog の部分でログの場所とフォーマットを指定しています。

先ほど見た combinedというフォーマット形式（自分で定義）で出力しています。

これで、アクセスの履歴を確認することができるというわけですね。

## 結局jsonが扱い易いってワケ

今後の展望として、ログファイルの一元管理、解析があります。

そのためには、何となく文字列で並んでいるのではなくて、明確に「ラベルー値」のセットになっていて欲しいと思いませんか？

そこで、アクセスログの出力形式をjson形式に変更してみましょう。

まず、Apacheの設定ファイルを編集して、新しいフォーマット形式を定義しましょう。

何かしらの編集手段（nanoやvim）で`/etc/apache2/apache2.conf`を開いてください。

次のような行を追加します（LogFormat 定義の近くが良いです）：

```apache
LogFormat "{ \"time\":\"%{%Y-%m-%dT%H:%M:%S%z}t\", \"host\":\"%h\", \"method\":\"%m\", \"uri\":\"%U\", \"query\":\"%q\", \"status\":\"%>s\", \"ua\":\"%{User-Agent}i\", \"referer\":\"%{Referer}i\", \"bytes\":\"%B\", \"reqtime\":\"%D\" }" json
```

これで、`json`という名前の新しいログ整形形式ができました。

次に、ハンズオンのための設定ファイルにアクセスログをこの形式で出力する設定をしましょう。

お好きな手段で`/etc/apache2/sites-available/infra-hands-on.conf`を開いてください。これは第04章で作成した設定ファイルです。

ここに、
```apache
CustomLog /var/log/apache2/infra_hands_on_access.log json
```
を追記してください。

次に、Apacheを再起動して設定を反映しましょう。

一応構文テストをします。
```bash
sudo apachectl configtest
```
よさそうなら、再起動して設定反映。
```
sudo systemctl restart apache2
```

これによって、Apache に向けたアクセスが、json形式で`infra_hands_on_access.log`というファイルに出力されるようになります。

実際にアプリケーションに訪れてからログを覗いてみてください。

```bash
less /var/log/apache2/json_access.log
```

記録がいくつも残っているかと思います。json形式であることも確認できると思います。

## アクセスログを解析してみる

ログを見れるようになりましたが、これを目視で確認して

「このURIへのアクセスが多いかも！！」

「このアクセスに対する返答が重いかも！！」

とするのは手間がかかって大変です。

これをなすために、様々な便利ツールが発明されていますがここでは`alp`という解析ツールを紹介します。

### alpのインストール

基本的には、バイナリを実行してインストールします。

まずはバイナリをダウンロードしましょう。

コードは[https://github.com/tkuchiki/alp/releases](https://github.com/tkuchiki/alp/releases)にあります。AssetsのLinux用のものをwgetとかで持ってきましょう。

```bash
wget https://github.com/tkuchiki/alp/releases/download/v1.0.21/alp_linux_amd64.tar.gz
```

持ってこれたら、解凍します。

```bash
tar zxvf alp_linux_amd64.tar.gz
```

最後に、インストールをしましょう。

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

* パイプを使って渡す

```bash
cat /var/log/apache2/json_access.log | alp json
```

* オプションを使う（`--file`）

```bash
alp json --file /var/log/apache2/json_access.log
```

いずれのやり方でも、標準出力に結果が出てきたと思います。

もし記録を残したい場合は` > log.txt`みたいにして、結果をメモしておくとよいでしょう。

### いろいろな切り口で眺める

ただ眺めるだけではなく、例えば時間がかかっている順に表示したり、特定のアクセスをグループとしてみたりできます。

```bash
alp json --file /var/log/apache2/json_access.log --sort sum -r
```

また、正規表現によってURIをまとめることができます。

このハンズオンのアプリでは、`/api/message/[部屋のID]`によって部屋ごとのメッセージを取得します。

解析時には、俯瞰して「部屋のメッセージ取得には全体に対してどのくらいの時間がかかっているか」を見たいときがあります。

例えば、

```bash
alp json --file /var/log/apache2/json_access.log --sort sum -r -m "/api/message/*"
```

のように、`-m`オプションの後に正規表現を書くと、それにマッチするURIをまとめてくれます。

## まとめ

Apacheのログを覗いて、整形方法を変えてから、alpでの解析までやりました。

お疲れ様でした。
