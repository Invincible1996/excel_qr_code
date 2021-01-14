package model

type ClassroomModel struct {
	Name string
	Code string
}

func (c *ClassroomModel) setName(name string) {
	c.Name = name
}

func (c *ClassroomModel) setCode(code string) {
	c.Code = code
}
