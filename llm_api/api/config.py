import os
from dotenv import load_dotenv

load_dotenv()

TAVILY_API_KEY = os.getenv("TAVILY_API_KEY")

URLS = [
    "https://tip.golang.org/tour/concurrency.article",
    "https://tip.golang.org/doc/effective_go"
]