package markov

import (
    "strings"
    "bytes"
    "math/rand"
    "io"
    "bufio"
)

type textLexer struct {}
func (lex textLexer) Next(b *bufio.Reader) (string, error) {
    c, err := peekByte(b)
    var buf bytes.Buffer

    for ; err == nil && (isLetter(c) || c == '\''); c, err = peekByte(b) {
        buf.WriteByte(c)
        b.ReadByte()
    }

    if buf.Len() > 0 {
        return buf.String(), nil
    }

    if err != nil {
        return "", err
    }

    b.ReadByte()
    if c == '.' || c == ',' || c == '-' {
        return string([]byte{c}), nil
    } else {
        return "", nil
    }
}

var TextLexer textLexer = textLexer{}

func isLetter(c byte) bool {
    return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func peekByte(b *bufio.Reader) (byte, error) {
    bytes, err := b.Peek(1)
    if err != nil {
        return 0, err
    } else {
        return bytes[0], nil
    }
}

type Lexer interface {
    Next(*bufio.Reader) (string, error)
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

const prefixSep = "\000"

func getKey(prev []string) string {
    return strings.Join(prev, prefixSep)
}

func (m *Markov) Add(prev []string, second string) {
    key := getKey(prev)
    entries := m.States[key]
    if entries == nil {
        entries = make([]*entry, 0, 1)
        m.States[key] = entries
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
    m.States[key] = append(entries, &entry)
}

func CreateMarkov(src io.Reader, lex Lexer, order int) *Markov {
    input := bufio.NewReader(src)

    m := Markov{States: make(map[string][]*entry)}
    m.Order = order

    prev := make([]string, order)
    for i := 0; i < order; i++ {
        prev[i] = ""
    }

    for tok, err := lex.Next(input); err == nil; tok, err = lex.Next(input) {
        if tok == "" {
            continue
        }

        m.Add(prev, tok)
        prev = append(prev[1:], tok)
    }

    return &m
}

func (m *Markov) Pick(prefix []string) string {
    key := getKey(prefix)
    entries := m.States[key]
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
    return ""
}

func (m *Markov) Generate() string {
    var buf bytes.Buffer

    // pick a random starting prefix
    var start string
    for k, _ := range(m.States) {
        start = k
        break
    }

    arr := strings.Split(start, prefixSep)
    first := true
    for {
        word := m.Pick(arr)
        if word == "" {
            break
        }
        if word != "," && word != "." {
            if first {
                first = false
            } else {
                buf.WriteString(" ")
            }
        }
        buf.WriteString(word)

        arr = append(arr[1:], word)

        if word == "." {
            break
        }
    }
    return buf.String()
}
