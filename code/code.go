package code

import "os"

type Code struct {
	frags []string
}

func (c *Code) AddLine(line string) {
	c.frags = append(c.frags, line)
}

func (c *Code) AddLines(other *Code) {
	c.frags = append(c.frags, other.frags...)
}

func (c *Code) SkipLine() {
	c.frags = append(c.frags, "")
}

func (c *Code) WriteTo(dst string) {
	var str string
	for _, s := range c.frags {
		str += s + "\n"
	}
	os.WriteFile(dst, []byte(str), 0644)
}

func NewCode() *Code {
	var c Code
	c.frags = make([]string, 1)
	return &c
}
