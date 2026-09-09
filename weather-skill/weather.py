# Weather Skill Source (Ported from clawhub.ai/steipete/skills/weather)
# Provides real-time weather and forecast using wttr.in (primary) and Open-Meteo (fallback)

import sys
import urllib.request
import json

def get_weather(location="New+York", format_type="3"):
    """
    Fetch weather from wttr.in or Open-Meteo.
    """
    url = f"https://wttr.in/{location}?format={format_type}"
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "curl/7.68.0"})
        with urllib.request.urlopen(req, timeout=10) as resp:
            return resp.read().decode('utf-8').strip()
    except Exception as e:
        return f"Error querying wttr.in: {e}"

if __name__ == "__main__":
    loc = sys.argv[1] if len(sys.argv) > 1 else "New+York"
    fmt = sys.argv[2] if len(sys.argv) > 2 else "3"
    print(get_weather(loc, fmt))
