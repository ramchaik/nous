from langchain_community.vectorstores import FAISS
from langchain_community.embeddings import OllamaEmbeddings
from langchain_community.docstore.in_memory import InMemoryDocstore
import faiss
import numpy as np
import os
from typing import TypedDict, List

def create_vectorstore(documents: List[str], force_rebuild: bool = False) -> FAISS:
    index_path = "data/faiss_index.bin"
    vectors_path = "data/vectors.npy"
    index_to_id_path = "data/index_to_id.npy"
    
    if not force_rebuild and os.path.exists(index_path) and os.path.exists(vectors_path) and os.path.exists(index_to_id_path):
        # Load existing index and vectors
        index = faiss.read_index(index_path)
        docstore_dict = np.load(vectors_path, allow_pickle=True).item()
        index_to_id = np.load(index_to_id_path, allow_pickle=True).item()
        
        # Create InMemoryDocstore from the loaded dictionary
        docstore = InMemoryDocstore(docstore_dict)
        
        ollama_host = os.getenv('OLLAMA_HOST', 'http://localhost:11434')
        ollama_embeddings = OllamaEmbeddings(model="nomic-embed-text", base_url=ollama_host)
        vectorstore = FAISS(
            embedding_function=ollama_embeddings,
            index=index,
            docstore=docstore,
            index_to_docstore_id=index_to_id
        )
        print("Loaded existing vectorstore from disk.")
    else:
        # Create new vectorstore
        ollama_host = os.getenv('OLLAMA_HOST', 'http://localhost:11434')
        ollama_embeddings = OllamaEmbeddings(model="nomic-embed-text", base_url=ollama_host)
        vectorstore = FAISS.from_documents(documents, ollama_embeddings)
        
        # Save the index, vectors, and index_to_docstore_id separately
        faiss.write_index(vectorstore.index, index_path)
        np.save(vectors_path, vectorstore.docstore._dict)
        np.save(index_to_id_path, vectorstore.index_to_docstore_id)
        print("Created new vectorstore and saved to disk.")
    
    return vectorstore