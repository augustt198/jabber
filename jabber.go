package main

import (
    "github.com/augustt198/jabber/markov"

    "fmt"
    "os"
    "io"
    "math/rand"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("usage: %s [text files]\n", os.Args[0])
        os.Exit(-1)
    }

    readers := make([]io.Reader, len(os.Args) - 1)
    for i := 1; i < len(os.Args); i++ {
        f, err := os.Open(os.Args[i])
        if err != nil {
            fmt.Println(err)
            os.Exit(-1)
        }
        readers[i - 1] = f
    }
    input := io.MultiReader(readers...)

    rand.Seed(time.Now().UTC().UnixNano())

    m := markov.CreateMarkov(input, markov.TextLexer, 2)

    for i := 0; i < 100; i++ {
        fmt.Printf("> %s\n", m.Generate())
    }
}
