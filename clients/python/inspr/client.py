import os
import json
import sys
from flask import Flask, request, url_for
from http import HTTPStatus
from .rest import *
from types import FunctionType

SIDECAR_READ_PORT = "INSPR_SCCLIENT_READ_PORT"
SIDECAR_WRITE_PORT = "INSPR_LBSIDECAR_WRITE_PORT"

class Client:
    def __init__(self) -> None:
        self.readPort = os.getenv(SIDECAR_READ_PORT)
        self.writeAddress = "http://localhost:" + str(os.getenv(SIDECAR_WRITE_PORT))
        self.app = Flask(__name__)

    
    def writeMessage(self, channel:str, msg) -> None:
        msgBody = {
            "data": msg
        }
        jsonObj = json.dumps(msgBody, indent=4)
        try:
            sendPostRequest(self.writeAddress + "/" + channel, jsonObj)
        except:
            print("Error while trying to write message")
            raise ValueError

    
    def handleChannel(self, channel:str, handleFunc:FunctionType) -> None:
        print("channel =", channel)
        def routeFunc():
            data = request.get_json(force=True)
            try:
                handleFunc(data["data"])
            except:
                err = "Error handling message"
                return err, HTTPStatus.INTERNAL_SERVER_ERROR

            return '', HTTPStatus.OK

        self.app.add_url_rule("/" + channel, endpoint = channel, view_func = routeFunc, methods=["POST"])


    def run(self) -> None:
        links = []
        for rule in self.app.url_map.iter_rules():
            links.append(rule.endpoint)
        print("registered routes =", links, file=sys.stderr)

        self.app.run(port=self.readPort)





