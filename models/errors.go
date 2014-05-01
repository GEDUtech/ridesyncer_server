package models

type Errors struct {
	Overall map[string]string
	Fields  map[string]string
}

func NewErrors() *Errors {
	return &Errors{map[string]string{}, map[string]string{}}
}

func (errors *Errors) Count() int {
	return len(errors.Fields) + len(errors.Overall)
}
