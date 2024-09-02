from langchain_community.vectorstores import FAISS
from langchain_community.embeddings import OllamaEmbeddings
from langchain_community.docstore.in_memory import InMemoryDocstore
import faiss
from llm_api.config import DATA_INDEX_PATH, DATA_INDEX_TO_ID_PATH, DATA_VECTOR_PATH, OllamaConfig
import numpy as np
import os
from typing import TypedDict, List

def create_vectorstore(documents: List[str], force_rebuild: bool = False) -> FAISS:
    index_path = DATA_INDEX_PATH
    vectors_path = DATA_VECTOR_PATH
    index_to_id_path = DATA_INDEX_TO_ID_PATH

    if not force_rebuild and os.path.exists(index_path) and os.path.exists(vectors_path) and os.path.exists(index_to_id_path):
        # Load existing index and vectors
        index = faiss.read_index(index_path)
        docstore_dict = np.load(vectors_path, allow_pickle=True).item()
        index_to_id = np.load(index_to_id_path, allow_pickle=True).item()

        # Create InMemoryDocstore from the loaded dictionary
        docstore = InMemoryDocstore(docstore_dict)

        ollama_embeddings = OllamaEmbeddings(model=OllamaConfig.EMBEDDINGS_MODEL, base_url=OllamaConfig.HOST)
        vectorstore = FAISS(
            embedding_function=ollama_embeddings,
            index=index,
            docstore=docstore,
            index_to_docstore_id=index_to_id
        )
        print("Loaded existing vectorstore from disk.")
    else:
        # Create new vectorstore
        ollama_embeddings = OllamaEmbeddings(model=OllamaConfig.EMBEDDINGS_MODEL, base_url=OllamaConfig.HOST)
        vectorstore = FAISS.from_documents(documents, ollama_embeddings)

        # Save the index, vectors, and index_to_docstore_id separately
        faiss.write_index(vectorstore.index, index_path)
        np.save(vectors_path, vectorstore.docstore._dict)
        np.save(index_to_id_path, vectorstore.index_to_docstore_id)
        print("Created new vectorstore and saved to disk.")

    return vectorstore
