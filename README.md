# Abundantlake

Abundantlake is a Go-based tool that applies a series of Pandoc filters to Markdown content, primarily focusing on link and text formatting. It's designed to clean up and standardize Markdown documents, especially those containing various types of links and formatting.

## Features

- Remove emphasis (bold and italic) from link text
- Remove emojis from link text while preserving them in regular text
- Trim spaces in link names
- Preserve Obsidian-style wiki links
- Handle bare URLs and code blocks appropriately

## Usage

To use Abundantlake, you need to have Go and Pandoc installed on your system.

1. Clone the repository:
   ```bash
   git clone https://github.com/gkwa/abundantlake.git
   cd abundantlake
   ```

2. Run the tests to ensure everything is working correctly:
   ```bash
   go test ./...
   ```

3. To process your own Markdown files, you can use the following Pandoc command:
   ```bash
   pandoc \
     --wrap=none \
     --from=gfm+wikilinks_title_after_pipe \
     --to=gfm+wikilinks_title_after_pipe \
     --lua-filter=remove_emphasis.lua \
     --lua-filter=remove_emoji.lua \
     --lua-filter=trim_link_names.lua \
     --lua-filter=link.lua \
     input.md -o output.md
   ```

   Replace `input.md` with your input file and `output.md` with your desired output file.

## Lua Filters

The project includes several Lua filters for Pandoc:

- `remove_emphasis.lua`: Removes bold and italic formatting from link text
- `remove_emoji.lua`: Removes emojis from link text while preserving them in regular text
- `trim_link_names.lua`: Trims spaces from link names
- `link.lua`: Handles various link formatting tasks

## Testing

The project includes a comprehensive test suite. To run the tests, use:

```bash
go test ./...
```

For verbose output, use:

```bash
go test ./... -v
```
