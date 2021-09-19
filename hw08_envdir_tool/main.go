package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func printHelp() {
	fmt.Println("go-envdir - runs another program with environment modified according to files in a specified directory.")
	fmt.Println("cmd: go-envdir <dir> <child>")
	fmt.Println("	dir - is a single argument (required)")
	fmt.Println("	child - consists of one or more arguments (required)")
	fmt.Println("cmd with debug: go-envdir -v <dir> <child>")
}

func main() {
	debug := flag.Bool("v", false, "enable debug mode")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		printHelp()
		os.Exit(111)
	}
	dir := args[0]
	child := args[1:]

	envs, err := ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	if *debug {
		fmt.Println("*** Debug mode is enabled ***")
		fmt.Println("additional env vars:")
		for k, v := range envs {
			if v.NeedRemove {
				fmt.Printf("%s will be unset\n", k)
			} else {
				fmt.Printf("%s=%s\n", k, v.Value)
			}
		}
	}

	code := RunCmd(child, envs)
	os.Exit(code)
}
