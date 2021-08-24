import sys
import rest
import controller_client as controller

APP_ROUTE = "apps"

class AppClient(controller.ControllerClient):
    def Get(self, scope:str) -> dict:
        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_get_request(self.url + "/" + APP_ROUTE, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Get App request: {e}")

    def Delete(self, scope:str, dryRun:bool) -> dict:
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

    def Post(self, scope:str, app:dict, dryRun:bool) -> dict:
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

    def Update(self, scope:str, app:dict, dryRun:bool) -> dict:
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