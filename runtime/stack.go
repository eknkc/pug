package runtime

type stack struct {
	values []interface{}
}

func (s *stack) pop() interface{} {
	if len(s.values) > 0 {
		val := s.values[len(s.values)-1]
		s.values = s.values[:len(s.values)-1]
		return val
	}

	return nil
}

func (s *stack) push(val ...interface{}) {
	s.values = append(s.values, val...)
}

func (s *stack) len() int {
	return len(s.values)
}
