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