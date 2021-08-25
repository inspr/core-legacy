import sys
from ..rest import *
from .controller_client import *

CHANNEL_ROUTE = "channels"

class ChannelClient(ControllerClient):
    def get(self, scope:str, channel_name:str) -> dict:
        msg_body = {
            "chname": channel_name
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Get Channel request: {e}")

    def delele(self, scope:str, channel_name:str, dryRun:bool) -> dict:
        msg_body = {
            "chname": channel_name,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Delete Channel request: {e}")

    def create(self, scope:str, channel:dict, dryRun:bool) -> dict:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Create Channel request: {e}")

    def update(self, scope:str, channel:dict, dryRun:bool) -> dict:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Update Channel request: {e}")