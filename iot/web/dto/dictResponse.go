package dto

type DictCodeItem struct {
	Id    int    `json:"id"`
	Pid   int    `json:"pid"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Sort  int    `json:"sort"`
	Extra string `json:"extra"`
}

type DictCodesResponse = map[string][]DictCodeItem
