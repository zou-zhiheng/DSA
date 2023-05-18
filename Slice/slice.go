package Slice

import "fmt"

type slice struct {
	Array string
	len int
	cap int
}

func appendTest(slice []int,data ...int) []int {
	m:=len(slice)
	n:=m+len(data)

	if n>cap(slice) {
		newSlice:=make([]int,(n+1)*2)
		copy(newSlice,slice)
		slice=newSlice
	}
	slice=slice[:n]
	copy(slice[m:n],data)

	return slice
}

func Demo(){
	test:=[]int{1,2,3}
	test=appendTest(test,3)

	fmt.Println(test)
}
