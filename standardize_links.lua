function Link(el)
  local url = el.target
  local text = pandoc.utils.stringify(el.content)
  if text == "" then
    text = url
  end
  return pandoc.RawInline('markdown', "[" .. text .. "](" .. url .. ")")
end

