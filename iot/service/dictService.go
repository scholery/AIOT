package service

import (
	"koudai-box/cache"

	"koudai-box/iot/db"
	"koudai-box/iot/web/dto"
)

const (
	DICT_CACHE_KEY string = "DICT_CACHE"
)

func InitDictCache() {
	cache.Delete(DICT_CACHE_KEY)
	dicts, _ := QueryAllDicts()
	c := make(map[string][]dto.DictCodeItem)
	for _, v := range dicts {
		c[v.Type] = append(c[v.Type], dto.DictCodeItem{
			Id:    v.Id,
			Pid:   v.Pid,
			Name:  v.Name,
			Value: v.Value,
			Sort:  v.Sort,
			Extra: v.Extra,
		})
	}
	cache.SetWithNoExpire(DICT_CACHE_KEY, c)
}

func GetDictCache() map[string][]dto.DictCodeItem {
	c, _ := cache.Get(DICT_CACHE_KEY)
	dictCache := c.(map[string][]dto.DictCodeItem)
	return dictCache
}

func GetDictName(typ string, value string) string {
	if typ == "" || value == "" {
		return ""
	}
	c := GetDictCache()
	tyCache := c[typ]
	for _, d := range tyCache {
		if d.Value == value {
			return d.Name
		}
	}
	return ""
}

func DictCodesService(codes []string) (dto.DictCodesResponse, error) {
	response := make(dto.DictCodesResponse)
	if len(codes) == 0 {
		return response, nil
	}
	c, _ := cache.Get(DICT_CACHE_KEY)
	dictCache := c.(map[string][]dto.DictCodeItem)
	for _, v := range codes {
		response[v] = dictCache[v]
	}
	return response, nil
}

func QueryAllDicts() ([]*db.Dict, error) {
	return db.ListDictAll()
}
