from fastapi import FastAPI
from pydantic import BaseModel
import time
import random
import socket

app = FastAPI()

class Parametros(BaseModel):
    tasa_aprendizaje: float
    profundidad: int

@app.post("/entrenar")
def entrenar_modelo(params: Parametros):
    print("Recibi tarea con tasa:", params.tasa_aprendizaje, "y profundidad:", params.profundidad)
    time.sleep(2) 
    precision_inventada = round(random.uniform(0.7, 0.99), 2)
    
    nombre_worker = socket.gethostname()
    
    return {
        "precision_final": precision_inventada,
        "worker_id": nombre_worker
    }