package main

import (
    "fmt"
    "os"
    "strconv"
    "math/rand"
    "math"
    "time"
)

func partial_pi(num_points int, link chan int){

    source := rand.NewSource(time.Now().UnixNano())
    generator := rand.New(source)

    count := 0
    for i := 0; i < num_points; i++{
        r1 := generator.Float64()
        r2 := generator.Float64()
        if math.Hypot(r1, r2) < 1{
            count++
        }
    }
    link <- count
}

func main(){
    /* Read arguments from command line */
    if len(os.Args) != 3 {
        fmt.Println("Wrong number of parameters")
        panic(fmt.Sprintf("Usage: $ %s num_gophers num_points\n", os.Args[0]));
    }
    num_gophers, err := strconv.Atoi(os.Args[1])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[1]))
    }
    num_points, err := strconv.Atoi(os.Args[2])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[2]))
    }

    /* Create go routines */
    link := make(chan int)
    var points_per_gopher int = num_points / num_gophers
    for i := 0; i < num_gophers; i++{
        go partial_pi(points_per_gopher, link)
    }

    /* fetch partial results */
    approx_pi := 0;
    for i := 0; i < num_gophers; i++{
        approx_pi += <-link 
    }

    /* Print result */
    result := 4.0 * float64(approx_pi) / (float64(num_gophers) * float64(points_per_gopher))
    fmt.Printf("Result: %f\n", result)
}
