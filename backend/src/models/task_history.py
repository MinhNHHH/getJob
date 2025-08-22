from backend.src.database import Base
from sqlalchemy import Column, Integer, Text, DateTime

class TaskHistory(Base):
  __tablename__   = "task_history"

  id               = Column(Integer, primary_key=True, index=True)
  name             = Column(Text)
  status           = Column(Text)
  details          = Column(Text)
  summary          = Column(Text)
  error            = Column(Text)
  start_time       = Column(DateTime)
  end_time         = Column(DateTime)
  duration         = Column(Integer)