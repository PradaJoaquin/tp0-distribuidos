from enum import Enum
from common.utils import Bet

class MessageType(Enum):
    AckMessage      = 0
    ErrMessage      = 1
    WaitMessage     = 2
    WinnersMessage  = 3
    BetMessageType  = 4
    DoneSendingBets = 5
    RequestWinners  = 6

class Message:
    def __init__(self, sender_id, message_type):
        self.sender_id = sender_id
        self.message_type = message_type

class BetMessage(Message):
    """ 
    Message that contains an array of bets
    
    Protocol: <sender_id>:<message_type>:<client_bet_1>\n<client_bet_2>\n...\n<client_bet_n>\n\r\n
    """
    def __init__(self, sender_id, bets):
        super().__init__(sender_id, MessageType.BetMessageType)
        self.bets = bets

class WinnersMessage(Message):
    """ 
    Message that contains an array of winners
    
    Protocol: <sender_id>:<message_type>:<client_winner_1>\n<client_winner_2>\n...\n<client_winner_n>\n\r\n
    """
    def __init__(self, sender_id, winners):
        super().__init__(sender_id, MessageType.WinnersMessage)
        self.winners = winners
    
    def to_bytes(self):
        """ Return a bytes message from a WinnersMessage object

            Protocol: <sender_id>:<message_type>:<client_winner_1>\n<client_winner_2>\n...\n<client_winner_n>\n\r\n
        """
        string_response = f"{self.sender_id}:{self.message_type.value}:"
        for winner in self.winners:

            string_response += f"{winner.agency},{winner.first_name},{winner.last_name},{winner.document},{winner.birthdate},{winner.number}\n"
        string_response += "\r\n"
        return string_response.encode('utf-8')

class ResponseMessage:
    def __init__(self, sender_id, message_type):
        self.sender_id = sender_id
        self.message_type = message_type

    def to_bytes(self):
        """ Return a bytes message from a ResponseMessage object

            Protocol: <sender_id>:<message_type>\r\n
        """
        string_response = f"{self.sender_id}:{self.message_type.value}\r\n"
        return string_response.encode('utf-8')

def client_message_from_bytes(bytes_message):
    """ Parse an array of bytes to a message and return it

        Protocol: <sender_id>:<message_type>:<body_opcional>\r\n
    """
    string_message = bytes_message.decode('utf-8').rstrip('\r\n')
    split = string_message.split(':')
    sender_id, message_type = split[:2]
    message_type = MessageType(int(message_type))
    if message_type == MessageType.BetMessageType:
        bets = []
        bets_info = split[2]
        for bet in bets_info.split('\n'):
            slipt = bet.split(',')
            # Remove leading and trailing whitespaces
            map(str.strip, slipt)
            first_name, last_name, document, birthdate, number = slipt
            bets.append(Bet(sender_id, first_name, last_name, document, birthdate, number))
        return BetMessage(sender_id, bets)
    else:
        return Message(sender_id, message_type)
