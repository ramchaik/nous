import os
import uuid
from langchain_core.output_parsers import StrOutputParser, JsonOutputParser
from graph import create_workflow

custom_graph = create_workflow()

os.environ['USER_AGENT'] = 'Nous_LLM/1.0'

class LangChainModel:
    @staticmethod
    def predict(input_text):
        config = {"configurable": {"thread_id": str(uuid.uuid4())}}
        state_dict = custom_graph.invoke(
            {"question": input_text, "steps": []}, config
        )
        return {
            "response": state_dict["generation"],
            "steps": state_dict["steps"]
        }

if __name__ == "__main__":
    model = LangChainModel()
    result = model.predict("Your question here")
    print(result)