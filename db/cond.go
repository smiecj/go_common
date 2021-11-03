// package db db common method: db query condition format, db connector...
package db

// db where condition json format
// example: cond_example.json

type whereCondition struct {
	Type   string `json:"type"`
	Key    string `json:"key"`
	Method string `json:"method"`
	Value  string `json:"value"`
}

type SearchCondition struct {
	Order struct {
		Field string `json:"field"`
		Sc    string `json:"sc"`
	} `json:"order"`
	Page struct {
		No    int `json:"no"`
		Limit int `json:"limit"`
	} `json:"page"`
	Where []whereCondition `json:"where"`
}

type UpdateCondition struct {
	Where []whereCondition `json:"where"`
}
