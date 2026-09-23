package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "serve":
		if err := runServe(os.Args[2:]); err != nil {
			log.Fatalf("nomad serve failed: %v", err)
		}
	case "mcp-serve":
		if err := runMCPServe(os.Args[2:]); err != nil {
			log.Fatalf("nomad mcp-serve failed: %v", err)
		}
	case "session":
		if err := runSession(os.Args[2:]); err != nil {
			log.Fatalf("nomad session failed: %v", err)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Nomad - AI Agent Framework")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  nomad <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  session    Start an interactive AI agent session")
	fmt.Println("  serve      Start an HTTP server")
	fmt.Println("  mcp-serve  Start an MCP HTTP server")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  nomad session                    # Start interactive session")
	fmt.Println("  nomad session --recipe my.yaml   # Start with recipe")
	fmt.Println("  nomad serve --port 8080          # Start HTTP server")
	fmt.Println()
	fmt.Println("Use 'nomad <command> -h' for command-specific help.")
}
