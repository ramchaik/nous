def create_retriever(vectorstore):
    return vectorstore.as_retriever(k=3, search_type="mmr")
