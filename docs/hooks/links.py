import re

def on_page_markdown(markdown: str, **kwargs):
    links = [el[0] for el in re.findall(r"((?<!!)\[((X[^\]]+?)|([^X].*?))\]\(http[s]?://.*?\))", markdown)]
    content = markdown
    for link in links:
        content = content.replace(link, link + '{: rel="nofollow noopener" target="_blank" }')
    return content
