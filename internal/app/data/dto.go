package data

import (
	"errors"
	"strconv"
	"time"
	"uiren/pkg/logger"
)

type DictionaryData struct {
	Data       map[string]interface{} `json:"data"`
	SystemData SystemData             `json:"system_data"`
}

type SystemData struct {
	Created_at float64 `json:"created_at"`
	Deleted_at float64 `json:"deleted_at"`
}

func (intData *DictionaryData) NormalizeTypes(metadataTypes map[string]string) error {
	for key, value := range intData.Data {
		if value == "" {
			continue
		}
		value, ok := value.(string)
		if !ok {
			return errors.New("bad request maybe")
		}
		switch metadataTypes[key] {
		case "number":
			newValue, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			intData.Data[key] = newValue

		case "boolean":
			newValue, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			intData.Data[key] = newValue

		case "datetime":
			newValue, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			intData.Data[key] = newValue
		case "string":
			//nothing
		case "reference":
			//some logic
		default:
			logger.Error(key, " ", value)
			logger.Error(metadataTypes)
			return errors.New("Неизвестный тип данных")
		}
	}
	return nil
}
