package main

import (
	"bytes"
	"fmt"
	"log"
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

[The BEST ways to make Gai Lan👌 + 2 Quick & Easy Stir Fry Chinese Broccoli Recipes - YouTube](https://www.youtube.com/watch?v=GQ7-tjp7wnA)

      `,
			expected: `
      
Check out this 😀 [Awesome Link](https://example.com) 🎉

[The BEST ways to make Gai Lan + 2 Quick & Easy Stir Fry Chinese Broccoli Recipes - YouTube](https://www.youtube.com/watch?v=GQ7-tjp7wnA)

`,
		},
		{
			name: "Remove multiple emojis from link text only",
			input: `
      
      Check out this 😀 [😀 Awesome Link 🎉 and things 😀](https://example.com) 🎉

      `,
			expected: `
      
Check out this 😀 [Awesome Link and things](https://example.com) 🎉
      
      `,
		},
		{
			name: "Remove multiple emojis from link text only 2",
			input: `
      
[Testing in Go: Golden Files · Ilija Eftimov 👨‍🚀](https://ieftimov.com)

      `,
			expected: `
      
[Testing in Go: Golden Files · Ilija Eftimov](https://ieftimov.com)

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
			name: "Typical markdown notes test",
			input: `

[[leadyspleen]]

[Character encoding: iconv -t utf-8 input.txt | pandoc | iconv -f utf-8](https://pandoc.org/chunkedhtml-demo/2.3-character-encoding.html)

run container without docker on macos

[Stream: Go 10 Week Backend Eng Onboarding](https://stream-wiki.notion.site/Stream-Go-10-Week-Backend-Eng-Onboarding-625363c8c3684753b7f2b7d829bcd67a)

[What is Bun?](https://bun.uptrace.dev/guide/)

[Go errors.Is now includes a nil check. | by the korean guy | Jul, 2024 | Medium](https://medium.com/@ojh031/go-errors-is-now-includes-a-nil-check-94f82fe4cc31)

[darrenburns/posting: The modern API client that lives in your terminal](https://github.com/darrenburns/posting?tab=readme-ov-file#posting)

[NewStore TechTalk - Advanced Testing with Go by Mitchell Hashimoto - YouTube](https://www.youtube.com/watch?v=yszygk1cpEc)

https://www.google.com/search?q=how+to+make+refried+beans

https://www.google.com/search?q=how+to+make+vegitarian+refried+beans

https://www.google.com/search?q=how+to+make+vegitarian+refried+beans

      `,
			expected: `

[[leadyspleen]]

[Character encoding: iconv -t utf-8 input.txt | pandoc | iconv -f utf-8](https://pandoc.org/chunkedhtml-demo/2.3-character-encoding.html)

run container without docker on macos

[Stream: Go 10 Week Backend Eng Onboarding](https://stream-wiki.notion.site/Stream-Go-10-Week-Backend-Eng-Onboarding-625363c8c3684753b7f2b7d829bcd67a)

[What is Bun?](https://bun.uptrace.dev/guide/)

[Go errors.Is now includes a nil check. | by the korean guy | Jul, 2024 | Medium](https://medium.com/@ojh031/go-errors-is-now-includes-a-nil-check-94f82fe4cc31)

[darrenburns/posting: The modern API client that lives in your terminal](https://github.com/darrenburns/posting?tab=readme-ov-file#posting)

[NewStore TechTalk - Advanced Testing with Go by Mitchell Hashimoto - YouTube](https://www.youtube.com/watch?v=yszygk1cpEc)

https://www.google.com/search?q=how+to+make+refried+beans

https://www.google.com/search?q=how+to+make+vegitarian+refried+beans

https://www.google.com/search?q=how+to+make+vegitarian+refried+beans

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
			name: "Bare link is output as is without angle brackets like <https://a.com>",
			input: `

https://example.com

`,
			expected: `

https://example.com

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
		{
			name: "Code block",
			input: `

Usage:

` + "```" + `bash
git clone https://github.com/gkwa/abundantlake.git
cd abundantlake
go test ./...
` + "```" + `

`,
			expected: `

Usage:

` + "``` " + `bash
git clone https://github.com/gkwa/abundantlake.git
cd abundantlake
go test ./...
` + "```" + `

`,
		},
		{
			name: "Bare link is left alone inside code block",
			input: `

Usage:

` + "```" + ` bash
git clone https://github.com/gkwa/abundantlake.git
` + "```" + `

`,
			expected: `

Usage:

` + "```" + ` bash
git clone https://github.com/gkwa/abundantlake.git
` + "```" + `

`,
		},
		{
			name: "Code block with link",
			input: `

[readme | **the** information source](https://readme.com)

Usage:

` + "```" + `bash
git clone https://github.com/gkwa/abundantlake.git
cd abundantlake
go test ./...
` + "```" + `

`,
			expected: `

[readme | the information source](https://readme.com)

Usage:

` + "``` " + `bash
git clone https://github.com/gkwa/abundantlake.git
cd abundantlake
go test ./...
` + "```" + `

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := applyPandocFilters(tt.input, []string{
				"remove_emphasis.lua",
				"remove_emoji.lua",
				"trim_link_names.lua",
				"link.lua",
			})
			if err != nil {
				t.Fatalf("failed to apply Pandoc filters for test '%s': %v", tt.name, err)
			}

			output = strings.ReplaceAll(output, "\\|", "|")

			if output != strings.TrimSpace(tt.expected) {
				diff := generateDiff(tt.input, tt.expected, output)
				t.Errorf("test '%s' failed.\nDiff:\n%s", tt.name, diff)
			}
		})
	}
}

func applyPandocFilters(input string, filterPaths []string) (string, error) {
	args := []string{
		"--wrap=none",
		"--from=gfm+wikilinks_title_after_pipe",
		"--to=gfm+wikilinks_title_after_pipe",
	}

	for _, filterPath := range filterPaths {
		args = append(args, "--lua-filter="+filterPath)
	}

	cmd := exec.Command("pandoc", args...)

	// Log the Pandoc CLI command being run
	log.Printf("Running Pandoc command: pandoc %s", strings.Join(args, " "))

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
	result.WriteString("Input:\n")
	result.WriteString(visualizeInvisibles(dmp.DiffText1(diffs1)))
	result.WriteString("\nExpected:\n")
	result.WriteString(visualizeInvisibles(dmp.DiffText2(diffs1)))
	result.WriteString("\nActual:\n")
	result.WriteString(visualizeInvisibles(dmp.DiffText2(diffs2)))
	return result.String()
}

func visualizeInvisibles(s string) string {
	s = strings.ReplaceAll(s, " ", "·")
	return s
}
