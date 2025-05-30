package common

type Paging struct {
	Page  int 	`json:"page" form:"page"`
	Limit int 	`json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func(p *Paging) ProcessPaging() {
	if p.Page < 1 {
		p.Page = 1
	}
	
	if p.Limit < 1 {
		p.Limit = 10
	}
}

