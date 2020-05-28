package main

import "fmt"

func main() {
	var x,y []int
	for i:=0;i<10;i++{
		y = append(x,i)
		fmt.Println(i,"----",x,"----",y,"----",cap(x),"----",cap(y),"----",&x,"----",&y)
	}
}
