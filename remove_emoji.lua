function Link(el)
    local new_content = {}
    for _, item in ipairs(el.content) do
        if item.t == "Str" then
            local new_text = item.text:gsub("([%s\xC2-\xF4][\x80-\xBF]*)", function(c)
                if c:match("^[\xC2-\xF4][\x80-\xBF]+$") then
                    local code_point = utf8_to_codepoint(c)
                    if is_emoji(code_point) then
                        return ""
                    else
                        return c
                    end
                else
                    return c
                end
            end)
            new_text = new_text:gsub("%s+", " "):gsub("^%s", ""):gsub("%s$", "")
            if new_text ~= "" then
                table.insert(new_content, pandoc.Str(new_text))
            end
        elseif item.t == "Space" then
            if #new_content > 0 and new_content[#new_content].t ~= "Space" then
                table.insert(new_content, item)
            end
        else
            table.insert(new_content, item)
        end
    end
    el.content = new_content
    return el
end

function Para(el)
    local new_content = {}
    for _, item in ipairs(el.content) do
        if item.t == "Str" then
            local new_text = item.text:gsub("([%s\xC2-\xF4][\x80-\xBF]*)", function(c)
                if c:match("^[\xC2-\xF4][\x80-\xBF]+$") then
                    local code_point = utf8_to_codepoint(c)
                    if is_emoji(code_point) then
                        return c
                    end
                end
                return c
            end)
            table.insert(new_content, pandoc.Str(new_text))
        else
            table.insert(new_content, item)
        end
    end
    el.content = new_content
    return el
end

function is_emoji(code_point)
    return (
        (code_point >= 0x1F600 and code_point <= 0x1F64F) or -- Emoticons
        (code_point >= 0x1F300 and code_point <= 0x1F5FF) or -- Miscellaneous Symbols and Pictographs
        (code_point >= 0x1F680 and code_point <= 0x1F6FF) or -- Transport and Map Symbols
        (code_point >= 0x1F900 and code_point <= 0x1F9FF) or -- Symbols and Pictographs Extended-A
        (code_point >= 0x1FA70 and code_point <= 0x1FAFF) or -- Symbols and Pictographs Extended-A
        (code_point == 0x200D) -- Zero Width Joiner
    )
end

function utf8_to_codepoint(utf8_char)
    local byte1, byte2, byte3, byte4 = utf8_char:byte(1, 4)
    if byte1 < 224 then
        return (byte1 - 192) * 64 + (byte2 - 128)
    elseif byte1 < 240 then
        return (byte1 - 224) * 4096 + (byte2 - 128) * 64 + (byte3 - 128)
    else
        return (byte1 - 240) * 262144 + (byte2 - 128) * 4096 + (byte3 - 128) * 64 + (byte4 - 128)
    end
end
