# back-challe-chara2022

## Qmatta とは

![image](https://user-images.githubusercontent.com/80093134/194593649-e5a5ee43-77cf-4ab3-a460-c09d98c07fa5.png)

ベアプログラミングをアプリ上で実現するサービスです．

### ベアプログラミングとは
解決できない問題に遭遇した時，誰かにその問題について話すことで思考が整理され，解決策を思いつく現象のことです．</br>
テディベア効果やラバーダック・デバッグといった名称もあります．</br>


## dockerの起動

### コンテナのビルド方法
- 初めてサービスを立ち上げる時に利用（imageがない時も同様）
- requirements.txtを更新した際は再度ビルドする必要がある
```bash
 $ make build-dev
```
### コンテナの作成，起動
- 基本的にはupを利用する
- docker-composeにサービスを追加，requirements.txtを更新した時は再ビルドする必要がある
```bash
 $ make up-dev // コンテナの作成、起動
```

## MinIOの利用方法
- [MinIO](https://min.io)
- MinIO: ローカル環境におけるAWS S3環境

### MinIOコンソールにアクセス
コンテナ起動後，[http://localhost:9090]にアクセスするとMinIOのコンソールが起動
- ログイン時のIDとパスワードは`.env`で指定
![スクリーンショット 2022-10-14 17 11 12](https://user-images.githubusercontent.com/29566903/195796672-66d43868-ea27-475a-9fed-eb310a1f5cb6.png)
![スクリーンショット 2022-10-14 17 11 29](https://user-images.githubusercontent.com/29566903/195796771-650a891d-e1a0-4f2f-85e1-44daf0f473ed.png)

## 開発ルール

1. issueを確認する or issueを作成する
2. ブランチを作成する（`[タグ]/#[issue番号]-[実装内容]`）

| タグ | 詳細 |
| --- | --- |
| feat | UI・機能実装 |
| bugfix | バグ修正 |
| setup | セットアップ |
| release | リリース |
| develop | ローカルでの実行環境 |
