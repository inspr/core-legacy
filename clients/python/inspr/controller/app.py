import sys
from ..models import *
from ..rest import *
from .client import *

APP_ROUTE = "apps"

class AppClient(Client):
    def get(self, scope:str) -> InsprStructure:
        headers = self.get_header_with_scope(scope)

        try:
            resp = send_get_request(self.url + "/" + APP_ROUTE, headers=headers)
            return InsprStructure(json.loads(resp.text))
        
        except Exception as e:
            raise Exception(f"Error while send a Get App request: {e}")

    def delete(self, scope:str, dryRun:bool) -> Changelog:
        msg_body = {
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_delete_request(self.url + "/" + APP_ROUTE, body=msg_body, headers=headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Delete App request: {e}")

    def create(self, scope:str, app:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "app": app,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_post_request(self.url + "/" + APP_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Create App request: {e}")

    def update(self, scope:str, app:dict, dryRun:bool) -> Changelog:
        msg_body = {
            "app": app,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = send_update_request(self.url + "/" + APP_ROUTE, msg_body, headers)
            return Changelog(json.loads(resp.text))
        except Exception as e:
            raise Exception(f"Error while send a Update App request: {e}")