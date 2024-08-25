from langchain_ollama import ChatOllama

def create_llm(model="llama3.1", temperature=0, format=''):
    return ChatOllama(model=model, temperature=temperature, format=format)