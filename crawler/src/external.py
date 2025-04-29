import requests 
import json

def notification(payload):
	try:
		requests.request("post", "http://localhost:8080/notification", 
			data=json.dumps(payload), 
			headers = {
				'Content-Type': 'application/json'
			}
		)
	except Exception as error:
		print(error)