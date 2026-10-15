from fastapi import FastAPI
from pydantic import BaseModel
import time
import random
import socket
import requests

app = FastAPI()

CONSUL_URL = "http://consul:8500/v1/agent/service"
SERVICE_NAME = "python-worker"
HOSTNAME = socket.gethostname()
IP_ADDRESS = socket.gethostbyname(HOSTNAME)
PORT = 8000

@app.on_event("startup")
def registrar_en_consul():
    payload = {
        "ID": f"{SERVICE_NAME}-{HOSTNAME}",
        "Name": SERVICE_NAME,
        "Address": IP_ADDRESS,
        "Port": PORT,
        "Check": {
            "HTTP": f"http://{IP_ADDRESS}:{PORT}/health",
            "Interval": "10s",
            "Timeout": "5s"
        }
    }
    try:
        requests.put(f"{CONSUL_URL}/register", json=payload)
        print("Registrado en el Consul")
    except Exception as e:
        print("Fallo al registrar en el Consul:", e)

@app.on_event("shutdown")
def borrar_de_consul():
    try:
        requests.put(f"{CONSUL_URL}/deregister/{SERVICE_NAME}-{HOSTNAME}")
        print("Borrado del Consul")
    except Exception as e:
        pass

@app.get("/health")
def health_check():
    return {"status": "ok", "worker": HOSTNAME}

class Parametros(BaseModel):
    tasa_aprendizaje: float
    profundidad: int

@app.post("/entrenar")
def entrenar_modelo(params: Parametros):
    time.sleep(2) 
    precision_inventada = round(random.uniform(0.7, 0.99), 2)
    return {
        "precision_final": precision_inventada,
        "worker_id": HOSTNAME
    }