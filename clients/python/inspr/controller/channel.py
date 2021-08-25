import sys
from ..models import *
from ..rest import *
from .client import *

CHANNEL_ROUTE = "channels"

class ChannelClient(Client):
    def get(self, scope:str, channel_name:str) -> InsprStructure:
        msg_body = {
            "chname": channel_name
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
            return InsprStructure(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Get Channel request: {e}")

    def delele(self, scope:str, channel_name:str, dryRun:bool) -> Changelog:
        msg_body = {
            "chname": channel_name,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Delete Channel request: {e}")

    def create(self, scope:str, channel:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Create Channel request: {e}")

    def update(self, scope:str, channel:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Update Channel request: {e}")