Usage:

``` bash
go test ./...
```


Example of a call

``` bash

pandoc --wrap=none --from=gfm+wikilinks_title_after_pipe --to=gfm+wikilinks_title_after_pipe --lua-filter=remove_emphasis.lua --lua-filter=remove_emoji.lua --lua-filter=trim_link_names.lua --lua-filter=link.lua

```