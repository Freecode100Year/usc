# Web Search Extended Skill (Ported from clawhub.ai/yejinlei/skills/web-search-ex-skill)
# Multi-engine web search (Baidu, Bing, DuckDuckGo) without API keys.

import sys
import json
import urllib.request
import urllib.parse
import re

def search_duckduckgo(query, num_results=5):
    """Fallback search using DuckDuckGo HTML scraping"""
    url = f"https://html.duckduckgo.com/html/?q={urllib.parse.quote(query)}"
    req = urllib.request.Request(url, headers={
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    })
    results = []
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            html = resp.read().decode('utf-8', errors='ignore')
            links = re.findall(r'<a[^>]+class="result__snippet[^>]+href="([^"]+)"[^>]*>(.*?)</a>', html)
            titles = re.findall(r'<a[^>]+class="result__url[^>]*>(.*?)</a>', html)
            for i in range(min(num_results, len(links))):
                results.append({
                    "title": re.sub(r'<[^>]+>', '', titles[i]) if i < len(titles) else "Result",
                    "url": links[i][0],
                    "snippet": re.sub(r'<[^>]+>', '', links[i][1])
                })
    except Exception as e:
        pass
    return results

if __name__ == "__main__":
    q = sys.argv[1] if len(sys.argv) > 1 else "New York free online courses"
    print(json.dumps({"query": q, "results": search_duckduckgo(q)}, ensure_ascii=False))
