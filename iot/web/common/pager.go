package common

func Page2Offset(pageNo, pageSize int) (offset, limit int) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset = (pageNo - 1) * pageSize
	limit = pageSize
	return offset, limit
}
