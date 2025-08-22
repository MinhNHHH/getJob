from backend.src.database import Base
import backend.src.models.utils as mutils
from sqlalchemy import Column, DateTime, Integer, String
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship


class Company(Base):
    __tablename__ = "company"

    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String(255), nullable=False, index=True)
    url = Column(String(255), nullable=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)

    # reverse relation
    jobs = relationship("Job", back_populates="company", cascade="all, delete-orphan")


def upsert(db, company_info):
    """
    Update if company exists, otherwise insert
    """
    existed = db.query(Company).filter(Company.name == company_info["name"]).first()
    if not existed:
        return mutils.insert(db, Company, company_info), "insert"
    else:
        return mutils.update(db, existed, company_info), "update"