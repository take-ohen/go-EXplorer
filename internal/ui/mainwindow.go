package ui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/take-ohen/go-EXplore/internal/core"
)

// Run initializes and starts the GUI application.
func Run() {
	var mw *walk.MainWindow
	var tv *walk.TableView
	var pathLE *walk.LineEdit
	var treeView *walk.TreeView
	fileListModel := NewFileInfoModel()
	treeModel := NewDirectoryTreeModel()

	// 起動時の自動遷移制御用フラグ
	isSpecialInit := true

	// ディレクトリ移動処理
	changeLocation := func(path string) {
		if path == "" {
			return
		}
		pathLE.SetText(path)
		mw.SetTitle("Go-EXplorer - " + path)

		// UIフリーズ回避のため非同期実行
		go func() {
			start := time.Now()
			items, err := core.ListDir(path)
			fmt.Printf("[DEBUG] UI ListDir Total: %v\n", time.Since(start))
			mw.Synchronize(func() {
				fileListModel.Reset()
				if err == nil {
					fileListModel.AddItems(items)
				} else {
					// エラー時はリストを空にする等の処理（今回は簡易実装）
					// Access denied or other errors are ignored for list view performance
				}
			})
		}()
	}

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Go-EXplorer",
		MinSize:  Size{Width: 800, Height: 600},
		Layout:   VBox{},
		// ウィンドウ表示時に初期パスへ移動
		OnSizeChanged: func() {
			if isSpecialInit {
				isSpecialInit = false
				// 初期パスへの移動
				changeLocation("C:\\Windows")
			}
		},
		Children: []Widget{
			// 1. 上部ツールバーエリア
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{Text: "Path:"},
					LineEdit{
						AssignTo: &pathLE,
						Text:     "C:\\Windows", // デフォルト値
					},
					PushButton{
						Text: "Go",
						OnClicked: func() {
							changeLocation(pathLE.Text())
						},
					},
					PushButton{
						Text: "Deep Scan",
						OnClicked: func() {
							mw.SetTitle("Go-EXplorer - Deep Scanning... " + pathLE.Text())
							fileListModel.Reset()
							go scanPath(pathLE.Text(), mw, fileListModel)
						},
					},
				},
			},
			// 2. メインエリア（左右分割）
			HSplitter{
				Children: []Widget{
					// 左ペイン：ツリービュー
					TreeView{
						AssignTo: &treeView,
						Model:    treeModel,
						MinSize:  Size{Width: 200},
						OnCurrentItemChanged: func() {
							// 初期化中はイベントを無視
							if isSpecialInit {
								return
							}
							if item := treeView.CurrentItem(); item != nil {
								node := item.(*DirectoryNode)
								changeLocation(node.Path())
							}
						},
					},
					// 右ペイン：ファイルリスト
					TableView{
						AssignTo:         &tv,
						Model:            fileListModel,
						AlternatingRowBG: true,
						Columns: []TableViewColumn{
							{Title: "Path", Width: 400},
							{Title: "Size", Width: 100, Format: "%d"},
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}

	// ModelにMainWindowの参照を渡し、ルートノードを初期化する
	treeModel.SetMainWindow(mw)
	mw.Run()
}

// scanPath performs the file scanning in a separate goroutine.
func scanPath(root string, mw *walk.MainWindow, model *FileInfoModel) {
	scanner := core.NewScanner()
	out, errc := scanner.Scan(root)
	var wg sync.WaitGroup
	var errorCount int

	// エラー表示用のGoroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range errc {
			// エラー個別の出力は行わず、件数のみカウントする
			errorCount++
		}
	}()

	// 結果受信ループ（バッファリング付き）
	var buffer []core.FileItem
	const batchSize = 1000

	for item := range out {
		buffer = append(buffer, item)

		if len(buffer) >= batchSize {
			chunk := make([]core.FileItem, len(buffer))
			copy(chunk, buffer)
			mw.Synchronize(func() { model.AddItems(chunk) })
			buffer = buffer[:0]
		}
	}
	// 残りのバッファを出力
	if len(buffer) > 0 {
		chunk := make([]core.FileItem, len(buffer))
		copy(chunk, buffer)
		mw.Synchronize(func() { model.AddItems(chunk) })
	}

	// エラー監視の終了を待つ
	wg.Wait()

	// UIスレッドでタイトルを更新
	mw.Synchronize(func() {
		mw.SetTitle("Go-EXplorer - Done (Errors ignored)")
	})
}
