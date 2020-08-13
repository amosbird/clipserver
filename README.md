# Extremely UNSAFE clipboard service (use it at your own risk)

ssh -R /tmp/clipserver.sock:/tmp/clipserver.sock ...

server: clipserver -S

client: echo hello world | clipserver
