"test UDP socket client, gnu smalltalk"
PackageLoader fileInPackage: 'TCP'.
FileStream fileIn: 'JSON.st'.

Eval[

| datagram dsock res json jsonString |
dsock := Sockets DatagramSocket local: '192.168.1.255' port: 4444.
dsock bufferSize: 1024.
"dsock := Sockets DatagramSocket local: (Sockets SocketAddress anyLocalAddress) port: 4444."
"dsock := Sockets DatagramSocket local: 'localhost' port: 4444."

datagram := dsock next.
datagram inspect.
datagram ifNotNil: [
    jsonString := datagram data asString.
    jsonString printNl.
    json := JSONReader fromJSON: jsonString.
    '??' printNl.
    "json inspect."
    (json at: 'pushId') printNl.
]
]
