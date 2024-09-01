#!/bin/sh

ollama serve &

until ollama list >/dev/null 2>&1; do
    echo "Waiting for Ollama to be ready..."
    sleep 5
done

ollama pull phi3:mini
ollama pull tinyllama
ollama pull nomic-embed-text

wait