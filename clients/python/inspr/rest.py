import requests
import json
import sys
from flask.wrappers import Response
from requests.models import HTTPError

def send_post_request(url:str, body = {}, headers = {}) -> Response:
    try:
        resp = requests.post(url, data = json.dumps(body), headers=headers)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e, file=sys.stderr)
        raise HTTPError
    
    return resp

def send_update_request(url:str, body = {}, headers = {}) -> Response:
    try:
        resp = requests.put(url, data = json.dumps(body), headers=headers)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e, file=sys.stderr)
        raise HTTPError
    
    return resp

def send_get_request(url:str, body = {}, headers = {}) -> Response:
    try:
        resp = requests.get(url, data=json.dumps(body), headers=headers)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e, file=sys.stderr)
        raise HTTPError
    
    return resp

def send_delete_request(url:str, body = {}, headers = {}) -> Response:
    try:
        resp = requests.delete(url, data=json.dumps(body), headers=headers)
        resp.raise_for_status()
    
    except requests.exceptions.RequestException as e:
        print(e, file=sys.stderr)
        raise HTTPError
    
    return resp
