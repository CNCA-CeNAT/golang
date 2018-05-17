package main

import (
    "fmt"
    "os"
    "strconv"
)

var iterations int
var stencil0, stencil1 []float64

func index_to_coord(index, size int) (x, y, z int) {
    x = index / (size * size)
    y = (index % (size*size)) / size
    z = (index % (size*size)) % size

    return x, y, z
}

func coord_to_index(x, y, z, size int) (int) {
    return x * (size*size) + y * size + z
}

func point_stencil(s []float64, x, y, z, size int) (float64) {
    return (s[coord_to_index(x, y, z, size)] +
            s[coord_to_index(x-1, y, z, size)] +
            s[coord_to_index(x+1, y, z, size)] +
            s[coord_to_index(x, y-1, z, size)] +
            s[coord_to_index(x, y+1, z, size)] +
            s[coord_to_index(x, y, z-1, size)] +
            s[coord_to_index(x, y, z+1, size)] ) / 7.0;
}

func stencil(index, size, ppg int, link chan int, done chan bool){
    var s_in, s_out []float64

    for it := <-link; it >= 0; it = <-link {
        if it % 2 == 0 {
            s_out = stencil0;
            s_in  = stencil1;
        } else {
            s_out = stencil1;
            s_in  = stencil0;
        }

        i := index
        for (i-index) < ppg {
            x, y, z := index_to_coord(index, size)
            if x == 0 || y == 0 || z == 0 {
               continue
            }
            if x == size-1 || y == size-1 || z == size-1{
                continue
            }
            s_in[index] = point_stencil(s_out, x, y, z, size)
            i++
        }
        done <- true
    }
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Wrong number of parameters")
        panic(fmt.Sprintf("Usage: $ %s size iterations num_gophers\n", os.Args[0]))
    }

    var err error
    size, err := strconv.Atoi(os.Args[1])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[1]))
    }
    iterations, err = strconv.Atoi(os.Args[2])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[2]))
    }
    num_gophers, err := strconv.Atoi(os.Args[3])
    if err != nil {
        panic(fmt.Sprintf("%s isn't an integer", os.Args[3]))
    }

    stencil0 = make([]float64, size*size*size)
    stencil1 = make([]float64, size*size*size)

    points_per_gopher := (size-1) * (size-1) * (size-1) / num_gophers

    next_it := make(chan int)
    done    := make(chan bool)
    index := 0
    for i := 0; i < num_gophers; i++ {
        x, y, z := index_to_coord(index, size-1)
        external_index := coord_to_index(x+1, y+1, z+1, size)
        if i == num_gophers-1 {
            points_per_gopher = (size-1)*(size-1)*(size-1) - index
            go stencil(external_index, size, points_per_gopher, next_it, done)
        } else {
            go stencil(external_index, size, points_per_gopher, next_it, done)
        }
        index += points_per_gopher
    }

    // sinc, equivalent to a barrier
    for it := 0; it < iterations; it ++{
        for i := 0; i < num_gophers; i++ {
            next_it <- it
        }
        for i := 0; i < num_gophers; i++ {
            <-done
        }
    }
    // signal gophers to en
    for i := 0; i < num_gophers; i++ {
        next_it <- -1
    }
}
