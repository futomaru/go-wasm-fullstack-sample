package main

// //go:wasmexport のすぐ下に「エクスポートしたい関数」を書く
// エクスポート名は "Sub"（ホスト側からこの名前で呼び出す）

//go:wasmexport Sub
func Sub(a int32, b int32) int32 {
	return a - b
}
