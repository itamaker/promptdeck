package app

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runTUI()
	}

	switch args[0] {
	case "render":
		return runRender(args[1:])
	case "matrix":
		return runMatrix(args[1:])
	case "optimize":
		return runOptimize(args[1:])
	case "tui", "interactive":
		return runTUI()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Println("promptdeck renders reusable prompt templates and optimizes prompt experiments.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  promptdeck                 # launch Bubble Tea TUI")
	fmt.Println("  promptdeck render -template examples/review.tmpl -vars examples/vars.json")
	fmt.Println("  promptdeck matrix -template examples/review.tmpl -matrix examples/matrix.json")
	fmt.Println("  promptdeck optimize -template examples/review.tmpl -matrix examples/matrix.json -scores examples/scores.json")
}
