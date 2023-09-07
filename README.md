# Ejercicio 7

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