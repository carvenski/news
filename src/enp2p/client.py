import websocket
import time

# client 1
def on_message(ws, message):
    print("recv message: %s" % message)
    ws.send("i have got your message %s" % message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    print("register user wechat info")
    ws.send("id1\nmale\ntest1\ntest1\nhahaha1")

if __name__ == "__main__":
    ws = websocket.WebSocketApp("ws://192.168.153.1:35000",
                              on_message = on_message,
                              on_error = on_error,
                              on_close = on_close)
    ws.on_open = on_open
    ws.run_forever()


# client 2
# from websocket import create_connection
# ws = create_connection("ws://192.168.153.1:35000") 
# ws.send("id2\nfemale\ntest2\ntest2\ntest2")



