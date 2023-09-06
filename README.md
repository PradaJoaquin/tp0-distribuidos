# Ejercicio 7

## Protocolo
El protocolo en este ejercicio fue modificado para que se efectue el sorteo y cada cliente/agencia pueda saber sus ganadores.

Para eso se agregaron nuevos tipos de mensajes:

- DoneSendingBets: Mensaje que envían las agencias para que el server sepa que ya terminaron de enviar apuestas.
    - `Protocol: <sender_id>:<message_type>\r\n`
- RequestWinners: Mensaje que envían las agencias para que el server les envíe los ganadores de su agencia.
    - `Protocol: <sender_id>:<message_type>\r\n`
- WaitMessage: Mensaje que se envía a los clientes para que esperen a que se efectue el sorteo.
    - `Protocol: <sender_id>:<message_type>\r\n`
- WinnersMessage: Mensaje que se envía a los clientes para que sepan quienes fueron los ganadores de su agencia.
    - `Protocol: <sender_id>:<message_type>:<winner_bet_1>\n<winner_bet_2>\n...\n<winner_bet_n>\n\r\n`

Ahora ya no existe separación entre mensajes del server y mensajes de las agencias, sino que todos los mensajes son iguales y se diferencian por el tipo de mensaje.
Por eso mismo ahora todos los mensajes terminan con los caracteres `\r\n`.

Cuando se recibe un WaitMessage, el cliente debe esperar un cierto tiempo para volver a enviar la solicitud de ganadores. Este tiempo es el loop period configurable en el archivo de configuración de los clientes.

El servidor no enviará los ganadores hasta que todas las agencias hayan enviado el mensaje DoneSendingBets. El numero de agencias es configurable en el archivo de configuración del servidor.

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para levantar server y clientes.
- Ejecutar `make docker-compose-logs` en otra terminal para empezar a ver los logs y verificar que el numero de ganadores sea el correcto.