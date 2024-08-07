function Link(el)
    local new_content = {}
    for _, item in ipairs(el.content) do
        if item.t == "Str" then
            table.insert(new_content, pandoc.Str(item.text))
        elseif item.t == "Emph" or item.t == "Strong" then
            for _, subitem in ipairs(item.content) do
                if subitem.t == "Str" then
                    table.insert(new_content, pandoc.Str(subitem.text))
                else
                    table.insert(new_content, subitem)
                end
            end
        else
            table.insert(new_content, item)
        end
    end

    -- Combine all content into a single string
    local combined_text = ""
    for _, item in ipairs(new_content) do
        if item.t == "Str" then
            combined_text = combined_text .. item.text
        elseif item.t == "Space" then
            combined_text = combined_text .. " "
        end
    end

    -- Trim the combined text
    combined_text = combined_text:match("^%s*(.-)%s*$")

    -- Replace the content with the trimmed text
    el.content = { pandoc.Str(combined_text) }

    return el
end

