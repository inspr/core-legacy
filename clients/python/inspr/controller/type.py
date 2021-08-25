import sys
from ..rest import *
from .controller_client import *

TYPE_ROUTE = "types"

class TypeClient(ControllerClient):
    def get(self, scope:str, type_name:str) -> dict:
        msg_body = {
            "typename": type_name
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + TYPE_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Get Type request: {e}")

    def delele(self, scope:str, type_name:str, dryRun:bool) -> dict:
        msg_body = {
            "typename": type_name,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + TYPE_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Delete Type request: {e}")

    def create(self, scope:str, type:dict, dryRun:bool) -> dict:
        msg_body = {
            "type": type,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + TYPE_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Create Type request: {e}")

    def update(self, scope:str, type:dict, dryRun:bool) -> dict:
        msg_body = {
            "type": type,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + TYPE_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Update Type request: {e}")