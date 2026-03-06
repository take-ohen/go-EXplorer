package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	// 1. テスト用の一時ディレクトリ構造を作成
	// root/
	//   file1.txt
	//   subdir/
	//     file2.txt
	rootDir := t.TempDir()

	subDir := filepath.Join(rootDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	file1 := filepath.Join(rootDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. スキャナーの実行
	scanner := NewScanner()
	out, errc := scanner.Scan(rootDir)

	// 3. 結果の検証
	// 期待されるパスのマップ (見つかったら true にする)
	expectedFiles := map[string]bool{
		file1:  false,
		file2:  false,
		subDir: false, // ディレクトリ自体もFileItemとして報告される仕様
	}

	// ファイル情報チャネルの受信
	for item := range out {
		if _, ok := expectedFiles[item.Path]; ok {
			expectedFiles[item.Path] = true
			t.Logf("Found: %s", item.Path)
		} else {
			t.Errorf("Unexpected item found: %s", item.Path)
		}
	}

	// エラーチャネルの確認 (正常系なのでエラーは期待しない)
	// チャネルがcloseされるまで読み切る
	for err := range errc {
		if err != nil {
			t.Errorf("Unexpected error received: %v", err)
		}
	}

	// 全ての期待するファイルが見つかったか確認
	for path, found := range expectedFiles {
		if !found {
			t.Errorf("Expected item not found: %s", path)
		}
	}
}

func TestScanner_Scan_Error(t *testing.T) {
	scanner := NewScanner()
	// 存在しないディレクトリを指定
	out, errc := scanner.Scan("path/to/non/existent")

	// アイテムは来ないはず
	for range out {
		t.Error("Expected no items, but got one")
	}

	// エラーが来るはず
	err := <-errc
	if err == nil {
		t.Error("Expected error, but got nil or channel closed without error")
	}
	t.Logf("Received expected error: %v", err)
}
