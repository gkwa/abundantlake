function Link(el)
 local url = el.target
 local text = pandoc.utils.stringify(el.content)
 if text == url then
   text = url
 end
 return pandoc.Link(text, url)
end

