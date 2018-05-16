package main

import "fmt"

func dummy(x int) int{
    return x
}

func main(){
    for i :=0; i < 10; i++ {
        fmt.Println(i)
    }

    test := dummy;
    test(3)
}
