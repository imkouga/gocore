package xlsx

import (
	"testing"
)

// 新增一行
func TestAddRow(t *testing.T) {
	f := NewFile()
	s, _ := f.AddSheet("test")
	s.AddRowByCells([]interface{}{"one", "two"}...)
}
