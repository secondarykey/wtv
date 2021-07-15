# wtv

wtv is webtoon viewer

## 

基本的にはebitenの調査レベルで行っていたのですが、ニーズがあるかな？と思って使えるところまで実装しました。

コンテンツ等の読み込みの負荷処理、GUI部品等の共通化を行う為、
開発は継続する予定でいます。

## 最適化

現在は画像サイズにより自動的に最適化を行います。

## Issue

- MenuのComponent化

  PlayerをSceneにする

- 右寄せ、下寄せのレイアウト設計

  Addかな、、、

- もう少し表示領域分しか利用しない仕組み

  - book の部分

- スクロールモードの座標を厳密にする

- 下メニュー
   - 一番上への実装（下メニュー

- デバッグモード

  - 領域に線を描画

- ボタンの実装

  トグル -> ソート のアクティブなものを活性状態にする
  ソートの上下設定

