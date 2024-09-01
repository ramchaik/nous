from typing import TypedDict, List
from langchain.schema import Document
from langchain_core.output_parsers import StrOutputParser, JsonOutputParser
from langgraph.graph import START, END, StateGraph
from functools import lru_cache
import time
import logging
import concurrent.futures

from config import TAVILY_API_KEY, URLS
from document_processor import load_and_split_documents
from embeddings import create_vectorstore
from retriever import create_retriever

from llm import create_llm
from prompts import rag_prompt, grading_prompt

from langchain_community.tools.tavily_search import TavilySearchResults

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

doc_splits = load_and_split_documents(URLS)
vectorstore = create_vectorstore(doc_splits, force_rebuild=False)
retriever = create_retriever(vectorstore)

llm = create_llm(model="phi3:mini")
rag_chain = rag_prompt | llm | StrOutputParser()
retrieval_grader = grading_prompt | create_llm(model="tinyllama", format="json") | JsonOutputParser()

web_search_tool = TavilySearchResults(api_key=TAVILY_API_KEY)

class GraphState(TypedDict):
    question: str
    generation: str
    search: str
    documents: List[str]
    steps: List[str]

@lru_cache(maxsize=100)
def cached_web_search(question: str):
    logger.info("Web search cache miss")
    return web_search_tool.invoke({"query": question})

def retrieve(state: GraphState) -> GraphState:
    start_time = time.time()
    question = state["question"]
    documents = retriever.invoke(question)
    state["documents"] = documents
    state["steps"].append("retrieve_documents")
    logger.info(f"Retrieve time: {time.time() - start_time:.2f} seconds")
    return state

@lru_cache(maxsize=100)
def cached_rag_chain_invoke(documents, question):
    return rag_chain.invoke({"documents": documents, "question": question})

def generate(state: GraphState) -> GraphState:
    start_time = time.time()
    docs_content = "\n".join(doc.page_content for doc in state["documents"][:3])  # Limit to top 3 documents
    state["generation"] = rag_chain.invoke({"documents": docs_content, "question": state["question"]})
    state["steps"].append("generate_answer")
    logger.info(f"Generate time: {time.time() - start_time:.2f} seconds")
    return state

def grade_documents(state: GraphState) -> GraphState:
    start_time = time.time()
    question = state["question"]
    documents = state["documents"]

    batch_size = 20
    filtered_docs = []
    search_needed = False

    def process_batch(batch):
        batch_inputs = [{"question": question, "document": d.page_content} for d in batch]
        return retrieval_grader.batch(batch_inputs)

    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = [executor.submit(process_batch, documents[i:i+batch_size]) for i in range(0, len(documents), batch_size)]

        for future in concurrent.futures.as_completed(futures):
            batch_scores = future.result()
            for doc, score in zip(documents, batch_scores):
                if score["score"] == "yes":
                    filtered_docs.append(doc)
                else:
                    search_needed = True

            if len(filtered_docs) >= 5: # Top 5 docs
                break

    state["documents"] = filtered_docs[:5]  # Limit to top 5 relevant documents
    state["search"] = "Yes" if search_needed and len(filtered_docs) < 3 else "No"
    state["steps"].append("grade_document_retrieval")
    logger.info(f"Grade documents time: {time.time() - start_time:.2f} seconds")
    return state

def web_search(state: GraphState) -> GraphState:
    start_time = time.time()
    web_results = cached_web_search(state["question"])
    state["documents"].extend([Document(page_content=d["content"], metadata={"url": d["url"]}) for d in web_results])
    state["steps"].append("web_search")
    logger.info(f"Web search time: {time.time() - start_time:.2f} seconds")
    return state

def decide_to_generate(state: GraphState) -> str:
    return "search" if state.get("search") == "Yes" and len(state["documents"]) < 3 else "generate"


def add_nodes(workflow: StateGraph):
    workflow.add_node("retrieve", retrieve)
    workflow.add_node("grade_documents", grade_documents)
    workflow.add_node("generate", generate)
    workflow.add_node("web_search", web_search)

def build_graph(workflow: StateGraph):
    workflow.set_entry_point("retrieve")
    workflow.add_edge("retrieve", "grade_documents")
    workflow.add_conditional_edges(
        "grade_documents",
        decide_to_generate,
        {
            "search": "web_search",
            "generate": "generate",
        },
    )
    workflow.add_edge("web_search", "generate")
    workflow.add_edge("generate", END)

def create_workflow() -> StateGraph:
    workflow = StateGraph(GraphState)
    add_nodes(workflow=workflow)
    build_graph(workflow=workflow)
    return workflow.compile()

# Initialize the workflow
rag_workflow = create_workflow()

def process_query(question: str) -> dict:
    initial_state = GraphState(question=question, steps=[])
    result = rag_workflow.invoke(initial_state)
    return result
