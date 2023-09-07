from enum import Enum
import logging
import multiprocessing
from multiprocessing.connection import Listener
import signal

SERVER_STATE_ADDRESS = ('localhost', 6001)

class ServerStateOperations(Enum):
    """
    Enum that contains the operations that can be performed by the ServerState
    """
    AddClientDone = 0
    IsAllClientsDone = 1
    Close = 2

class ServerState:
    """
    Thread-Safe/Process-Safe Server status containing information about server state.
    """
    def __init__(self, number_of_clients, address):
        self.clients_done = set()
        self.total_number_of_clients = number_of_clients
        self.address = address
        self.listener = Listener(self.address)
        self.running = True
        # Register signal handler for SIGTERM
        signal.signal(signal.SIGTERM, self.__stop)

    def listen(self):
        """
        Listen for server state operations
        """
        while self.running == True:
            try:
                conn = self.listener.accept()
                msg = conn.recv()
                operation = msg[0]
                if operation == ServerStateOperations.AddClientDone:
                    self.add_client_done(msg[1])
                elif operation == ServerStateOperations.IsAllClientsDone:
                    conn.send(self.is_all_clients_done())
                elif operation == ServerStateOperations.Close:
                    break
            except OSError:
                return
        self.listener.close()

    def add_client_done(self, client_id):
        """
        Add client to the list of clients that finished sending bets
        """
        self.clients_done.add(client_id)

    def is_all_clients_done(self):
        """
        Check if all clients have finished sending bets
        """
        result = len(self.clients_done) == self.total_number_of_clients
        return result
    
    def __stop(self, signum, frame):
        """
        Stop the bets handler
        """
        logging.info("action: server_state_shutdown | result: in_progress")
        self.running = False
        self.listener.close()
        logging.info("action: server_state_shutdown | result: success")
