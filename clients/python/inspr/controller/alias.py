import sys
from ..rest import *
from .controller_client import *

ALIAS_ROUTE = "alias"

class AliasClient(ControllerClient):
    def get(self, scope:str, key:str) -> dict:
        msg_body = {
            "key": key
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + ALIAS_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Get Alias request: {e}")

    def delele(self, scope:str, key:str, dryRun:bool) -> dict:
        msg_body = {
            "key": key,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + ALIAS_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Delete Alias request: {e}")

    def create(self, scope:str, target:str, alias:dict, dryRun:bool) -> dict:
        msg_body = {
            "alias": alias,
            "target": target,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + ALIAS_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Create Alias request: {e}")

    def update(self, scope:str, target:str, alias:dict, dryRun:bool) -> dict:
        msg_body = {
            "alias": alias,
            "target": target,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + ALIAS_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Update Alias request: {e}")