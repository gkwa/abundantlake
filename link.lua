function Link(el)
    local url = el.target
    local text = pandoc.utils.stringify(el.content)

    -- Check if it's a wikilink
    if el.title == "wikilink" then
      -- If it's a wikilink, return it unchanged
      return el
    end

    -- Check if the link is already formatted (i.e., text != url)
    if text ~= url then
      -- If it's already formatted, return it unchanged
      return el
    else
      -- If it's a bare URL, leave it alone
      return pandoc.RawInline('markdown', url)
    end
  end
