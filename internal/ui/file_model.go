package ui

import (
	"github.com/lxn/walk"
	"github.com/take-ohen/go-EXplore/internal/core"
)

// FileInfoModel is a data model for the TableView.
type FileInfoModel struct {
	walk.TableModelBase
	items []*core.FileItem
}

// NewFileInfoModel creates a new FileInfoModel.
func NewFileInfoModel() *FileInfoModel {
	return &FileInfoModel{}
}

// RowCount returns the number of rows in the model.
func (m *FileInfoModel) RowCount() int {
	return len(m.items)
}

// Value returns the value for the specified row and column.
func (m *FileInfoModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Path
	case 1:
		return item.Size
	}
	return nil
}

// AddItems appends new items to the model and notifies the view.
func (m *FileInfoModel) AddItems(items []core.FileItem) {
	for _, item := range items {
		// Create a copy to safely take the address
		c := item
		m.items = append(m.items, &c)
	}
	m.PublishRowsReset()
}

// Reset clears the model.
func (m *FileInfoModel) Reset() {
	m.items = nil
	m.PublishRowsReset()
}
