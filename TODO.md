# 課題管理 (TODO)

## 🚨 直近の修正 (High Priority)
- [ ] `internal/ui/directory_tree_model.go` のコンパイルエラー修正
    - `walk.TreeItem` インターフェースの実装 (`Child` -> `ChildAt`)
    - `walk.TreeModel` インターフェースの実装 (`Root` -> `RootAt`)
- [ ] `internal/ui/mainwindow.go` の修正
    - `TreeView` へのモデル適用部分の型整合性確認

## 🚀 機能実装 (Features)
- [ ] TreeViewでのディレクトリ選択とTableViewの連動
- [ ] TableViewのカラム追加（更新日時、ファイル種別）
- [ ] Grep検索機能 (EX-grep)
- [ ] ディレクトリ容量可視化 (Visual Sum)

## 🔧 改善 (Improvements)
- [ ] エラーハンドリングの強化
- [ ] テストコードの拡充