from client import *

PONG_INPUT_CHANNEL = "ponginput"
PONG_OUTPUT_CHANNEL = "pongoutput"

def main():
    client = Client()
    msg = "Pong!"

    def readPingAndSendPong(data):
        if data == "Ping!":
            print(data)
        else:
            print("Not received pong :(")
            print("received = ", data)
            raise Exception
        
        try:
            client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
        except:
            raise Exception


    try:
        client.writeMessage(PONG_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured")
        return

    client.handleChannel(PONG_INPUT_CHANNEL, readPingAndSendPong)
    client.run()