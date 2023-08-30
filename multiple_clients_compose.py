import sys

def create_file(clients):
    with open("docker-compose-dev.yaml", "w") as file:
        file.write("version: '3.9'\n")
        file.write("name: tp0\n")
        write_services(file, clients)
        write_networks(file)

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
    file.write("\n")

def write_client(file, i):
    file.write("  client" + str(i) + ":\n")
    file.write("    container_name: client" + str(i) + "\n")
    file.write("    image: client:latest\n")
    file.write("    entrypoint: /client\n")
    file.write("    environment:\n")
    file.write("      - CLI_ID=" + str(i) + "\n")
    file.write("      - CLI_LOG_LEVEL=DEBUG\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n")
    file.write("    depends_on:\n")
    file.write("      - server\n")
    file.write("\n")

def write_networks(file):
    file.write("networks:\n")
    file.write("  testing_net:\n")
    file.write("    ipam:\n")
    file.write("      driver: default\n")
    file.write("      config:\n")
    file.write("        - subnet: 172.25.125.0/24\n")

def main():
    if len(sys.argv) != 2:
        print("Please especify how many clients you want to create: python3 multiple_clients_compose.py <number_of_clients>")
        sys.exit(1)
    try:
        clients = int(sys.argv[1])
    except ValueError:
        print("Please especify how many clients you want to create with a number: python3 multiple_clients_compose.py <number_of_clients>")
        sys.exit(1)
    create_file(clients)

if __name__ == '__main__':
    main()