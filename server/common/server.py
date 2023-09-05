import socket
import logging
import signal
import common.communication as communication
import common.utils as utils

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self.running = True
        # Register signal handler for SIGTERM
        signal.signal(signal.SIGTERM, self.__stop)


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self.running:
            try:
                client_sock = self.__accept_new_connection()
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
            client_bets = communication.receive_client_bets(client_sock)
            
            # Store bet in database
            utils.store_bets(client_bets)
            
            # Send ack response
            communication.send_ack_response(client_sock)
            logging.info(f"action: receive_message | result: success")
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            communication.send_err_response(client_sock)
        except ValueError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            communication.send_err_response(client_sock)
        finally:
            client_sock.close()

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

