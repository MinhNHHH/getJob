import requests
from datetime import datetime
from bs4 import BeautifulSoup
import uuid
import json
from src.database import r
from src.external import notification


def different_time_in_miliseconds(start, end):
	return end - start

def create_task_history(task_name, task_details):
	task_id = str(uuid.uuid4())
	task_info = {
		"id": task_id,
		"task_name": task_name,
		"start_time": datetime.now().timestamp(),
		"end_time": None,
		"details": task_details,
		"error": None,
		"status": "started",
	}
	r.set(task_id, json.dumps(task_info))
	return task_id

def on_task_done(task_id, results):
	task_info = json.loads(r.get(task_id))
	now = datetime.now().timestamp()
	task_info.update({
		"end_time": now,
		"status": "done",
		"results": results,
		"duration": different_time_in_miliseconds(task_info["start_time"], now),
	})
	print(task_id)
	r.set(task_id, json.dumps(task_info))


def on_task_error(task_id, error):
	task_info = json.loads(r.get(task_id))
	now = datetime.now().timestamp()
	task_info.update({
		"end_time": now,
		"status": "error",
		"duration": different_time_in_miliseconds(task_info["start_time"], now),
		"error": error
	})
	r.set(task_id, json.dumps(task_info))

def runner(task_name):
	def decorator(func):
		def wrapper(*args, **kwargs):
			task_detail = kwargs.get("task_details")
			task_id = create_task_history(task_name, task_detail)
			kwargs["task_id"] = task_id
			result = None
			try:
				print(f"[{task_name}] Starting task...")
				result = func(*args, **kwargs)
			except Exception as e:
				on_task_error(task_id, f"[{task_name}] Task failed with error: {e}")
				print(f"[{task_name}] Task failed with error: {e}")
				raise
			finally:
				on_task_done(task_id, result)
				notification({
					"task_id": task_id,
					"task_name": task_name,
				})
		return wrapper
	return decorator



class JobCrawlBase:
	def __init__(self, pages=1, job=''):
		self.pages = pages
		self.job = job
		self.headers = {
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept-Language": "en-US,en;q=0.9",
		}
		
	def get_request(self, url):
		try:
			response = requests.get(url, headers=self.headers)
			return response
		except Exception as error:
			print(f"Request error: {error}")
			return None

	def parser_html(self, url):
		response = self.get_request(url)
		if not response:
			return None
		soup = BeautifulSoup(response.text, "lxml")
		return soup

