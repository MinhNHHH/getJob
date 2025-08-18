import sys,os
sys.path.append(os.path.abspath(".."))
sys.path.append(os.path.abspath("."))

from datetime import datetime
from fastapi import Depends, FastAPI, HTTPException, APIRouter

from backend.src.router import company, jobs

app = FastAPI(
  title="My API",
  description="This is my FastAPI project",
  version="1.0.0",
  docs_url="/swagger",      # Swagger UI
  redoc_url="/redoc",       # ReDoc
  openapi_url="/api/openapi.json",  # OpenAPI schema
)

app_router = APIRouter(prefix="/api")
app_router.include_router(company.router)
app_router.include_router(jobs.router)

app.include_router(app_router)

@app.get("/health")
def health():
  return {
    "message": "Server is working now",
    "time": datetime.now()
  }
