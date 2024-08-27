import os
from dotenv import load_dotenv

load_dotenv()

TAVILY_API_KEY = os.getenv("TAVILY_API_KEY")

URLS = [
    "https://tip.golang.org/tour/concurrency.article",
    "https://tip.golang.org/doc/effective_go",
    # "https://socrates.acadiau.ca/courses/engl/rcunningham/resources/Shpe/Hamlet.pdf",
    # "https://gosafir.com/mag/wp-content/uploads/2019/12/Tolkien-J.-The-lord-of-the-rings-HarperCollins-ebooks-2010.pdf"
]