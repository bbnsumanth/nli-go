package main

import "fmt"
import "strings"
import "os"
import "nli-go/lib"

// Provide a sentence as command line parameters (or as a single parameter within quotes)
// and this app will provide the tokens, separated by slashes
func main() {

    rawInput := strings.Join(os.Args[1:], " ")

    if len(rawInput) == 0 {
        fmt.Print("use: tokenizer \"Provide a sentence here\"\n")
        return
    }

    inputSource := lib.NewSimpleRawInputSource(rawInput)

    tokenizer := lib.NewSimpleTokenizer()
    wordArray := tokenizer.Process(inputSource)

    fmt.Print(strings.Join(wordArray, "/"))
    fmt.Print("\n")
}
