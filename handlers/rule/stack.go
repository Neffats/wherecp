package rulehandler

type fnStack struct {
	data []string
}

func NewfnStack() *fnStack{
	return &fnStack{
		data: make([]string, 0),
	}
}

func (s *fnStack) Push(item string) {
	s.data = append(s.data, item)
}

func (s *fnStack) Pop() string {
	if len(s.data) == 0 {
		return ""
	}

	item := s.data[len(s.data)-1]
	s.data = s.data[0:len(s.data)-1]
	return item
}

func (s *fnStack) Peek() string {
	if len(s.data) == 0 {
		return ""
	}
	return s.data[len(s.data)-1]
}
