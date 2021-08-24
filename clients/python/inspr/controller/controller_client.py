import os

INSPRD_ADDRESS = "INSPR_INSPRD_ADDRESS"
SCOPE = "INSPR_CONTROLLER_SCOPE"
TOKEN = "INSPR_CONTROLLER_TOKEN"

class ControllerClient:
    def __init__(self) -> None:
        self.url = os.getenv(INSPRD_ADDRESS)
        self.controller_scope = os.getenv(SCOPE)
    
    def get_header_with_scope(self, scope:str) -> dict:
        headers = {
            "content-type": "application/json",
            "Scope": scope,
            "Authorization": self.get_token() 
        }
        return headers

    def get_token(self) -> str:
        token = "Bearer " + os.getenv(TOKEN)
        return token

    def set_token(self, token:str) -> None:
        os.environ[TOKEN] = token[len("Bearer "):]
