package model

type ClassroomModel struct {
	Name string
	Code string
}

func (c *ClassroomModel) SetName(name string) {
	c.Name = name
}

func (c *ClassroomModel) SetCode(code string) {
	c.Code = code
}
