from fastapi import FastAPI
from pydantic import BaseModel
import os
import cohere

app = FastAPI()

cohere_api_key = os.getenv('COHERE_API_KEY')
if not cohere_api_key:
    raise ValueError("COHERE_API_KEY environment variable not set")
cohere_client = cohere.Client(cohere_api_key)

class QuestionRequest(BaseModel):
    question: str

@app.post("/get-answer")
async def get_answer(request: QuestionRequest):
    prompt = f"Question: {request.question}\nAnswer:"
    response = cohere_client.generate(
        model='command-xlarge-nightly',
        prompt=prompt,
        max_tokens=100
    )
    answer = response.generations[0].text.strip()
    return {"question": request.question, "answer": answer}
