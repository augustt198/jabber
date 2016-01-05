package markov

import (
    "strings"
    "bytes"
    "math/rand"
    "io"

    "container/list"
)

type Lexer interface {
    Next() (string, error)
}

type entry struct {
    Freq int
    Word string
}

type Markov struct {
    States map[string][]*entry
    // order/memory of this markov chain, meaning how
    // many previous states the next state depends on
    Order int
}

func (m *Markov) Add(first, second string) {
    entries := m.States[first]
    if entries == nil {
        entries = make([]*entry, 0, 1)
        m.States[first] = entries
    }

    for _, e := range(entries) {
        if strings.EqualFold(e.Word, second) {
            e.Freq += 1
            return
        }
    }

    entry := entry{
        Freq: 1,
        Word: second,
    }
    m.States[first] = append(entries, &entry)
}

func isLetter(c byte) bool {
    return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func CreateMarkov(src io.Reader, lex Lexer, order int) *Markov {
    m := Markov{States: make(map[string][]*entry)}
    m.Order = order

    prev := list.New()
    for i := 0; i < order; i++ {
        prev.PushBack("")
    }

    for tok, err := lex.Next(src); err != nil {
        if tok == "" {
            continue
        }

        var buf bytes.Buffer
        for i < len(source) && isLetter(source[i]) {
            buf.WriteByte(source[i])
            i += 1
        }
        if buf.Len() > 0 {
            word := buf.String()
            m.Add(prev, word)
            prev = word
            continue
        }
        if i >= len(source) { break }

        current := source[i]
        if current != ' ' {
            word := string([]byte{current})
            m.Add(prev, word)
            prev = word
        }
        i += 1
    }

    return &m
}

func (m *Markov) Pick(word string) string {
    entries := m.States[word]
    if entries == nil {
        return ""
    }

    sum := 0
    for _, e := range(entries) {
        sum += e.Freq
    }

    r := rand.Intn(sum)

    curr := 0
    for _, e := range(entries) {
        if r >= curr && r < (curr + e.Freq) {
            return e.Word
        }
        curr += e.Freq
    }
    return "ERR"
}

func (m *Markov) Generate(start string) string {
    var buf bytes.Buffer
    prev := start
    for {
        word := m.Pick(prev)
        if word == "" {
            break
        }
        if isLetter(word[0]) {
            buf.WriteByte(' ')
        }
        buf.WriteString(word)

        if word == "." {
            break
        }
        prev = word
    }
    return buf.String()
}
