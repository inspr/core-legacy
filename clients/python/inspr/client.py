import os
import json
import sys
from flask import Flask, request
from http import HTTPStatus
from .rest import *
from typing import Callable, Any

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
        except Exception as e:
            print(f"Error while trying to write message: {e}")
            raise Exception("failed to deliver message: channel: {}".format(channel))

    def handleChannel(self, channel:str) -> Callable[[Callable[[Any], Any]], Callable[[Any], Any]]:
        def wrapper(handleFunc: Callable[[Any], Any]):
            def routeFunc():
                data = request.get_json(force=True)
                try:
                    handleFunc(data["data"])
                except:
                    err = "Error handling message"
                    return err, HTTPStatus.INTERNAL_SERVER_ERROR

                return '', HTTPStatus.OK

            self.app.add_url_rule("/" + channel, endpoint = channel, view_func = routeFunc, methods=["POST"])
            return handleFunc
        return wrapper

    def run(self) -> None:
        links = []
        for rule in self.app.url_map.iter_rules():
            links.append(rule.endpoint)
        print("registered routes =", links, file=sys.stderr)

        self.app.run(port=self.readPort)
