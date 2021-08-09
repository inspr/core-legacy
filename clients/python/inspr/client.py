import os
from flask import Flask, request
from http import HTTPStatus
from .rest import *
from types import FunctionType

SIDECAR_READ_PORT = "INSPR_SCCLIENT_READ_PORT"
SIDECAR_WRITE_PORT = "INSPR_LBSIDECAR_WRITE_PORT"

class Client:
    def __init__(self) -> None:
        self.readPort = os.getenv(SIDECAR_READ_PORT)
        self.writeAddress = ":" + str(os.getenv(SIDECAR_WRITE_PORT))
        self.app = Flask(__name__)

    
    def writeMessage(self, channel:str, msg) -> None:
        msgBody = {
            "data": msg
        }
        try:
            sendPostRequest(self.writeAddress + "/" + channel, msgBody)
        except:
            print("Error while trying to write message")
            raise ValueError

    
    def handleChannel(self, channel:str, handleFunc:FunctionType) -> None:
        def routeFunc():
            data = request.json
            try:
                handleFunc(data)
            except:
                err = "Error handling message"
                return err, HTTPStatus.INTERNAL_SERVER_ERROR

            return '', HTTPStatus.OK

        self.app.add_url_rule("/" + channel, view_func = routeFunc)


    def run(self) -> None:
        self.app.run(port=self.readPort)





