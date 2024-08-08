#!/bin/bash
pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_title_after_pipe -t native input1.md
pandoc --lua-filter=standardize_links.lua --wrap=none --from=gfm+wikilinks_title_after_pipe -t native input2.md
