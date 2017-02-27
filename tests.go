package main

import (
	"fmt"

	"bytes"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/reconquest/loreley"
)

func main2() {
	buffer := &bytes.Buffer{}

	writer := tabwriter.NewWriter(buffer, 2, 4, 2, ' ', tabwriter.FilterHTML)

	writer.Write([]byte(strings.Join(
		[]string{
			"<underline>CORES<reset>",
			"<underline>DESCRIPTION<reset>\n",
		}, "\t",
	)))

	writer.Write([]byte(strings.Join(
		[]string{
			"<fg 15><bg 1><bold> 1 <reset> <fg 15><bg 243><bold> 3 <reset>",
			"test\n",
		}, "\t",
	)))

	writer.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	result, err := loreley.CompileAndExecuteToString(
		buffer.String(),
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Print(result)
}

func test() {
	// Use handy standard colors
	color.Set(color.FgYellow)

	fmt.Println("Existing text will now be in yellow")
	fmt.Printf("This one %s\n", "too")

	color.Unset() // Don't forget to unset

	// You can mix up parameters
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset() // Use it in your function

	fmt.Println("All text will now be bold magenta.")

	// Print with default helper functions
	color.Cyan("Prints text in cyan.")

	// A newline will be appended automatically
	color.Blue("Prints %s in blue.", "text")

	// These are using the default foreground colors
	color.Red("We have red")
	color.Magenta("And many others ..")

	// Create a new color object
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")

	// Or just add them to New()
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("This prints bold cyan %s\n", "too!.")

	// Mix up foreground and background colors, create new mixes!
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")

	whiteBackground := red.Add(color.BgWhite)
	whiteBackground.Println("Red text with white background.")

	fmt.Println("This", color.RedString("warning"), "should be not neglected.")
}

/*
func main() {
	main2()
}
*/
