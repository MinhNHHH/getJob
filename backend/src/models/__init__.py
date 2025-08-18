# NOTE: when create new models, we need to include them here
# so that alembic could correctly generate the migration
import sys,os
sys.path.append(os.path.abspath(".."))
from backend.src.database import Base
from backend.src.models.company import Company
from backend.src.models.job import Job

__all__ = [
    Company,
    Job
]