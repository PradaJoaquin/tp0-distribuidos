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