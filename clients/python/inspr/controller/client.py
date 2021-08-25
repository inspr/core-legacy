import os

INSPRD_ADDRESS = "INSPR_INSPRD_ADDRESS"
SCOPE = "INSPR_CONTROLLER_SCOPE"
TOKEN = "INSPR_CONTROLLER_TOKEN"
TOKEN_PATH = ".inspr/token"

class Client:
    def __init__(self, insprd_url, scope = "") -> None:
        self.url = insprd_url
        self.scope = scope
    
    def get_header_with_scope(self, scope:str) -> dict:
        headers = {
            "content-type": "application/json",
            "Scope": scope,
            "Authorization": self.get_token() 
        }
        return headers

    def get_token(self) -> str:
        home_dir = os.path.expanduser('~')
        token_path = os.path.join(home_dir, TOKEN_PATH)
        f = open(token_path, "r")
        token = f.read().strip()

        if token is None:
            raise Exception("cant get token")
        
        token = "Bearer " + token
        return token

    def set_token(self, token:str) -> None:
        os.environ[TOKEN] = token[len("Bearer "):]
