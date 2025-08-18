from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, scoped_session, Session

# from pymongo import MongoClient

from src import cfgs

############################### ANALYTICS DB ###############################
engine = create_engine(cfgs.get_default("DB_CONNECTION_URI"),
    # http://docs.sqlalchemy.org/en/latest/core/pooling.html
    pool_size=20, max_overflow=-1)

session_factory = sessionmaker(bind=engine)
SessionLocal = scoped_session(session_factory)

Base = declarative_base()

############################### CRAWLER DB ###############################

# _crawler_client = MongoClient(cfg.get_default("CRAWLER_CONNECTION_URI"))
# CRAWLER_DB = _crawler_client[cfg.get_default("CRAWLER_DBNAME")]
