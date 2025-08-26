package mocks

import (
	"errors"
	"reflect"
)

// MockRow представляет собой мок для sql.Row
type MockRow struct {
	values []interface{}
	index  int
}

// NewMockRow создает новый MockRow с заданными значениями
func NewMockRow(values ...interface{}) *MockRow {
	return &MockRow{
		values: values,
		index:  0,
	}
}

// Scan реализует интерфейс sql.Scanner
func (r *MockRow) Scan(dest ...interface{}) error {
	if r.index >= len(r.values) {
		return errors.New("no more values to scan")
	}

	for i, d := range dest {
		if r.index+i >= len(r.values) {
			return errors.New("not enough values to scan")
		}

		// Получаем значение из мока
		val := r.values[r.index+i]

		// Проверяем, что dest[i] является указателем
		destValue := reflect.ValueOf(d)
		if destValue.Kind() != reflect.Ptr {
			return errors.New("destination must be a pointer")
		}

		// Получаем значение, на которое указывает указатель
		destElem := destValue.Elem()

		// Устанавливаем значение
		if val != nil {
			valValue := reflect.ValueOf(val)
			if valValue.Type().AssignableTo(destElem.Type()) {
				destElem.Set(valValue)
			} else {
				// Попытка конвертации типов
				switch dest := d.(type) {
				case *string:
					if str, ok := val.(string); ok {
						*dest = str
					} else {
						*dest = ""
					}
				case *int:
					if num, ok := val.(int); ok {
						*dest = num
					} else {
						*dest = 0
					}
				case *int64:
					if num, ok := val.(int64); ok {
						*dest = num
					} else {
						*dest = 0
					}
				case *uint64:
					if num, ok := val.(uint64); ok {
						*dest = num
					} else {
						*dest = 0
					}
				case *bool:
					if b, ok := val.(bool); ok {
						*dest = b
					} else {
						*dest = false
					}
				case *float64:
					if f, ok := val.(float64); ok {
						*dest = f
					} else {
						*dest = 0.0
					}
				default:
					// Для других типов просто игнорируем
				}
			}
		}
	}

	r.index += len(dest)
	return nil
}

// MockRowWithError создает MockRow который возвращает ошибку при Scan
func MockRowWithError(err error) *MockRow {
	return &MockRow{
		values: []interface{}{err},
		index:  0,
	}
}

// MockRowNoRows создает MockRow который симулирует sql.ErrNoRows
func MockRowNoRows() *MockRow {
	return &MockRow{
		values: []interface{}{errors.New("sql: no rows in result set")},
		index:  0,
	}
}
