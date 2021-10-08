import inspr
import sys
import time

ADD_PATH = "add"
SUB_PATH = "sub"
MUL_PATH = "mul"

API_NAME = "api"

def main():
    client = inspr.Client()
    msg_body = {
        "op1": 5,
        "op2": 3,
    }

    print("Python Client")

    while True:
        try:
            resp = client.send_request(API_NAME, ADD_PATH, method="POST", body=msg_body)
            res = resp.json()
            print("The result of 5 + 3 = ", res)
        except:
            print("An error has occured", file=sys.stderr)
        
        time.sleep(3)

        try:
            resp = client.send_request(API_NAME, SUB_PATH, method="POST", body=msg_body)
            res = resp.json()
            print("The result of 5 - 3 = ", res)
        except:
            print("An error has occured", file=sys.stderr)
        
        time.sleep(3)

        try:
            resp = client.send_request(API_NAME, MUL_PATH, method="POST", body=msg_body)
            res = resp.json()
            print("The result of 5 * 3 = ", res)
        except:
            print("An error has occured", file=sys.stderr)
        
        time.sleep(3)




if __name__ == "__main__":
    main()