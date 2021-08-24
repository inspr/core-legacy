from inspr import *
import sys

PING_INPUT_CHANNEL = "pinginput"
PING_OUTPUT_CHANNEL = "pingoutput"

def main():
    client = Client()
    msg = "Ping!"

    @client.handle_channel(PING_INPUT_CHANNEL)
    def read_pong_and_send_ping(data):

        if data == 'Pong!':
            print(data, file=sys.stderr)
        else:
            print('Not received Pong', file=sys.stderr)

        try:
            client.write_message(PING_OUTPUT_CHANNEL, msg)
            return Response(status=200)
        except:
            raise Exception
        


    try:
        client.write_message(PING_OUTPUT_CHANNEL, msg)
    except:
        print("An error has occured", file=sys.stderr)
        return

    client.run()

if __name__ == "__main__":
    main()
