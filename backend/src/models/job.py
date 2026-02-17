from backend.src.database import Base
import backend.src.models.utils as mutils
from sqlalchemy import Column, Integer, String, Text, DateTime, ForeignKey, func
from sqlalchemy.orm import relationship

class Job(Base):
    __tablename__ = "jobs"

    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String(255), nullable=False)
    description = Column(Text, nullable=False)
    company_id = Column(Integer, ForeignKey("company.id"), nullable=False)
    location = Column(String, nullable=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)

    # relationship to Company
    company = relationship("Company", back_populates="jobs")

def upsert(db, job_info):
    """
    Update if job exists, otherwise insert
    """
    existed = db.query(Job).filter(
        Job.title == job_info["title"],
        Job.location == job_info["location"],
        Job.company_id == job_info["company_id"]
    ).first()
    if not existed:
        return mutils.insert(db, Job, job_info), "insert"
    else:
        return mutils.update(db, existed, job_info), "update"