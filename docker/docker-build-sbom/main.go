package main

import "fmt"
import "github.com/reactivex/rxgo/v2"

func main() {
    observable := rxgo.Just("Hello, World!")()
    ch := observable.Observe()
    item := <-ch
    fmt.Println(item.V)
}
