package filter

type Options struct {
	SortOptions   SOptions
	FilterOptions IFOptions
}

type FOptions struct {
	apply  bool
	limit  int
	fields []Field
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
