from langchain_community.vectorstores import FAISS
from langchain_community.embeddings import OllamaEmbeddings

def create_vectorstore(documents):
    ollama_embeddings = OllamaEmbeddings(model="nomic-embed-text")
    vectorstore = FAISS.from_documents(documents, ollama_embeddings)
    return vectorstore