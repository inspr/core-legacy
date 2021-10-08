from flask import request
import sys
from http import HTTPStatus
import inspr

def main():
    client = inspr.Client()
    
    @client.handle_route("/add")
    def add_handler():
        data = request.get_json(force=True)
        op1 = data["op1"]
        op2 = data["op2"]

        total = op1 + op2
        print(total, file=sys.stderr)

        return str(total), HTTPStatus.OK
    
    @client.handle_route("/sub")
    def sub_handler():
        data = request.get_json(force=True)
        op1 = data["op1"]
        op2 = data["op2"]

        total = op1 - op2
        print(total, file=sys.stderr)

        return str(total), HTTPStatus.OK

    @client.handle_route("/mul")
    def sub_handler():
        data = request.get_json(force=True)
        op1 = data["op1"]
        op2 = data["op2"]

        total = op1 * op2
        print(total, file=sys.stderr)

        return str(total), HTTPStatus.OK

    client.run()
    

if __name__ == "__main__":
    main()