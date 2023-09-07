import socket
import logging
import signal
import multiprocessing
import common.client_handler as client_handler
from common.bets_handler import BetsHandler, BETS_HANDLER_ADDRESS
from common.server_state import ServerState
from common.server_state import SERVER_STATE_ADDRESS

class Server:
    def __init__(self, port, listen_backlog, number_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self.state = ServerState(number_of_clients, SERVER_STATE_ADDRESS)
        self.bets_handler = BetsHandler(BETS_HANDLER_ADDRESS)

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
        # Start bets handler
        multiprocessing.Process(target=self.bets_handler.listen).start()
        # Start server state
        multiprocessing.Process(target=self.state.listen).start()

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                client_sock.settimeout(connection_timeout)
                # Active children call is necessary to avoid zombie processes
                multiprocessing.active_children()
                
                c_handler = client_handler.ClientHandler(SERVER_STATE_ADDRESS, BETS_HANDLER_ADDRESS)
                multiprocessing.Process(target=c_handler.handle_client, args=(client_sock,)).start()
            except OSError:
                return

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

