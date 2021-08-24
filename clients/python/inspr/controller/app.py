import sys
import rest
import controller_client as controller

APP_ROUTE = "apps"

class AppClient(controller.ControllerClient):
    def get(self, scope:str) -> dict:
        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_get_request(self.url + "/" + APP_ROUTE, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Get App request: {e}")

    def delete(self, scope:str, dryRun:bool) -> dict:
        msg_body = {
            "dryRun": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_delete_request(self.url + "/" + APP_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Delete App request: {e}")

    def post(self, scope:str, app:dict, dryRun:bool) -> dict:
        msg_body = {
            "app": app,
            "dryRun": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_post_request(self.url + "/" + APP_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Post App request: {e}")

    def update(self, scope:str, app:dict, dryRun:bool) -> dict:
        msg_body = {
            "app": app,
            "dryRun": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_update_request(self.url + "/" + APP_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Update App request: {e}")