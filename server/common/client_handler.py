import logging
import signal
import socket
import common.communication as communication
import common.message as message
from multiprocessing.connection import Client
from common.bets_handler import BetsHandlerOperations
from common.server_state import ServerStateOperations

class ClientHandler:
    def __init__(self, client_sock, server_state_address, bets_handler_address):
        self.client_sock = client_sock
        self.server_state_address = server_state_address
        self.bets_handler_address = bets_handler_address
        self.running = True
        # Register signal handler for SIGTERM
        signal.signal(signal.SIGTERM, self.__stop)

    def handle_client(self):
        while self.running == True:
            # Receive message from client
            try:
                client_message = communication.receive_client_message(self.client_sock)
                self.handle_client_message(self.client_sock, client_message)
            except OSError as e:
                return
            except ValueError as e:
                logging.error(f"action: receive_message | result: fail | error: {e}")
                communication.send_err_response(self.client_sock)
                self.client_sock.close()
                self.running = False
        self.client_sock.close()

    def handle_client_message(self, client_sock, client_message):
        if client_message.message_type == message.MessageType.BetMessageType:
            # Store bet in database
            self.store_bets(client_message.bets)
            
            communication.send_ack_response(client_sock)
            logging.info(f"action: receive_message_bet | result: success")
            
        elif client_message.message_type == message.MessageType.DoneSendingBets:
            # Add client to the list of clients that finished sending bets
            self.add_client_done(client_message.sender_id)
            communication.send_ack_response(client_sock)
            logging.info(f"action: receive_message_done | result: success")

        elif client_message.message_type == message.MessageType.RequestWinners:
            # Send winners to client
            if not self.is_all_clients_done():
                communication.send_wait_response(client_sock)
            else:
                winners = self.filter_winners(client_message.sender_id)
                communication.send_winners_response(client_sock, winners)
            logging.info(f"action: receive_message_request | result: success")

    def store_bets(self, bets):
        with Client(self.bets_handler_address) as conn:
            conn.send([BetsHandlerOperations.StoreBets, bets])

    def filter_winners(self, client_id):
        with Client(self.bets_handler_address) as conn:
            conn.send([BetsHandlerOperations.FilterWinners, client_id])
            winners = conn.recv()
        return winners
    
    def add_client_done(self, client_id):
        with Client(self.server_state_address) as conn:
            conn.send([ServerStateOperations.AddClientDone, client_id])

    def is_all_clients_done(self):
        with Client(self.server_state_address) as conn:
            conn.send([ServerStateOperations.IsAllClientsDone])
            result = conn.recv()
        return result
    
    def __stop(self, *args):
        """
        Stop server closing the client socket.
        """
        logging.info("action: client_handler_shutdown | result: in_progress")
        self.client_sock.shutdown(socket.SHUT_RDWR)
        self.client_sock.close()
        self.running = False
        logging.info("action: client_handler_shutdown | result: success")
    
    