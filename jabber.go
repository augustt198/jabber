package main

import (
    "github.com/augustt198/jabber/markov"

    "fmt"
    "os"
    "io/ioutil"
    "math/rand"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("usage: %s [text file]\n", os.Args[0])
        os.Exit(-1)
    }

    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }

    rand.Seed(time.Now().UTC().UnixNano())


    m := markov.CreateMarkov(string(data))
    /*
    for k, v := range(m.States) {
        fmt.Printf("'%s'\n", k)
        for _, e := range(v) {
            fmt.Printf("   '%s' -> %d\n", e.Word, e.Freq)
        }
    }
    */

    for i := 0; i < 100; i++ {
        fmt.Println(m.Generate("\n"))
    }
}
