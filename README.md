# Maximum-Infra-hands-on

**教材は`/workshops`にまとめてあります**

**遊ぶ際は、mainブランチの内容で遊んでください。**

インフラ班で扱うコンテンツにしたい。

小規模なISUCONみたいな感じ。

## アイデア
- サーバーにこのリポジトリを置いて、MySQLやRedis、Nginxなどの設定や、FireWallの設定をする
- バックエンドが最初はsqliteやインメモリ（map）で実装されていて、それぞれの弱みを補うためにセッティングを進めていくシナリオ
- チャットアプリを想定。WebSocketで実装している

## 開発を手伝ってくれる方へ
```bash
cd frontend
npm install
```
で環境構築をお願いします。

プロジェクトのルートで
```
bash ./scripts/start-dev
```
で、開発環境が動きます。
