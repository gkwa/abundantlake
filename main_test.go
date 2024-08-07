package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestPandocFilters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "RemoveEmphasisFromLink",
			input: `
			
			I like [*Google* Search](https://google.com)
			
			`,
			expected: `
			
			I like [Google Search](https://google.com)
			
			`,
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
			name: "Remove emojis from link text only",
			input: `
      
      Check out this 😀 [😀 Awesome 🎉 Link 🌟](https://example.com) 🎉
      
      `,
			expected: `
      
      Check out this 😀 [Awesome Link](https://example.com) 🎉
      
      `,
		},

		{
			name: "Preserve emojis outside links",
			input: `

This is a 🌟 test.  [They call this drug eggs in Korea because these are so addictive!! Pt. 2 😳🥚🔥 - YouTube](https://www.youtube.com/shorts/MBnJsEbDflA)

      `,
			expected: `

This is a 🌟 test. [They call this drug eggs in Korea because these are so addictive!! Pt. 2 - YouTube](https://www.youtube.com/shorts/MBnJsEbDflA)

`,
		},
		{
			name: "Preserve emphasis outside links and clear within links",
			input: `

      This is an __italics__ test with [Some __italics__ Link](https://example.com)

      `,
			expected: `

 		This is an **italics** test with [Some italics Link](https://example.com)

`,
		},
		{
			name: "Trim spaces in link names",
			input: `
      [ Spaced  Link ](https://example.com)
      `,
			expected: `
      [Spaced Link](https://example.com)
      `,
		},
		{
			name: "Obsidian links",
			input: `

[[test]]

[[test#asdf]]

[[xxx^section|alias]]			

`,
			expected: `

[[test]]

[[test#asdf]]

[[xxx^section|alias]]			



`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := applyPandocFilter(tt.input, "remove_emphasis.lua")
			if err != nil {
				t.Fatalf("failed to apply Pandoc filter for test '%s': %v", tt.name, err)
			}
			output, err = applyPandocFilter(output, "remove_emoji.lua")
			if err != nil {
				t.Fatalf("failed to apply Pandoc filter for test '%s': %v", tt.name, err)
			}
			output, err = applyPandocFilter(output, "trim_link_names.lua")
			if err != nil {
				t.Fatalf("failed to apply Pandoc filter for test '%s': %v", tt.name, err)
			}

			output = strings.ReplaceAll(output, "\\|", "|")

			if output != strings.TrimSpace(tt.expected) {
				diff := generateDiff(tt.input, tt.expected, output)
				t.Errorf("test '%s' failed.\nDiff:\n%s", tt.name, diff)
			}
		})
	}
}

func applyPandocFilter(input, filterPath string) (string, error) {
	cmd := exec.Command(
		"pandoc",
		"--wrap=none",
		"--from=gfm+wikilinks_title_after_pipe",
		"--to=gfm+wikilinks_title_after_pipe",
		"--lua-filter="+filterPath,
	)

	cmd.Stdin = strings.NewReader(strings.TrimSpace(input))
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
	return s
}
