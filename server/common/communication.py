import common.message as message

# TODO: if more servers are added, this should be in the config file.
SERVER_ID = "server"

BUFFER_SIZE = 8192 # 8 KiB

def receive_client_bets(client_sock):
    """ Receive a batch of bets from a client socket and return an array of Bet objects"""
    bytes_read = client_sock.recv(BUFFER_SIZE)
    while bytes_read[-2:] != b'\r\n':
        bytes_read += client_sock.recv(BUFFER_SIZE)
    return message.bets_from_bytes(bytes_read)

def send_ack_response(client_sock):
    """ Send an ack response to a client socket"""
    response_message = message.ResponseMessage(SERVER_ID, message.ResponseType.AckMessage)
    _send_response(client_sock, response_message)

def send_err_response(client_sock):
    """ Send an err response to a client socket"""
    response_message = message.ResponseMessage(SERVER_ID, message.ResponseType.ErrMessage)
    _send_response(client_sock, response_message)

def _send_response(client_sock, response_message):
    """ Send a response message to a client socket"""
    bytes_sent = client_sock.send(response_message.to_bytes())
    while bytes_sent < len(response_message.to_bytes()):
        bytes_sent += client_sock.send(response_message.to_bytes()[bytes_sent:])
    