from backend.src.task.base import runner
from backend.src.database import SessionLocal
import backend.src.models.company as c
import backend.src.models.job as j

from pydantic import BaseModel
import logging

class JobDetail(BaseModel):
    title: str
    company_name: str
    company_uri: str
    location: str
    description: str

class TaskInfo(BaseModel):
    id: str
    task_name: str
    start_time: float
    end_time: float
    results: list[JobDetail]

class Summary(BaseModel):
    company_inserted: int = 0
    company_updated: int = 0
    job_inserted: int = 0
    job_updated: int = 0

_logger = logging.getLogger(__name__)

def save_company(db, job_detail: JobDetail, summary: Summary):
    company_info = {
        "name": job_detail.company_name,
        "url": job_detail.company_uri,
    }
    [company, method] = c.upsert(db, company_info)
    if method == "insert":
        summary.company_inserted += 1
    else:
        summary.company_updated += 1
    return company

def save_job(db, job_detail: JobDetail, summary: Summary):
    job_info = {
        "title": job_detail.title,
        "company_id": job_detail.company_id,
        "location": job_detail.location,
        "description": job_detail.description,
    }
    [job, method] = j.upsert(db, job_info)
    if method == "insert":
        summary.job_inserted += 1
    else:
        summary.job_updated += 1
    return job

@runner(task_name="crawl-job")
def crawl_job(task_info: TaskInfo, **kwargs):
    db = SessionLocal()
    summary = Summary()

    for job in task_info.results:
        job_detail = JobDetail(**job)
        save_company(db, job_detail, summary)
        save_job(db, job_detail, summary)

    return summary
    
    