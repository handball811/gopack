# Gopack

## 概要
データを連続で追加していきデータを取得し、取得にみすが発生したときにそれを伝えることで、必要なものだけ再取得できる用にする

取得条件をしていしたり削除できるようにする

## 情報の種類
- 一時情報
  - 一回だけ送信を試み送信の失敗を無視する
- 固定送信情報
  - 一度だけ届くまで再送を繰り返すデータ
- 状態情報
  - 変化を検出して最新の状態を送信する
  - 未検知状態と検知状態がある
  - それぞれの切り替わりを通知してくれる
  - 知っている状態に乖離があることを検知するたびに送信経路に乗せてくれる

## 使い方

```
// Packの生成
pack, err := NewPack()

// データ追加用のHandlerを取得
// // 次のデータ群に一度だけ含められるデータを登録する
onetimeHandler := GenerateOnetimeHandler(pack)
err := onetimeHandler.Set(segment, 0) // データ群に含めたいビット列を渡す

// // 届くことが確認されるまで繰り返し送信されるデータを登録する
singleHandler := GenerateSingleHandler(pack)
err := singleHandler.Set(segment, 1000) // データ群に含めたいビット列を渡す

// // 状態が切り替わっていくようなデータを登録する
statefulHandler := GenerateStatefulHandler(pack)
// // cancelは通知範囲外になった時に伝えるような情報, checkerは現在状態を伝えたい対象かを返すような関数である必要がある
id, err := statefulHandler.Set(segment, cancel , 300, checker) // 初期状態を渡す
statefulHandler.Update(id, segment) // 現在の状態が切り替わった時に呼び出す
statefulHandler.Unset(id) // 状態管理が不必要になった時に呼び出す

// 今回伝えるのに使用するビット列を回収する
index := pack.Get(segment) // 0未満は失敗

// 回収したビット列の到達が確認されたときにそのことを通知する
pack.Complete(index) // この値を渡すことで、同じデータを何度も伝えずに済むようになる

// 回収したビット列が到達しなかったことを伝え再度伝える必要があれば伝えるための機能
pack.Drop(index) // 届かなかったことを通知できないといつまでも不確定なものが積み重なるので注意
```
