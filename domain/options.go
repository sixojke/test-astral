package domain

import "strconv"

const (
	defaultLimit = 10
	defaultPage  = 1
)

type FilterParams struct {
	Key    string
	Value  string
	Limit  int
	Offset int
}

func PrepareFillterParams(key, value string, limit, page string) *FilterParams {
	lim, _ := strconv.Atoi(limit)
	if lim <= 0 {
		lim = defaultLimit
	}

	pag, _ := strconv.Atoi(page)
	if pag <= 0 {
		pag = defaultPage
	}

	return &FilterParams{
		Key:    key,
		Value:  value,
		Limit:  lim,
		Offset: (pag - 1) * lim,
	}
}
