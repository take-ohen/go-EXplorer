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

## 2026-03-08

### ✅ 実施タスク

- [x] `scanner.go` のエラーハンドリングを改善
- [x] `scanner.go` のテストコード (`scanner_test.go`) を作成

### 📝 メモ

- `Scan` メソッドが `(<-chan FileItem, <-chan error)` を返すようにシグネチャを変更。
- `os.ReadDir` や `entry.Info` で発生したエラーを、黙って無視するのではなく新設したエラーチャネルに送信するようにした。
- `testing` パッケージと `t.TempDir()` を利用して、一時的なファイル構造を生成し、正常系のテスト (`TestScanner_Scan`) を実装。
- 存在しないパスを指定し、エラーチャネルからエラーが正しく報告されることを確認する異常系のテスト (`TestScanner_Scan_Error`) も実装。

## 2026-03-09

### ✅ 実施タスク

- [x] UIをエクスプローラ風レイアウトに変更 (`HSplitter`による左右分割)
- [x] `TextEdit`でのログ表示を廃止し、`TableView`によるリスト表示へ移行
- [x] `TreeView`を導入し、ディレクトリツリー表示の基盤を実装
- [x] `DirectoryTreeModel`を新規作成
- [x] `walk.TreeItem`, `walk.TreeModel`インターフェースのメソッド名不整合によるコンパイルエラーを修正
- [x] 課題管理を `TODO.md` から `TASK.csv` へ移行
- [x] ドライブ列挙処理を `os.Stat` ループから Windows API (`GetLogicalDrives`) へ変更し、起動時の遅延を解消
- [x] パフォーマンス調査のため、`ListDir` およびUI処理に計測ログを追加

### 📝 メモ

- 課題管理方法を`TASK.csv`での手動管理に、開発ログを`DEV_LOG.md`での記録に、それぞれ再定義。
- `logger`ツールおよび関連する旧管理ファイル(`TODO.md`, `task.md`, `DEV_LOG.csv`)は廃止。
- `TODO.md` に記載されていた未着手タスクを `TASK.csv` に転記し、同ファイルを削除。