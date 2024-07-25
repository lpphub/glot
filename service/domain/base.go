package domain

type PageQuery struct {
	Pn int `json:"current" form:"current"`
	Ps int `json:"size" form:"size"`
}

type Pager struct {
	Page  int         `json:"current,omitempty"`
	Size  int         `json:"size,omitempty"`
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type EmptyList []interface{}

func WrapPager(total int64, list interface{}) *Pager {
	return &Pager{
		Total: total,
		List:  list,
	}
}
