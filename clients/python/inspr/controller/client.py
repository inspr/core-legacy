import os

INSPRD_ADDRESS = "INSPR_INSPRD_ADDRESS"
SCOPE = "INSPR_CONTROLLER_SCOPE"
TOKEN = "INSPR_CONTROLLER_TOKEN"
TOKEN_PATH = ".inspr/token"

class Client:
    def __init__(self, insprd_url, global_scope = "") -> None:
        self.url = insprd_url
        self.global_scope = global_scope
    
    def get_header_with_scope(self, scope:str) -> dict:
        headers = {
            "content-type": "application/json",
            "Scope": self.join_scopes(self.global_scope, scope),
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

    def join_scopes(self, s1:str, s2:str) -> str:
        if s2 == "":
            return s1
        
        separator = ""
        if s1 != "":
            separator = "."
        
        new_scope = s1 + separator + s2
        return new_scope