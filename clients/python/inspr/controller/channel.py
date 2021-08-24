import sys
import rest
import controller_client as controller

CHANNEL_ROUTE = "channels"

class ChannelClient(controller.ControllerClient):
    def get(self, scope:str, channel_name:str) -> dict:
        msg_body = {
            "chname": channel_name
        }
        
        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_get_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
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
            resp = rest.send_delete_request(self.url + "/" + CHANNEL_ROUTE, body=msg_body, headers=headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Delete Channel request: {e}")

    def post(self, scope:str, channel:dict, dryRun:bool) -> dict:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_post_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Post Channel request: {e}")

    def update(self, scope:str, channel:dict, dryRun:bool) -> dict:
        msg_body = {
            "channel": channel,
            "dry": dryRun
        }

        headers = self.get_header_with_scope(scope)

        try:
            resp = rest.send_update_request(self.url + "/" + CHANNEL_ROUTE, msg_body, headers)
            print(resp, file=sys.stderr)
            return resp
        except Exception as e:
            raise Exception(f"Error while send a Update Channel request: {e}")