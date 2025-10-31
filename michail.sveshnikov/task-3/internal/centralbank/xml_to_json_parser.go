package centralbank

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type (
	Currency struct {
		NumCode  string  `json:"num_code"`
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
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read XML file: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("cannot convert encoding: %w", err)
	}

	var tempValCursData tempValCurs
	err = xml.Unmarshal(utf8Data, &tempValCursData)
	if err != nil {
		return nil, fmt.Errorf("cannot parse XML: %w", err)
	}

	var currencies []Currency
	for _, tempValute := range tempValCursData.Valutes {
		valueStr := strings.Replace(tempValute.Value, ",", ".", -1)

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse Value for currency %s: %w", tempValute.CharCode, err)
		}

		currencies = append(currencies, Currency{
			NumCode:  tempValute.NumCode,
			CharCode: tempValute.CharCode,
			Value:    value,
		})
	}

	return currencies, nil
}
