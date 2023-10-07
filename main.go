package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/user"

    "monkey/evaluator"
    "monkey/lexer"
    "monkey/object"
    "monkey/parser"
    "monkey/repl/repl"
)

func main() {
    out := os.Stdout
    if len(os.Args) > 1 {
        env := object.NewEnvironment()
        filename := os.Args[1]
        file, err := os.Open(filename)
        if err != nil {
            fmt.Println("Error opening file: ", err)
            return
        }
        defer file.Close()
        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanBytes)
        scanner.Scan()
        line := scanner.Text()
        for scanner.Scan() {
            line += scanner.Text()
        }
        l := lexer.New(line)
        p := parser.New(l)

        program := p.ParseProgram()
        if len(p.Errors()) != 0 {
            repl.PrintParserErrors(out, p.Errors())
            return
        }

        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
        return
    }

    currentUser, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! This is the Monkey programming language!\n", currentUser.Username)
    fmt.Printf("Feel free to type in commands\n")
    repl.Start(os.Stdin, out)
}
