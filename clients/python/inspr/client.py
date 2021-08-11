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
        self.read_port = os.getenv(SIDECAR_READ_PORT)
        self.write_address = "http://localhost:" + str(os.getenv(SIDECAR_WRITE_PORT))
        self.app = Flask(__name__)

    def write_message(self, channel:str, msg) -> None:
        msg_body = {
            "data": msg
        }
        json_obj = json.dumps(msg_body, indent=4)
        try:
            send_post_request(self.write_address + "/" + channel, json_obj)
        except Exception as e:
            print(f"Error while trying to write message: {e}")
            raise Exception("failed to deliver message: channel: {}".format(channel))

    def handle_channel(self, channel:str) -> Callable[[Callable[[Any], Any]], Callable[[Any], Any]]:
        def wrapper(handle_func: Callable[[Any], Any]):
            def route_func():
                data = request.get_json(force=True)
                try:
                    handle_func(data["data"])
                except:
                    err = "Error handling message"
                    return err, HTTPStatus.INTERNAL_SERVER_ERROR

                return '', HTTPStatus.OK

            self.app.add_url_rule("/" + channel, endpoint = channel, view_func = route_func, methods=["POST"])
            return handle_func
        return wrapper

    def run(self) -> None:
        links = []
        for rule in self.app.url_map.iter_rules():
            links.append(rule.endpoint)
        print("registered routes =", links, file=sys.stderr)

        self.app.run(port=self.read_port)
