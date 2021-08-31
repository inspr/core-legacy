from .app import *
from .channel import *
from .type import *
from .alias import *

class ControllerClient:
    def __init__(self, insprd_url, scope) -> None:
        self.app = AppClient(insprd_url, scope)
        self.channel = ChannelClient(insprd_url, scope)
        self.type = TypeClient(insprd_url, scope)
        self.alias = AliasClient(insprd_url, scope)