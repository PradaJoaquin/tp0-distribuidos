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
    
def bet_from_bytes(bytes_message):
    """ Parse a bet from a bytes message and return a Bet object

        Protocol: <sender_id>:<first_name>,<last_name>,<document>,<birthdate>,<number>\n
    """
    string_message = bytes_message.decode('utf-8').rstrip()
    sender_id, bet_info = string_message.split(':')
    first_name, last_name, document, birthdate, number = bet_info.split(',')
    return Bet(sender_id, first_name, last_name, document, birthdate, number)