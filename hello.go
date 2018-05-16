package main

import "fmt"

func say_hello(i int, c chan int) {
    fmt.Printf("Hello from %d!\n", i)
    c <- i
}

func main() {

    c := make(chan int)

    for i := 0; i < 10; i++ {
        go say_hello(i, c)
    }
    for i := 0; i < 10; i++ {
        <- c
    }
}
