package types

type Stack[T any] struct{
	Elements []T
}

func (s *Stack[T]) Push(item T){
	s.Elements = append(s.Elements, item)
}

func (s *Stack[T]) Pop() *T{
  if len(s.Elements) == 0{
		return nil
	}

  n := len(s.Elements)-1 
  item := s.Elements[n]
  s.Elements = s.Elements[:n]
  return &item
}

func (s *Stack[T]) TopElement() *T{
	if len(s.Elements) == 0{
		return nil
	}

	return &s.Elements[len(s.Elements)-1]
}

func NewStack[T any]() Stack[T]{
	return Stack[T]{}
}
