package main

// The zero value for Stack is an empty stack ready to use.
type Stack struct {
	data []interface{}
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

// Push adds x to the top of the stack.
func (s *Stack) Push(x interface{}) {
	s.data = append(s.data, x)
}

//Returns item at top of stack without removing
func (s *Stack) Peek() (interface{}, bool) {
	if s.IsEmpty() {
		return "", false
	}
	return s.data[len(s.data) - 1], true	
}

// Pop removes and returns the top element of the stack.
// It’s a run-time error to call Pop on an empty stack.
func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return "", false
	}
	index := len(s.data) - 1   // Get the index of the top most element.
	element := s.data[index] // Index into the slice and obtain the element.
	s.data[index] = nil // to avoid memory leak
	s.data = s.data[:index] // Remove it from the stack by slicing it off.
	return element, true
}

// Size returns the number of elements in the stack.
func (s *Stack) Size() int {
	return len(s.data)
}

type MoveStack struct { Stack }

func (s *MoveStack) Push(item Move) { s.Stack.Push(item) }
func (s *MoveStack) Pop() (Move, bool)  {
	val, bol := s.Stack.Pop()
	if bol {
		return val.(Move), bol
	}
	return Move{}, bol
}
func (s *MoveStack) Peek() (Move, bool) {
	val, bol := s.Stack.Peek()
	if bol {
		return val.(Move), bol
	}
	return Move{}, bol
}
