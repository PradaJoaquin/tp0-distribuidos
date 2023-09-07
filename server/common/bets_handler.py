from enum import Enum
import common.utils as utils
from multiprocessing.connection import Listener

BETS_HANDLER_ADDRESS = ('localhost', 6000)

class BetsHandlerOperations(Enum):
    """
    Enum that contains the operations that can be performed by the BetsHandler
    """
    StoreBets = 0
    FilterWinners = 1
    Close = 2

class BetsHandler:
    """
    Thread-Safe/Process-Safe Bets hanlder to access the bets information.
    """
    def __init__(self, address):
        self.address = address
        self.listener = Listener(self.address)

    def listen(self):
        """
        Listen for bets
        """
        while True:
            conn = self.listener.accept()
            msg = conn.recv()
            operation = msg[0]
            if operation == BetsHandlerOperations.StoreBets:
                self.store_bets(msg[1])
            elif operation == BetsHandlerOperations.FilterWinners:
                conn.send(self.filter_winners(msg[1]))
            elif operation == BetsHandlerOperations.Close:
                break
        self.listener.close()

    def store_bets(self, bets):
        """
        Store bets in "file"
        """
        utils.store_bets(bets)
    
    def filter_winners(self, client_id):
        """
        Filter winners by client id / agency

        return: list of Bet objects that won the lottery
        """
        bets = utils.load_bets()
        winners = [bet for bet in bets if utils.has_won(bet)]
        return [winner for winner in winners if winner.agency == int(client_id)]
