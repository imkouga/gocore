package xlsx

import xlsxv3 "github.com/tealeg/xlsx/v3"

type (
	File struct {
		*xlsx.File
	}

	Sheet struct {
		*xlsx.Sheet
	}
)

func NewFile() *File {
	f := &File{
		File: xlsxv3.NewFile(),
	}
	return f
}

func (f *File) AddSheet(sheetName string) (*Sheet, error) {
	s, e := f.File.AddSheet(sheetName)
	if e != nil {
		return nil, e
	}
	ss := &Sheet{
		Sheet: s,
	}
	return ss, nil
}

func (s *Sheet) AddRowByCells(cells ...interface{}) {
	row := s.AddRow()
	for _, cell := range cells {
		rowCell := row.AddCell()
		switch cell.(type) {
		case string:
			rowCell.Value = cell.(string)
		case int64:
			rowCell.SetInt64(cell.(int64))
		case int:
			rowCell.SetInt(cell.(int))
		case float64:
			rowCell.SetFloat(cell.(float64))
		case bool:
			rowCell.SetBool(cell.(bool))
		default:
			rowCell.SetString("unSupport type")
		}
	}
}
