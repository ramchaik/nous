import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# API Keys
TAVILY_API_KEY = os.getenv("TAVILY_API_KEY")

# Vectorstore paths
DATA_INDEX_PATH = os.getenv('DATA_INDEX_PATH', "data/faiss_index.bin")
DATA_VECTOR_PATH = os.getenv('DATA_VECTOR_PATH', "data/vectors.npy")
DATA_INDEX_TO_ID_PATH = os.getenv('DATA_INDEX_TO_ID_PATH', "data/index_to_id.npy")

# Ollama configuration
class OllamaConfig:
    HOST = os.getenv('OLLAMA_HOST', 'http://localhost:11434')
    EMBEDDINGS_MODEL = os.getenv('OLLAMA_EMBEDDINGS_MODEL', "nomic-embed-text")
    MODEL_RAG = os.getenv('OLLAMA_MODEL_RAG', "phi3:mini")
    MODEL_RETRIEVER_GRADER = os.getenv('OLLAMA_MODEL_RETRIEVER_GRADER', "thinyllama:mini")

# Personal Knowledge base resources
KNOWLEDGE_BASE_URLS = [
    "https://tip.golang.org/tour/concurrency.article",
    "https://tip.golang.org/doc/effective_go",
    "https://gosafir.com/mag/wp-content/uploads/2019/12/Tolkien-J.-The-lord-of-the-rings-HarperCollins-ebooks-2010.pdf",
    "https://gist.github.com/silver-xu/1dcceaa14c4f0253d9637d4811948437",
]

# Ensure all required environment variables are set
def validate_env_vars():
    required_vars = ["TAVILY_API_KEY"]
    for var in required_vars:
        if not os.getenv(var):
            raise EnvironmentError(f"Missing required environment variable: {var}")

validate_env_vars()
