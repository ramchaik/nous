from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.document_loaders import WebBaseLoader
from concurrent.futures import ThreadPoolExecutor

def load_documents(urls):
    with ThreadPoolExecutor() as executor:
        docs = list(executor.map(lambda url: WebBaseLoader(url).load(), urls))
    return [item for sublist in docs for item in sublist]

def load_and_split_documents(urls):
    print("Loading and splitting documents")
    docs_list = load_documents(urls)
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=50)
    return text_splitter.split_documents(docs_list)
