function Link(el)
  local new_content = {}
  local combined_text = ""

  for _, item in ipairs(el.content) do
    if item.t == "Str" then
      combined_text = combined_text .. item.text
    elseif item.t == "Space" then
      combined_text = combined_text .. " "
    end
  end

  combined_text = combined_text:match("^%s*(.-)%s*$")

  el.content = { pandoc.Str(combined_text) }
  return el
end
