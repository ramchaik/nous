from langchain_ollama import ChatOllama
from functools import lru_cache

@lru_cache(maxsize=2)
def create_llm(model="llama3.1", temperature=0, format=''):
    return ChatOllama(model=model, temperature=temperature, format=format)