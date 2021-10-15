# Natsgen it is a publisher/subscriber implimentation with "NATS" message broker written in GO.
Simple project written to try NATS message broker.  

Can be used as a regular GO module 
`go get github.com/kdimtri/natsgen && go install $_`
Or get a copy with `git clone https://github.com/kdimtri/natsgen` and 
build it with `go build`

To run a  publisher use: `natsgen pub pub.json`
This will run a publisher that uses "default" subject and nats 'demo'  server.
Publisher will post messages with interval of 1 second.
Messages prints in logs and in new file `pub.json'.
On connection loss, It will  try  to reconnect to server.
Last message from file 'pub.json' will be used as a starting point for next runs of  publisher.

For a subscriber run: `natsgen sub sub.json`
This will run a subscriber that uses the same "default" subject and nats 'demo' server.
It gets messages from broker and  stores them in its own file 'sub.json'

Get a  list of additional run options with: `natsgen -h`.
