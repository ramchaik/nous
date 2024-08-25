from langchain_community.vectorstores import SKLearnVectorStore
from langchain_community.embeddings import OllamaEmbeddings

def create_vectorstore(documents):
    ollama_embeddings = OllamaEmbeddings(model="nomic-embed-text")
    vectorstore = SKLearnVectorStore.from_documents(
        documents=documents,
        embedding=ollama_embeddings,
    )
    return vectorstore