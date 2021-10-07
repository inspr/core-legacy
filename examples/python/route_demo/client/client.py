import inspr
from flask import request
import sys
import time

ADD_PATH = "add"
SUB_PATH = "sub"
MUL_PATH = "mul"

API_NAME = "api"

def main():
    client = inspr.Client()
    msg_body = {
        "op1": 1,
        "op2": 2,
    }

    print("Python Client")

    for i in range(0,5):
        resp = client.send_request(API_NAME, ADD_PATH, method="POST", body=msg_body)
        print("Request ", i)
        print("Response () = ", resp, file=sys.stderr)
        print(resp.json())
        time.sleep(10)


if __name__ == "__main__":
    main()