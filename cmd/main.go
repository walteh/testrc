package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/walteh/snake"
	"github.com/walteh/testrc/cmd/root"
)

func main() {

	ctx := context.Background()

	rootCmd := snake.NewRootCommand(ctx, &root.Root{})

	if err := snake.DecorateRootCommand(ctx, rootCmd, &snake.DecorateOptions{
		Headings: color.New(color.FgCyan, color.Bold),
		ExecName: color.New(color.FgHiGreen, color.Bold),
		Commands: color.New(color.FgHiRed, color.Faint),
	}); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "[%s] (error) %+v\n", rootCmd.Name(), err)
		if err != nil {
			panic(err)
		}
	}

	rootCmd.SilenceErrors = true

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "[%s] (error) %+v\n", rootCmd.Name(), err)
		if err != nil {
			panic(err)
		}
	}

}
