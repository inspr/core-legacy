import json

class Changelog:
    def __init__(self, changelog:dict) -> None:
        self.obj = changelog
    
    def __str__(self) -> str:
        return f'Changelog({ json.dumps(self.obj, sort_keys=False, indent=4) })'

class InsprStructure:
    def __init__(self, structure:dict) -> None:
        self.obj = structure
    
    def __str__(self) -> str:
        return f'InsprStructure({ json.dumps(self.obj, sort_keys=False, indent=4) })'