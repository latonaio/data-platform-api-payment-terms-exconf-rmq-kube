package dpfm_api_output_formatter

type MetaData struct {
	// ConnectionKey        string               `json:"connection_key"`
	// Result               bool                 `json:"result"`
	// RedisKey             string               `json:"redis_key"`
	// Filepath             string               `json:"filepath"`
	// APIStatusCode        int                  `json:"api_status_code"`
	// RuntimeSessionID     string               `json:"runtime_session_id"`
	// ServiceLabel         string               `json:"service_label"`
	// APISchema            string               `json:"api_schema"`
	// Accepter             []string             `json:"accepter"`
	// OrderID              *int                 `json:"order_id"`
	// Deleted              bool                 `json:"deleted"`
	PaymentTerms PaymentTerms `json:"PaymentTerms"`
}

type PaymentTerms struct {
	PaymentTerms       string `json:"PaymentTerms"`
	ExistenceConf bool   `json:"ExistenceConf"`
}
