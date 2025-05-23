## この章でやること
- ストレージサービスについて知る
- MinIOをセットアップしてみる

## アイコンがないと、チャットがつまらん
開発をしているうちに、Nさんは「インフラチャット」の虜になってしまいました。

「会社用アカウント」「趣味用アカウント」「愚痴用アカウント」の３つを使ってチャットを最大限楽しんでいましたが…

> N「やっぱり、アイコンがないとどのアカウントでログインしているか分かりづらいよ。社用アカウントで愚痴を吐かないためにも、早急に実装の必要があるね」

あなたはNさんの要望に答えるために、ユーザーが好きなアイコンを設定できる構成を考えなくてはならなくなりました。

## データベースに入れる？　ご乱心ですか
あなたは最初に、MySQLの中に画像データを入れることができないか、と考えました。

Nさんに早速提案しましたが、どうも反応が芳しくありません。

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

## MinIO導入
<<<<<<< HEAD
まずはMinIOのサーバー（単一ノードのやり方）を導入しましょう。

[https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd](https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd)

このサイトに従って行います。
=======
[https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd](https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html#minio-snsd)

このサイトに従って行います。

>>>>>>> 9e98878ce1894dcfc84c96f60ab347a42599cebd
