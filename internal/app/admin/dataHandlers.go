package admin

/*
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"uiren/internal/app/data"
	"uiren/pkg/logger"

	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func (app *App) exportData(w http.ResponseWriter, r *http.Request) {
	var (
		dictIDString = r.URL.Query().Get("dict_id")
		ctx          = r.Context()
	)
	if dictIDString == "" {
		logger.Error("Нет dict_id")
		http.Error(w, "dictID required", http.StatusBadRequest)
		return
	}
	dictID, err := strconv.Atoi(dictIDString)
	if err != nil {
		logger.Error("Ошибка Atoi: ", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=data.csv")
	w.Header().Set("Content-Type", "text/csv")

	/*
		4 - кадастриал
		1- цвета
		2 - тест всего
		5 - тест адресов
*/ /*
	dataList, err := app.dataService.GetDictionaryData(ctx, dictID)
	if err != nil {
		logger.Error("Ошибка GetDictionaryData: ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	if len(dataList) == 0 {
		logger.Error("Нет данных")
		return
	}
	metadataSortOrder, err := app.metadataService.GetMetadataSortOrder(ctx, dictID)
	if err != nil {
		logger.Error("Ошибка GetMetadata: ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	keys := make([]string, 0, len(dataList[0].Data))
	for key := range dataList[0].Data {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return metadataSortOrder[keys[i]] < metadataSortOrder[keys[j]]
	})
	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write(keys)
	for _, dataEntity := range dataList {
		row := make([]string, 0, len(dataEntity.Data))
		for _, key := range keys {
			if dataEntity.Data[key] == nil {
				row = append(row, "")
			} else {
				row = append(row, fmt.Sprintf("%v", dataEntity.Data[key]))

			}
		}
		writer.Write(row)
	}
}

func (app *App) importData(w http.ResponseWriter, r *http.Request) {
	var (
		dictIDString = r.URL.Query().Get("dict_id")
		ctx          = r.Context()
	)
	dictID, err := strconv.Atoi(dictIDString)
	if err != nil {
		logger.Error("Ошибка Atoi: ", err)
		http.Error(w, "bad request", http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Error("Ошибка FormFile: ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		logger.Error("Ошибка Read(заголовки): ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	/*
		4 - кадастриал
		1- цвета
		2 - тест всего
		5 - тест адресов
*/ /*
	metadataTypes, err := app.metadataService.GetMetadataTypes(ctx, dictID)
	if err != nil {
		logger.Error("Ошибка GetMetadataTypes: ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	var dictDataList []data.DictionaryData
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("Ошибка Read: ", err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		var dictData data.DictionaryData
		dictData.Data = make(map[string]interface{})
		for i, value := range record {
			dictData.Data[headers[i]] = value
		}
		if err := dictData.NormalizeTypes(metadataTypes); err != nil {
			logger.Error("Ошибка Unmarshal: ", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		dictData.SystemData.Created_at = float64(time.Now().Unix())
		dictDataList = append(dictDataList, dictData)
	}

	for _, dictData := range dictDataList {
		id, err := app.dataService.InsertData(ctx, dictID, dictData)
		if err != nil {
			logger.Error("Ошибка InsertData: ", err)
			http.Error(w, "Внутренняя ошибка", http.StatusBadRequest)
			return
		}
		logger.Error(id)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dictDataList); err != nil {
		logger.Error("Ошибка Encode: ", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}
}
*/
