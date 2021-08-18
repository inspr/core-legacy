import requests
import sys
from flask.wrappers import Response
from requests.models import HTTPError

def send_post_request(url:str, body) -> Response:
    try:
        resp = requests.post(url, data = body)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e, file=sys.stderr)
        raise HTTPError
    
    return resp