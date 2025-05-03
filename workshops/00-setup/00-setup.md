## はじめに

このリポジトリは、簡易的な技術で作られたチャットアプリを強化しながら、いろいろなサービスのセットアップを学ぶハンズオンです。

VMやAWSなど、隔離されたサーバーを確保して行っていただければと思います。

以下に、Ubuntu上で`qemu-kvm`を用いてVMを立てる手順を軽く紹介しておきます。

## qemu-kvmでUbuntuディストリビューションのVMを立てる

まずは必要なものをinstallしましょう。
```bash
sudo apt update
sudo apt install virt-manager
sudo apt install qemu-kvm
```
このあたりがあればOK。

次に、Ubuntuサーバーのisoをゲットしましょう。

[https://ubuntu.com/download/server](https://ubuntu.com/download/server)

ダウンロード後、以下のようにisoのファイルのパーミッションを調整します
```bash
sudo chmod 644 ubuntu-24.04-live-server-amd64.iso
```
また、VMがKVMにアクセスできるように以下の設定も必要な場合があります。
```bash
sudo chmod a+rw /dev/kvm
```
### virt-manager立ち上げ
```bash
virt-manager
```
上記で仮想マシンマネージャ（GUI）が起動します。

### VM立てる
新規作成を押下して、ローカル参照で先程入れたisoを選択しましょう。

isoが選択できればあとはひたすら続けるだけ。

メモリ割り当てなどはアプリケーションが動くくらいのものにしてください。

Ubuntuセットアップ時にメモリの割り当てで空白（Free）を生まないように編集するのを忘れないようにしましょう。

## sshできるように準備する
基本的にsshでリモートログインして操作したほうが便利なので、**openssh-server**を入れておきましょう
```bash
sudo apt update
sudo apt install -y openssh-server
```
入れ終わったら、一応
```bash
sudo systemctl status ssh
```
で状態確認をしましょう。
