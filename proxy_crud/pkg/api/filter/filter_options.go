package filter

type IFOptions interface {
	Limit() int
	Apply() bool
	AddField(key string, value any, operator string, dtype string)
	Fields() []Field
}

func (o *FOptions) Limit() int {
	return o.limit
}
func (o *FOptions) Apply() bool {
	return o.apply
}
func (o *FOptions) AddField(key string, value any, operator, dtype string) {
	o.fields = append(o.fields, Field{
		Key:      key,
		Value:    value,
		Operator: operator,
		Type:     dtype,
	})
}
func (o *FOptions) Fields() []Field {
	return o.fields
}
