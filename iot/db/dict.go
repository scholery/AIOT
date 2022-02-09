package db

import (
	"strconv"
	"strings"
)

func ListDictWithCodes(types []string) ([]*Dict, error) {
	var placeholders []string
	for i := 0; i < len(types); i++ {
		placeholders = append(placeholders, "?")
	}
	var dicts []*Dict
	_, err := webOrm.Raw("select id,type,num,key,value,tips,extra from dict where type in ("+strings.Join(placeholders, ",")+")", types).QueryRows(&dicts)
	return dicts, err
}

func ListDictWithIds(ids []int) ([]*Dict, error) {
	var placeholders []string
	for i := 0; i < len(ids); i++ {
		placeholders = append(placeholders, "?")
	}
	var dicts []*Dict
	_, err := webOrm.Raw("select id,type,num,key,value,tips,extra from dict where id in ("+strings.Join(placeholders, ",")+") order by type", ids).QueryRows(&dicts)
	return dicts, err
}

func ListDictAll() ([]*Dict, error) {
	var dicts []*Dict
	_, err := webOrm.QueryTable("dict").All(&dicts)

	return dicts, err
}

//删除算法关联的事件类型和告警类型
func DeleteDictByAlg(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	strArr := make([]string, 0)
	for _, v := range ids {
		strArr = append(strArr, strconv.Itoa(v))
	}
	types := []string{"event_type", "warning_type"}
	return webOrm.QueryTable("dict").Filter("extra__in", strArr).Filter("type__in", types).Delete()
}

func InsertDictByAlg(algId int, dicts []Dict) {
	_, err := webOrm.InsertMulti(len(dicts), &dicts)
	if err != nil {
		logger.Errorln(err)
	}
}
