from fastapi import APIRouter

router = APIRouter(
    prefix="/jobs",          # all routes will start with /jobs
    tags=["jobs"],           # for docs grouping
)

@router.get("/")
def list_companies():
    return [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]

@router.get("/{job_id}")
def get_jobs(job_id: int):
    return {"id": job_id, "name": f"Comapny {job_id}"}

@router.post("/")
def create_user(job: dict):
    return {"message": "User created", "user": job}
