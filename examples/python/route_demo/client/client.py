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
        "op1": 1,
        "op2": 2,
    }

    print("Python Client")

    while True:
        try:
            resp = client.send_request(API_NAME, ADD_PATH, method="POST", body=msg_body)
        except:
            print("An error has occured", file=sys.stderr)
            return

        print("The result is ", resp.json())
        time.sleep(10)


if __name__ == "__main__":
    main()