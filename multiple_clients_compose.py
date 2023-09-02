import sys

CLIENT_EXAMPLES = [{"nombre": "Santiago Lionel", "apellido": "Lorca", "DOCUMENTO": "30904465", "NACIMIENTO": "1999-03-17", "numero": "7547"},
                   {"nombre": "Juan pablo", "apellido": "Sancho", "DOCUMENTO": "36674459", "NACIMIENTO": "1990-05-14", "numero": "1000"},
                   {"nombre": "Mauricio", "apellido": "Ponche", "DOCUMENTO": "41924361", "NACIMIENTO": "2000-01-23", "numero": "4231"},
                   {"nombre": "Lionel Andres", "apellido": "Messi", "DOCUMENTO": "35654839", "NACIMIENTO": "1987-06-24", "numero": "10"},
                   {"nombre": "Ramiro", "apellido": "Pachino", "DOCUMENTO": "42914245", "NACIMIENTO": "2000-09-23", "numero": "9876"}]

def create_file(clients):
    file = open("docker-compose-dev.yaml", "w")
    file.write("version: '3.9'\n")
    file.write("name: tp0\n")
    write_services(file, clients)
    write_networks(file)
    file.close()

def write_services(file, clients):
    file.write("services:\n")
    write_server(file)
    for i in range(1, clients + 1):
        write_client(file, i)

def write_server(file):
    file.write("  server:\n")
    file.write("    container_name: server\n")
    file.write("    image: server:latest\n")
    file.write("    entrypoint: python3 /main.py\n")
    file.write("    environment:\n")
    file.write("      - PYTHONUNBUFFERED=1\n")
    file.write("      - LOGGING_LEVEL=DEBUG\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n")
    file.write("    volumes:\n")
    file.write("      - ./server/config.ini:/config.ini\n")
    file.write("\n")
    file.write("\n")

def write_client(file, i):
    file.write("  client" + str(i) + ":\n")
    file.write("    container_name: client" + str(i) + "\n")
    file.write("    image: client:latest\n")
    file.write("    entrypoint: /client\n")
    write_client_environment(file, i)
    file.write("    networks:\n")
    file.write("      - testing_net\n")
    file.write("    depends_on:\n")
    file.write("      - server\n")
    file.write("    volumes:\n")
    file.write("      - ./client/config.yaml:/config.yaml\n")
    file.write("\n")

def write_client_environment(file, i):
    client_example = CLIENT_EXAMPLES[i % len(CLIENT_EXAMPLES)]

    file.write("    environment:\n")
    file.write("      - CLI_ID=" + str(i) + "\n")
    file.write("      - CLI_LOG_LEVEL=DEBUG\n")
    file.write("      - CLI_NOMBRE=" + client_example["nombre"] + "\n")
    file.write("      - CLI_APELLIDO=" + client_example["apellido"] + "\n")
    file.write("      - CLI_DOCUMENTO=" + client_example["DOCUMENTO"] + "\n")
    file.write("      - CLI_NACIMIENTO=" + client_example["NACIMIENTO"] + "\n")
    file.write("      - CLI_NUMERO=" + client_example["numero"] + "\n")

def write_networks(file):
    file.write("networks:\n")
    file.write("  testing_net:\n")
    file.write("    ipam:\n")
    file.write("      driver: default\n")
    file.write("      config:\n")
    file.write("        - subnet: 172.25.125.0/24\n")

def main():
    #arguments
    if len(sys.argv) != 2:
        print("Please especify how many clients you want to create: python3 multiple_clients_compose.py <number_of_clients>")
        sys.exit(1)
    
    clients = int(sys.argv[1])
    create_file(clients)

if __name__ == '__main__':
    main()