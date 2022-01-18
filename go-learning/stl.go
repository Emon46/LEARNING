package main

import (
	"fmt"
	"sort"
	"strings"
)

// ---------------------start stack--------------------
type Stack []string

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// return top element of stack. Return false if stack is empty.
func (s *Stack) Top() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element, true
	}
}

// ---------------------start stack--------------------

// ----------------------
func main() {
	var stack Stack // create a stack variable of type Stack

	stack.Push("this")
	stack.Push("is")
	stack.Push("sparta!!")

	top, ok := stack.Top()
	if ok == true {
		fmt.Println(top)
	}

	for len(stack) > 0 {
		x, y := stack.Pop()
		if y == true {
			fmt.Println(x)
		}
	}

	//--------------------------sort------------------------
	s := "hello world"
	elements := strings.Split(s, " ")
	fmt.Println(elements)
	//output :  [hello, world]

	arr := []int{6, 7, 1, 2, 3, 4}
	sort.Ints(arr)
	fmt.Println(arr)
	// [1 2 3 4 6 7]

	sort.SliceStable(arr, func(i, j int) bool {
		//i,j are represented for two value of the slice .
		return arr[i] > arr[j]
	})
	fmt.Println(arr)
	// [7 6 4 3 2 1]

	/*
		sort.Slice(matrix[:], func(i, j int) bool {
		    for x := range matrix[i] {
		        if matrix[i][x] == matrix[j][x] {
		            continue
		        }
		        return matrix[i][x] < matrix[j][x]
		    }
		    return false
		})

		fmt.Println(matrix)
	*/

	var arr2d [7][4]int
	fmt.Println(arr2d)
	// [[0 0 0 0] [0 0 0 0] [0 0 0 0] [0 0 0 0] [0 0 0 0] [0 0 0 0] [0 0 0 0]]

}
