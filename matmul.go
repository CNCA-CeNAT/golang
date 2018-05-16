package main

import (
    "fmt"
    "os"
    "strconv"
)

var m1, m2, m3 []float64

func matmul(begin, end int, size int, link chan bool){

    for i := begin; i < end; i++{
        for j := 0; j < size; j++{
            sum := 0.0
            for k := 0; k < size; k++{
                sum += m1[i*size+k] * m2[j*size+k]
            }
            m3[i*size+j] = sum;
        }
    }
    link<- true
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Wrong number of parameters")
        panic(fmt.Sprintf("Usage: $ %s size num_gophers\n", os.Args[0]))
    }
    size, err := strconv.Atoi(os.Args[1])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[1]))
    }

    num_gophers, err := strconv.Atoi(os.Args[2])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[2]))
    }

    m1 = make([]float64, size * size)
    m2 = make([]float64, size * size)
    m3 = make([]float64, size * size)

    rows_per_gopher := size / num_gophers
    index := 0
    link := make(chan bool)
    for g := 0; g < num_gophers; g++ {
        go matmul(index, index + rows_per_gopher, size, link)
        index += rows_per_gopher
    }

    for g := 0; g < num_gophers; g++ {
        <-link
    }

}
