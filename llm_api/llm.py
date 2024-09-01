import os
from langchain_ollama import ChatOllama
from functools import lru_cache

@lru_cache(maxsize=2)
def create_llm(model="llama3.1", temperature=0, format=''):
    ollama_host = os.getenv('OLLAMA_HOST', 'http://localhost:11434')
    return ChatOllama(
        model=model, 
        temperature=temperature, 
        format=format,
        base_url=ollama_host
    )