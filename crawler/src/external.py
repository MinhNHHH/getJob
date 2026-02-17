import requests 
import json

def notification(payload):
	try:
		requests.request("post", "http://localhost:8000/api/crawler/notify", 
			data=json.dumps(payload), 
			headers = {
				'Content-Type': 'application/json'
			}
		)
	except Exception as error:
		print(error)