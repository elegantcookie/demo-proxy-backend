package filter

type Options struct {
	SortOptions   SOptions
	FilterOptions IFOptions
}

type FOptions struct {
	apply  bool
	limit  int
	page   int
	fields []Field
}

func NewFOptions(apply bool, limit, page int, fields []Field) *FOptions {
	return &FOptions{
		apply:  apply,
		limit:  limit,
		page:   page,
		fields: fields,
	}
}

type Field struct {
	Key      string
	Value    any
	Operator string
	Type     string
}

type SOptions struct {
	Field string
	Order string
}
