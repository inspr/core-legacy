import requests
from flask.wrappers import Response
from requests.models import HTTPError

def sendPostRequest(url:str, body) -> Response:
    print(body)
    try:
        resp = requests.post(url, data = body)
        print("resp = ", resp)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e)
        raise HTTPError
    
    return resp