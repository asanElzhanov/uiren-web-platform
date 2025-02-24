package metadata

type (
	Data struct {
		DataType         map[string]interface{} `json:"data_type"`
		Description      string                 `json:"description"`
		IsMarkedToDelete bool                   `json:"is_marked_to_delete"`
		IsRequired       bool                   `json:"is_required"`
		IsUnique         bool                   `json:"is_unique"`
		RequisiteName    string                 `json:"requisite_name"`
		SortOrder        float64                `json:"sort_order"`
		TitleKaz         string                 `json:"title_kaz"`
		TitleRus         string                 `json:"title_rus"`
	}
	SystemData struct {
		Action    string  `json:"action"`
		CreatedAt float64 `json:"created_at"`
		DeletedAt float64 `json:"deleted_at"`
	}
	Metadata struct {
		Data       Data       `json:"data"`
		SystemData SystemData `json:"system_data"`
	}
)
