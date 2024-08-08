#!/bin/bash
pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_title_after_pipe --to=gfm+wikilinks_title_after_pipe input1.md
pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_
