from enum import Enum
from common.utils import Bet

class ResponseType(Enum):
    AckMessage = 0
    ErrMessage = 1

class ResponseMessage:
    def __init__(self, sender_id, response_type):
        self.sender_id = sender_id
        self.response_type = response_type

    def to_bytes(self):
        """ Return a bytes message from a ResponseMessage object

            Protocol: <sender_id>:<response_type>\n
        """
        string_response = f"{self.sender_id}:{self.response_type.value}\n"
        return string_response.encode('utf-8')
    
def bets_from_bytes(bytes_message):
    """ Parse an array of bets from a bytes message and return an array of Bet objects

        Protocol: <sender_id>:<first_name>,<last_name>,<document>,<birthdate>,<number>\n
    """
    string_message = bytes_message.decode('utf-8').rstrip('\r\n')
    sender_id, bets_info = string_message.split(':')

    bets = []
    for bet in bets_info.split('\n'):
        slipt = bet.split(',')
        # Remove leading and trailing whitespaces
        map(str.strip, slipt)
        first_name, last_name, document, birthdate, number = slipt
        bets.append(Bet(sender_id, first_name, last_name, document, birthdate, number))
    return bets