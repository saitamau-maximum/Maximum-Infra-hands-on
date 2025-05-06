## この章でやること
- MySQLのセットアップ

### 注意
04章までとリポジトリの内容が変わっています。v2.0.0以前のバージョンで遊んでいた方は、最新の内容の反映をお願いします。
```bash
git stash
```
で、変更した環境変数を一時的にどけてから、
```bash
git pull
```
で、最新の状態を反映してください、そのあと
```bash
git stash apply
```
で、一時的にどけた環境変数設定が戻ります

また、Goのアプリケーションの内容をビルドし直す必要があります。
```bash
cd ~/Maximum-Infra-hands-on/backend
go build -o InfraHandsOn ./cmd/main.go
```
## 手軽さには罠がある
Nginxを導入したことで、順調な滑り出しを見せたInfraChat。

このままスムーズに動いて……！

そんな願いもむなしく、ある日上司に不具合を告げられてしまいます。

> 上司「なんか最近、めちゃくちゃチャット機能が重いらしいんだけど。修正しといて。」

(T_T)

あなたは早速、Nさんと一緒に改善に望みます。

> N「まずは、データベースのログを見たいところだけれど…、そもそもSQLiteだからログ機能がないね。
> 
> 多分、チャットのログが肥大化しすぎてSQLiteのパフォーマンスが悪くなっちゃってるのが原因な気がする。
>
> データベースをMySQLに変更するのが良さそうかも。」

あなたはNさんと手分けしてデータベースを変更することにしました。

---

SQLiteは軽量で手軽ですが、ファイルベースのシンプルさ故に同時アクセスや、多量のデータ処理に不向きです。

SQLite以外のRDBMS（Relational DataBase Management System）としては、
- MySQL
- MariaDB
- PostgreSQL

などがあります。

この章ではMySQLのセットアップをしてみましょう。

※以下の記述は、次の２つの記事を参考にしています。自力でセットアップが可能な方はやってみてください。

[UbuntuのMySQLインストール記事](https://documentation.ubuntu.com/server/how-to/databases/install-mysql/index.html)

[MySQL公式ドキュメント](https://dev.mysql.com/doc/mysql-getting-started/en/)

## まずは入手するところから。

まずはMySQLを入手しましょう
```bash
sudo apt install mysql-server
```
入手できたら、まずは状態を確認してみましょう。
```bash
sudo systemctl status mysql
```
また、MySQLのサーバーが立ち上がっているはずなので、
```bash
sudo ss -tap | grep mysql
```
で、正しくポートをリッスンしているか確認しましょう。
```
$ sudo ss -tap | grep mysql
LISTEN 0      70                  127.0.0.1:33060                   0.0.0.0:*     users:(("mysqld",pid=999,fd=21))                                                                                               
LISTEN 0      151                 127.0.0.1:mysql                   0.0.0.0:*     users:(("mysqld",pid=999,fd=23))
```
このようにでたらOKです。

標準では`3306`番ポートをリッスンしています。

設定変更をする場合には`/etc/mysql`の中にある設定ファイルをいじることが多いです。

将来的に複数台構成にする際に、データベースを分割することなどがあると思います。

そのときにはデータベースに外部からアクセスする必要が出てきますよね。

その際には、`/etc/mysql/mysql.conf.d/mysql.cnf`にある`bind-address`を変更したりします。

標準では、ループバックのみを許可するようになっています。

## ゴールはDSNを作ること
GoからMySQLを操作するために、MySQLの情報を渡さなくてはなりません。

そのために渡すのがDSNです。以下のように構成されます。
```bash
[ユーザー名]:[パスワード]@[通信方法]([IPアドレスとポート])/[データベース名]?[オプション]
```
すでに埋められるところを埋めると、
```bash
root:[パスワード]@tcp(localhost:3306)/[データベース名]?charset=utf8&parseTime=true"
```
となります。

後ろのオプションが気になる方は、ぜひ調べてみてください。

つまり、これかわあなたがNさんと協力して行うのは、
- **MySQLにログインしてパスワードを設定**
- **データベースを新しく立ち上げる**

の２つです。

順にこなしていきましょう。

### ログイン
```bash
sudo mysql -u root -p
```
で、rootユーザーでログインできます。

MySQLのパスワードを求められます。設定しましょう。

### データベース作成

```bash
SHOW DATABASES;
```
とコマンドプロンプトに打ち込むことで、現在あるデータベースを確認できます。
```
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
```

ログインできたら、次にデータベースを作成しましょう。

新しく`infra_hands_on`というデータベースを作ってみましょう
```bash
CREATE DATABASE infra_hands_on;
```
で、`infra_hands_on`という名前のデータベースを作ることができます。

先程のコマンドで、追加されたデータベースを確認してみてください。

### 認証方法の変更
MySQLでは、デフォルトでおそらく`auth_socket`による認証が行われています。確認してみましょう
```bash
SELECT user, host ,plugin FROM mysql.user WHERE user = 'root' AND host = 'localhost';
```
これは、`mysql`というデータベースの`user`テーブルから、`user`カラムが`root`、`host`カラムが`localhost`の行の、`user`、`host`、`plugin`の情報を取得しています。

`plugin`は認証プラグインです。これが`auth_socket`のままだと、Goのアプリケーションからアクセスできません。

ここを、`caching_sha2_password`に変えましょう。

```bash
ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY '[新しいパスワード]';
FLUSH PRIVILEGES;
```
これで、最終目標であるDSNが作れます。

もう一度載せておくと

```bash
root:[パスワード]@tcp(localhost:3306)/[データベース名]?charset=utf8&parseTime=true"
```
です。

## Nさんが一晩でやってくれました
さて、バックエンドのアプリケーションがMySQLに接続するためのMySQL側の準備が整いました。

**次に、いままではSQLiteで実装されていたGoアプリケーションを、MySQLで動くように変更しましょう！！**

でも、データベースにアクセスする操作はアプリケーションの様々な場所で行われており、そのすべてを変更するのは非常に大変です。

どうしたらいいでしょうか……

ご安心ください。Nさんが一晩でやってくれました。

なんと、**偶然にも** このGoアプリケーションは非常にクリーンな設計がなされていて、カンタンに技術の差し替えが可能な状態でした！

> Nさん「環境変数MYSQL_DSNにdsnを設定すれば、自動的にMySQLに切り替わるように実装したよ」

それでは、環境変数を設定しにいきましょう。

Goアプリケーションは常駐化しているので、`.service`ファイルを編集することで環境変数を設定できます。

```bash
sudo vim /etc/systemctl/system/InfraHandsOn.service
```
などで、編集します。編集に使うエディタは何でもいいです。

```
Environment=CORS_ORIGIN=http://192.168.123.8:5173
↓
Environment=MYSQL_DSN="あなたのDSN"
```
このように変更しましょう。

ちなみに、`CORS_ORIGIN`については、なくしてしまって大丈夫です。

なぜなら、04章の変更によって、フロントエンドからバックエンドへの通信がクロスオリジンではなくなったからです。

serviceファイルに書き込めたら、まずは変更を反映させます。
```bash
sudo systemctl daemon-reload
sudo systemctl restart InfraHandsOn.service
```

ここまでやったら、
```bash
sudo systemctl status InfraHandsOn.service
```
を見てみてください。

ログに、`MySQL connected`と出ていたら、データベースが切り替わっています。成功です！！

お疲れ様でした。
