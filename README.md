# Ejercicio 5

## Protocolo
El protocolo implementado es el siguiente:

Hay 2 tipos de mensajes, los de tipo *BetMessage* y los *ResponseMessage*.
Cuando los clientes se conectan al servidor, el cliente envía un mensaje de tipo *BetMessage* con los valores necesarios en el siguiente formato:

`<sender_id>:<first_name>,<last_name>,<document>,<birthdate>,<number>\n`

Donde:
- `<sender_id>` es el id del cliente que envía el mensaje, o agencia.
- Y luego de los dos puntos, se envían los datos del cliente separados por comas.
- El mensaje termina con un salto de línea.

El servidor responde con un mensaje de tipo *ResponseMessage* con los valores necesarios en el siguiente formato:

`<sender_id>:<response_type>\n`

Donde:
- `<sender_id>` es el id del servidor que envía el mensaje.
- `<response_type>` es el tipo de respuesta que envía el servidor, puede ser `ACK` si salió todo bien o `ERROR` si hubo algún error en el mensaje enviado por el cliente.

Para asegurarnos de que no haya *short reads* o *short writes*, se lee el socket hasta que se reciba un salto de línea, y se escribe hasta que se hayan escrito todos los bytes del mensaje.


## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para levantar server y clientes.
- Ejecutar `make docker-compose-logs` en otra terminal para empezar a ver los logs y verificar que las apuestas se estén procesando correctamente.