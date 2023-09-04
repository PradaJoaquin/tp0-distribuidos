import common.message as message

# TODO: if more servers are added, this should be in the config file.
SERVER_ID = "server"


def receive_client_bet(client_sock):
    """ Receive a bet from a client socket and return a Bet object"""
    bytes_read = client_sock.recv(1024)
    while bytes_read[-1:] != b'\n':
        bytes_read += client_sock.recv(1024)
    return message.bet_from_bytes(bytes_read)

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
    