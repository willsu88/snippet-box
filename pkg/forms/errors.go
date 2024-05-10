package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	errs := e[field]
	if len(errs) == 0 {
		return ""
	}
	return errs[0]
}
