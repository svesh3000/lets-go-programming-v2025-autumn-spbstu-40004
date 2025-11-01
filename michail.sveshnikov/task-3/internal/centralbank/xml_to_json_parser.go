package centralbank

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/net/html/charset"
)

type (
	Currency struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}

	tempValCurs struct {
		XMLName xml.Name     `xml:"ValCurs"`
		Valutes []tempValute `xml:"Valute"`
	}

	tempValute struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
)

func ParseXMLFile(filename string) ([]Currency, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open XML file: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			_ = closeErr
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read XML file: %w", err)
	}

	content = bytes.ReplaceAll(content, []byte(","), []byte("."))

	decoder := xml.NewDecoder(bytes.NewReader(content))
	decoder.CharsetReader = charset.NewReaderLabel

	var tempValCursData tempValCurs

	err = decoder.Decode(&tempValCursData)
	if err != nil {
		return nil, fmt.Errorf("cannot parse XML: %w", err)
	}

	currencies := make([]Currency, 0, len(tempValCursData.Valutes))

	for _, tempValute := range tempValCursData.Valutes {
		// Обработка NumCode - если пустой, устанавливаем 0
		numCode := 0
		if tempValute.NumCode != "" {
			var err error
			numCode, err = strconv.Atoi(tempValute.NumCode)
			if err != nil {
				return nil, fmt.Errorf("cannot parse NumCode for currency %s: %w", tempValute.CharCode, err)
			}
		}

		// Обработка CharCode - если пустой, оставляем пустую строку
		charCode := tempValute.CharCode
		if charCode == "" {
			charCode = ""
		}

		value, err := strconv.ParseFloat(tempValute.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse Value for currency %s: %w", tempValute.CharCode, err)
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: charCode,
			Value:    value,
		})
	}

	return currencies, nil
}
