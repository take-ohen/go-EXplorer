# 開発ログ

## 2026-03-03

### ✅ 実施タスク

- [x] プロジェクト初期設定
- [x] `.gitignore` の作成
- [x] Standard Project Layout に基づくディレクトリ構成の作成
- [x] `internal/core/scanner.go` に並行スキャン処理のプロトタイプを実装

### 📝 メモ

- GUIライブラリとして `lxn/walk` を採用決定。
- コアロジックとUIロジックを `internal` パッケージに分離する方針。
- `scanner.go` では、Goroutineとチャネルを利用して非同期にファイル情報をストリームする。
- `sync.WaitGroup` で全Goroutineの完了を待ち、チャネルを `close` する。
- セマフォ (`chan struct{}`) を使用して、同時に実行されるGoroutineの数を制限し、リソース枯渇を防ぐ。

## 2026-03-27

### ✅ 実施タスク

- [x] `scanner.go` のエラーハンドリングを改善
- [x] `scanner.go` のテストコード (`scanner_test.go`) を作成

### 📝 メモ

- `Scan` メソッドが `(<-chan FileItem, <-chan error)` を返すようにシグネチャを変更。
- `os.ReadDir` や `entry.Info` で発生したエラーを、黙って無視するのではなく新設したエラーチャネルに送信するようにした。
- `testing` パッケージと `t.TempDir()` を利用して、一時的なファイル構造を生成し、正常系のテスト (`TestScanner_Scan`) を実装。
- 存在しないパスを指定し、エラーチャネルからエラーが正しく報告されることを確認する異常系のテスト (`TestScanner_Scan_Error`) も実装。