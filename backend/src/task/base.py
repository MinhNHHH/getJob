import json
import logging
import traceback
from datetime import datetime

from backend.src.database import SessionLocal
import backend.src.models.task_history as models
import backend.src.models.utils as mutils

_logger = logging.getLogger(__name__)

def diff_time_in_microseconds(time1, time2):
  """
  Returns the difference between 2 time in microseconds.
  Use to compute duration of task.
  """
  # apply an abs to ensure the different is always positive
  return abs((time1 - time2).total_seconds() * 1000)

def create_task_history(task_name, task_details):
  db  = SessionLocal()
  task_info = {
    "name": task_name,
    "start_time": datetime.now(),
    "end_time": None,
    "duration": None,
    "status": "started",
    "details": json.dumps(task_details),
    "summary": None,
    "error": None
  }
  task_history = mutils.insert(db, models.TaskHistory, task_info)
  return task_history

def on_task_done(task_history, summary):
  db  = SessionLocal()
  now = datetime.now()
  update_task_history = {
    "summary": json.dumps(summary),
    "status": "done",
    "end_time": now,
    "duration": diff_time_in_microseconds(now, task_history.start_time)
  }
  mutils.update(db, task_history, update_task_history)

def on_task_error(task_history, error):
  db  = SessionLocal()
  now = datetime.now()
  update_task_history = {
    "error": str(error),
    "status": "error",
    "end_time": now,
    "duration": diff_time_in_microseconds(now, task_history.start_time)
  }
  mutils.update(db, task_history, update_task_history)

def runner(task_name):
  def decorator(func):
    def wrapper(*args, **kwargs):
      task_details = kwargs.get("task_details", {})
      task_history = create_task_history(task_name, task_details)
      _logger.info(f"Task started {task_name} {task_history.id}")
      result = None
      try:
        # populate task_id to the task function
        kwargs["task_id"] = task_history.id
        result = func(*args, **kwargs)
      except Exception as e:
        _logger.error(f"Task error {task_name} {task_history.id} with error: {e}")
        _logger.error(traceback.format_exc())
        on_task_error(task_history, e)
      else:
        _logger.info(f"Task done {task_name} {task_history.id}: {result}")
        on_task_done(task_history, result)
      return result
    return wrapper
  return decorator