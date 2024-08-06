package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "RemoveEmphasisFromLink",
			input:    strings.TrimSpace(`I like [*Google* Search](https://google.com)`),
			expected: `I like [Google Search](https://google.com)`,
		},
		{
			name:     "Multiple spaces are squeezed",
			input:    `I like ice  cream`,
			expected: `I like ice cream`,
		},
		{
			name: "RemoveEmphasisFromLink 2",
			input: `
I like [**Google** Search](https://google.com)`,
			expected: `
I like [Google Search](https://google.com)
`,
		},
		{
			name: "Bulleted list alongside some links",
			input: `
lorum ipsum lorum ipsum lorum ipsum lorum ipsum lorum ipsum

I like [**Google** Search](https://google.com)
`,
			expected: `
lorum ipsum lorum ipsum lorum ipsum lorum ipsum lorum ipsum

I like [Google Search](https://google.com)
`,
		},
		{
			name: "Links with aliases",
			input: `

         [a | b](https://a.com/b)

         `,
			expected: `

         [a | b](https://a.com/b)

         
         `,
		},
		{
			name: "Links with aliases and emphasis on name",
			input: `

         [**a** | b](https://a.com/b)

         `,
			expected: `

         [a | b](https://a.com/b)

         
         `,
		},
		{
			name: "Links with aliases and emphasis on name and alias",
			input: `

         [**a** | **b**](https://a.com/b)

         `,
			expected: `

         [a | b](https://a.com/b)

         
         `,
		},
		{
			name: "Italics are removed from name",
			input: `

         [Google __Search__](https://google.com)

         `,
			expected: `

         [Google Search](https://google.com)

         
         `,
		},
		{
			name: "Italics are removed from alias",
			input: `

         [Google Search|__GSearch__](https://google.com)

         `,
			expected: `

         [Google Search|GSearch](https://google.com)

         
         `,
		},

      {
			name: "Italics are removed from name and alias",
			input: `

         [Google __Search__|__GSearch__](https://google.com)

         `,
			expected: `

         [Google Search|GSearch](https://google.com)

         
         `,
		},

      {
			name:     "Remove emojis from link text only",
			input:    `
         
         Check out this 😀 [😀 Awesome 🎉 Link 🌟](https://example.com) 🎉
         
         
         `,
			expected: `
         
         Check out this 😀 [Awesome Link](https://example.com) 🎉
         
         `,
		},
		{
			name:     "Preserve emojis outside links",
			input:    `

         This is a 🌟 test with [Some 🎉 Link](https://example.com) and more 🎈 emojis

         `,
			expected: `

         This is a 🌟 test with [Some Link](https://example.com) and more 🎈 emojis

`,
		},
	}
	for _, tt := range tests {
		output, err := applyPandocFilter(strings.TrimSpace(tt.input), "remove_emphasis.lua")
		if err != nil {
			return fmt.Errorf("failed to apply Pandoc filter for test '%s': %v", tt.name, err)
		}
		if output != strings.TrimSpace(tt.expected) {
			diff := generateDiff(strings.TrimSpace(tt.input), strings.TrimSpace(tt.expected), output)
			return fmt.Errorf("test '%s' failed.\nDiff:\n%s", tt.name, diff)
		}
		fmt.Printf("Test '%s' passed.\n", tt.name)
	}
	return nil
}

func applyPandocFilter(input, filterPath string) (string, error) {
	cmd := exec.Command("pandoc", "--from=gfm", "--to=gfm", "--lua-filter="+filterPath)
	cmd.Stdin = strings.NewReader(input)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pandoc command failed: %v\nStderr: %s", err, errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateDiff(input, expected, actual string) string {
	dmp := diffmatchpatch.New()
	diffs1 := dmp.DiffMain(input, expected, false)
	diffs2 := dmp.DiffMain(input, actual, false)
	var result strings.Builder
	result.WriteString("Input:    ")
	result.WriteString(visualizeInvisibles(dmp.DiffText1(diffs1)))
	result.WriteString("\nExpected: ")
	result.WriteString(visualizeInvisibles(dmp.DiffText2(diffs1)))
	result.WriteString("\nActual:   ")
	result.WriteString(visualizeInvisibles(dmp.DiffText2(diffs2)))
	return result.String()
}

func visualizeInvisibles(s string) string {
	s = strings.ReplaceAll(s, " ", "·")
	s = strings.ReplaceAll(s, "\n", "↵\n")
	return s
}
