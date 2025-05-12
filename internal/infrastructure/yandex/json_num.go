package yandex

import (
	"encoding/json"
	"strconv"
)

// Num — строка-контейнер для координаты, умеет читать "37.27" или 37.27.
type Num string

func (n *Num) UnmarshalJSON(b []byte) error {
	// null → пустая строка
	if string(b) == "null" {
		*n = ""
		return nil
	}
	// число без кавычек
	if b[0] != '"' {
		*n = Num(b)
		return nil
	}
	// обычная строка
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*n = Num(s)
	return nil
}

// Float64 конвертирует в число; пустая строка → 0, false.
func (n Num) Float64() (float64, bool) {
	if n == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(string(n), 64)
	return f, err == nil
}
