import sys,os
sys.path.append(os.path.abspath(".."))
sys.path.append(os.path.abspath("."))
import logging
import threading
from datetime import datetime

from fastapi import FastAPI, APIRouter

from backend.src.router import company, jobs
from backend.src.task.crawl_job import crawl_job

_logger = logging.getLogger(__name__)

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

@app.post("/api/crawler/notify")
def notify_crawl_job(task_info: TaskInfo):
  _logger.info(f"Got a task {task_info.__dict__}")
  if task_info.task_name == "crawl-job":
    crawl_thread = threading.Thread(target=crawl_job, args=(task_info,))
    crawl_thread.start()
  else:
    _logger.error(f"Invalid task name {task_info.task_name}")
    return {"error": "Invalid task name"}