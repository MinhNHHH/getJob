from fastapi import APIRouter

router = APIRouter(
    prefix="/company",          # all routes will start with /company
    tags=["company"],           # for docs grouping
)

@router.get("/")
def list_companies():
    return [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]

@router.get("/{company_id}")
def get_company(company_id: int):
    return {"id": company_id, "name": f"Comapny {company_id}"}

@router.post("/")
def create_user(company: dict):
    return {"message": "User created", "user": company}
