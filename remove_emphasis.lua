REMOVE_EMOJIS = true

local function is_emoji(code_point)
    return (
        (code_point >= 0x1F600 and code_point <= 0x1F64F) or -- Emoticons
        (code_point >= 0x1F300 and code_point <= 0x1F5FF) or -- Miscellaneous Symbols and Pictographs
        (code_point >= 0x1F680 and code_point <= 0x1F6FF) or -- Transport and Map Symbols
        (code_point >= 0x1F900 and code_point <= 0x1F9FF) or -- Supplemental Symbols and Pictographs
        (code_point >= 0x1FA70 and code_point <= 0x1FAFF)    -- Symbols and Pictographs Extended-A
    )
end

local function utf8_to_codepoint(utf8_char)
    local byte1, byte2, byte3, byte4 = utf8_char:byte(1, 4)
    if byte1 < 128 then
        return byte1
    elseif byte1 < 224 then
        return (byte1 - 192) * 64 + (byte2 - 128)
    elseif byte1 < 240 then
        return (byte1 - 224) * 4096 + (byte2 - 128) * 64 + (byte3 - 128)
    else
        return (byte1 - 240) * 262144 + (byte2 - 128) * 4096 + (byte3 - 128) * 64 + (byte4 - 128)
    end
end

local function is_utf8_emoji(utf8_char)
    local code_point = utf8_to_codepoint(utf8_char)
    return is_emoji(code_point)
end

local function remove_emojis(s)
    if not REMOVE_EMOJIS then
        return s
    end
    return s:gsub("[\xC2-\xF4][\x80-\xBF]+", function(c)
        if is_utf8_emoji(c) then
            return ""
        else
            return c
        end
    end)
end

function Link(el)
    local new_content = {}
    for _, item in ipairs(el.content) do
        if item.t == "Str" then
            table.insert(new_content, pandoc.Str(remove_emojis(item.text)))
        elseif item.t == "Emph" or item.t == "Strong" then
            for _, subitem in ipairs(item.content) do
                if subitem.t == "Str" then
                    table.insert(new_content, pandoc.Str(remove_emojis(subitem.text)))
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
