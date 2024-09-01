FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN sed -i '/gpt4all==/d' requirements.txt && \
    echo "gpt4all" >> requirements.txt && \
    sed -i 's/faiss-cpu==1.8.0.post1/faiss-cpu/' requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "app.py"]