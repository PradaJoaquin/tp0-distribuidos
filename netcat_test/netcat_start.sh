if [ $(echo "test" | nc -q 1 server 12345) = "test" ]; then
    echo "Netcat test passed"
else
    echo "Netcat test failed"
    exit 1
fi