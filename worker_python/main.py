from fastapi import FastAPI
from pydantic import BaseModel
import time
import random

app = FastAPI()

class ConfiguracionModelo(BaseModel):
    tasa_aprendizaje: float
    profundidad_arbol: int
    estimadores: int

@app.post("/entrenar")
def entrenar_modelo(config: ConfiguracionModelo):
    print(f"Recibiendo tarea: {config}")
    time.sleep(2) 
    precision_ejemplo = random.uniform(0.70, 0.98)
    print(f"Precisión: {precision_ejemplo:.4f}")
    
    return {
        "estado": "completado",
        "precision_final": round(precision_simulada, 4),
        "parametros_usados": config.dict()
    }