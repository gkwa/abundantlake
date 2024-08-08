#!/bin/bash
echo "https://google.com" | pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_title_after_pipe -t native
echo '[aaaaaaaaaaaaa](bbbbbbbbbbbbbb)' | pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_title_after_pipe -t native

