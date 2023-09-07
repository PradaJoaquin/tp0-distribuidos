# Ejercicio 1

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal
- Verificar que se crearon 2 clientes ejecutando `make docker-compose-logs`

# Ejercicio 1.1

## Ejecución
- Ejecutar `python3 multiple_clients_compose.py <number_of_clients>` en la terminal reemplazando *<number_of_clients>* por la cantidad de clientes que se desee crear.
- Verificar que se crearon la cantidad correcta de clientes ejecutando `make docker-compose-logs`.

# Ejercicio 2

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal.
- Luego de que se terminen de ejecutar los clientes modificar el archivo de config del cliente en `client/config.yaml`.
- Verificar que al volve a ejecutar `make docker-compose-up` no se creó ninguna imagen nueva y se utilizan los nuevos cambios.

# Ejercicio 3

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para que este levantado el server.
- Ejecutar `bash netcat_test.sh` en la terminal y si el resultado es *Netcat test passed* el servidor esta funcionando.
- **Atencion**: Si se ejecuta desde windows es necesario cambiar los 2 archivos bash al EFL LF.

# Ejercicio 4

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para levantar server y clientes.
- Ejecutar `make docker-compose-logs` en otra terminal para empezar a ver los logs.
- Ejecutar `make docker-compose-down` en la primera terminal mientras se ejecutan los clientes y servidor y observar los logs para ver el graceful shutdown.

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

# Ejercicio 6

## Protocolo
El protocolo en este ejercicio fue modificado para que los clientes puedan enviar las apuestas de a *batches*.

Para eso se modificó el *BetMessage* para poder enviar más de una apuesta en un mismo mensaje. El formato del mensaje es el siguiente:

`<sender_id>:<bet1>\n<bet2>\n,...,<betN>\n\r\n`

Donde:
- `<sender_id>` es el id del cliente que envía el mensaje, o agencia.
- `<bet1>`, `<bet2>`, ..., `<betN>` son las apuestas que envía el cliente, separadas por saltos de línea.
- El mensaje termina con un `\r\n`, que indica que no hay más apuestas en el *batch*.

El *batch size* es configurable en el archivo de configuración. Con un *batch size* de 125 cada *BetMessage* tiene un tamaño promedio de 6KB.

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para levantar server y clientes.
- Ejecutar `make docker-compose-logs` en otra terminal para empezar a ver los logs y verificar que las apuestas se estén procesando correctamente.

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

# Ejercicio 8

## Paralelismo
Para este ejercicio pensé en distintas soluciones sobre como lograr la concurrencia/paralelismo.

- La primera fue crear un thread por cada cliente, pero esto lo terminé descartando porque no me pareció una buena solución teniendo en cuenta el global interpreter lock de python, el cual no permite que dos threads ejecuten código python al mismo tiempo.
- La segunda fue crear un proceso por cada cliente y pasarle el socket del cliente, pero el problema que tenía era que no había forma de volver a unir los procesos al proceso padre y por lo tanto se creaba un memory leek de procesos "zombies".
    - Para solucionar estó pense en otras soluciones:
        - La primera era hacer una pool de procesos donde de una cola agarren los trabajos y ejecuten hasta terminar y agarrar otro, esta solución tampoco me convenció ya que me parecía compleja para el problema a resolver y tampoco me gustaba que al crear el pool le tenga que especificar un número de workers.
        - Y de ahí me llevo a la última solución que es simplemente llamar a la función [multiprocessing.active_children()](https://docs.python.org/3/library/multiprocessing.html#multiprocessing.active_children) que resulve el problema de los zombies ya que como dice su documentación *"Calling this has the side effect of “joining” any processes which have already finished."* Entonces cuando llega una solicitud de conexión nueva primero se liberan todos los procesos que ya terminaron y luego se crea el proceso nuevo con la nueva conexión, de esta manera no desperdiciando recursos y logrando una implementación sencilla. Esta implementacion tiene como limitación que no hay límite para la creación de procesos, cosa que puede ser un problema, para eso sería mejor usar la pool o una combinación de ambas. Pero para resolver este problema decidí ir por este camino no tan complejo.

## Mecanismos de sincronización
Primero como ya aclaramos antes a partir de que llega una nueva conexión de un cliente se crea un nuevo proceso que se va a encargar de manejar el socket de la conexión, el cliente fue modificado para que no se cree una nueva conexión por cada petición, sino manteniendola hasta finalizar.

Para los mecanismos de sincronización se crearon 2 nuevas entidades:
- El `BetsHandler` encargado de ocuparse del gaurdado y obtención de ganadores del archivo de bets.
- El `ServerState` encargado de mantener una cuenta de cuantos clientes entregaron ya todas sus bets.

Ambas entidades corren en un proceso aparte y su comunicación se hace a través de `multiprocessing.connection` simulando una conexión de sockets a nivel local. De esta manera se asegura que solo un proceso a la vez pueda modificar algunas de las secciones criticas.

El protocolo implementado para comunicarse con estas nuevas entidades es mandandole una lista con el primer valor siendo la operación que se quiere efectuar, en forma de un ENUM, y el segundo valor lo que se quiere agregar (bets en el caso de el BetsHandler).

## Ejecución
- Ejecutar `make docker-compose-up` en la terminal para levantar server y clientes.
- Ejecutar `make docker-compose-logs` en otra terminal para empezar a ver los logs y verificar que el numero de ganadores sea el correcto.
