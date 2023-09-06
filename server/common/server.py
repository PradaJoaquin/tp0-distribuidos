import socket
import logging
import signal
import common.communication as communication
import common.message as message
import common.utils as utils

class Server:
    def __init__(self, port, listen_backlog, number_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.clients_done = set()
        self.total_number_of_clients = number_of_clients

        self.running = True
        # Register signal handler for SIGTERM
        signal.signal(signal.SIGTERM, self.__stop)


    def run(self, connection_timeout=0):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again

        To set a timeout for the accept operation, set the connection_timeout parameter
        """
        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                client_sock.settimeout(connection_timeout)
                self.__handle_client_connection(client_sock)
            except OSError:
                return

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # Receive bet from client
            client_message = communication.receive_client_message(client_sock)
            if client_message.message_type == message.MessageType.BetMessageType:
                # Store bet in database
                utils.store_bets(client_message.bets)
                communication.send_ack_response(client_sock)
                logging.info(f"action: receive_message_bet | result: success")
                
            elif client_message.message_type == message.MessageType.DoneSendingBets:
                # Add client to the list of clients that finished sending bets
                self.clients_done.add(client_message.sender_id)
                communication.send_ack_response(client_sock)
                logging.info(f"action: receive_message_done | result: success")

            elif client_message.message_type == message.MessageType.RequestWinners:
                # Send winners to client
                if len(self.clients_done) < self.total_number_of_clients:
                    communication.send_wait_response(client_sock)
                else:
                    winners = self.filter_winners(client_message.sender_id)
                    communication.send_winners_response(client_sock, winners)
                logging.info(f"action: receive_message_request | result: success")
                
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            communication.send_err_response(client_sock)
        except ValueError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            communication.send_err_response(client_sock)
        finally:
            client_sock.close()

    def filter_winners(self, client_id):
        """
        Filter winners by client id

        winners: list of Bet objects
        client_id: client id to filter winners
        """
        bets_generator = utils.load_bets()
        winners = [bet for bet in bets_generator if utils.has_won(bet)]
        return [winner for winner in winners if winner.agency == int(client_id)]

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def __stop(self, *args):
        """
        Stop server closing the server socket.
        """
        logging.info("action: server_shutdown | result: in_progress")
        self.running = False
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info("action: server_socket_closed | result: success")
        logging.info("action: server_shutdown | result: success")

