package common

import (
	"encoding/json"
	"strconv"
)

type Response struct {
	Datas      interface{} `json:"datas"`
	ResultCode string      `json:"resultCode"`
	ResultMsg  string      `json:"resultMsg"`
	Success    bool        `json:"success"`
}

type PageData struct {
	PageNo    int         `json:"pageNo"`
	PageSize  int         `json:"pageSize"`
	Rows      interface{} `json:"rows"`
	TotalSize int64       `json:"totalSize"`
}

type PageResponse struct {
	Datas      PageData `json:"datas"`
	ResultCode string   `json:"resultCode"`
	ResultMsg  string   `json:"resultMsg"`
	Success    bool     `json:"success"`
}

type ImportResult struct {
	SuccessCount int
	FailCount    int
	FailedRows   []int
}

func (res ImportResult) String() string {
	msg := "导入完成"
	if res.SuccessCount > 0 {
		msg += "，导入成功 " + strconv.Itoa(res.SuccessCount) + " 条"
	}
	if res.FailCount > 0 {
		msg += "，导入失败 " + strconv.Itoa(res.FailCount) + " 条，失败的行号为："
		bytes, _ := json.Marshal(res.FailedRows)
		msg += string(bytes)
	}
	return msg
}

func Error(msg string, datas interface{}) Response {
	return Response{datas, STATUS_ERROR, msg, false}
}

func Ok(msg string, datas interface{}) Response {
	return Response{datas, STATUS_OK, msg, true}
}

func OkPage(pageNo, pageSize int, totalSize int64, msg string, rows interface{}) PageResponse {
	pageData := PageData{
		PageNo:    pageNo,
		PageSize:  pageSize,
		Rows:      rows,
		TotalSize: totalSize,
	}
	return PageResponse{
		Datas:      pageData,
		ResultCode: "0",
		ResultMsg:  msg,
		Success:    true,
	}
}

func ErrorPage(pageNo, pageSize int, totalSize int64, msg string) PageResponse {
	pageData := PageData{
		PageNo:    pageNo,
		PageSize:  pageSize,
		Rows:      nil,
		TotalSize: totalSize,
	}
	return PageResponse{
		Datas:      pageData,
		ResultCode: "1",
		ResultMsg:  msg,
		Success:    false,
	}
}
