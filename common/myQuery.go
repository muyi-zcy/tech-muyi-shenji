package common

var MAX_PAGE_SIZE = 2000

type MyQuery struct {
	Size    int   `json:"size"`
	Current int   `json:"current"`
	Total   int64 `json:"total"`
}

func (q *MyQuery) GetSize() int {
	if q.Size == 0 {
		return 20
	}
	if q.Size > MAX_PAGE_SIZE {
		return MAX_PAGE_SIZE
	}
	return q.Size
}

func (q *MyQuery) GetCurrent() int {
	if q.Current <= 0 {
		return 1
	}
	return q.Current
}

func (q *MyQuery) GetOffset() int {
	return (q.GetCurrent() - 1) * q.GetSize()
}
