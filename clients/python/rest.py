import requests
from flask.wrappers import Response
from requests.models import HTTPError

def sendPostRequest(url:str, body) -> Response:
    try:
        resp = requests.post(url, data = body)
        resp.raise_for_status()
    
    except requests.HTTPError as exception:
        print(exception)
        raise HTTPError
    
    return resp