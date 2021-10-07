import inspr
from flask import request
import sys
import time

ADD_PATH = "add"
SUB_PATH = "sub"
MUL_PATH = "mul"

API_NAME = "api"

def main():
    client = inspr.Client
    msg_body = {
        "op1": 1,
        "op2": 2,
    }

    while True:
        client.send_request(API_NAME, ADD_PATH, "POST", msg_body)
        data = request.get_json(force=True)
        print(data, file=sys.stderr)
        
        time.sleep(2)


if __name__ == "__main__":
    main()