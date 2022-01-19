import sys
from ..models import *
from ..rest import *
from .client import *

ALIAS_ROUTE = "alias"

class AliasClient(Client):
    def get(self, scope:str, name:str) -> InsprStructure:
        msg_body = {
            "name": name
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + ALIAS_ROUTE, body=msg_body, headers=headers)
            return InsprStructure(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Get Alias request: {e}")

    def delele(self, scope:str, name:str, dryRun:bool) -> Changelog:
        msg_body = {
            "name": name,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + ALIAS_ROUTE, body=msg_body, headers=headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Delete Alias request: {e}")

    def create(self, scope:str, alias:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "alias": alias,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + ALIAS_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Create Alias request: {e}")

    def update(self, scope:str, alias:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "alias": alias,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + ALIAS_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Update Alias request: {e}")