"test UDP socket client, gnu smalltalk"
PackageLoader fileInPackage: 'TCP'.
FileStream fileIn: 'JSON.st'.

!Sockets Datagram methodsFor: 'accessing'!
dataAsString
    "Return a string of the proper length from the data in the datagram"
    | string |
    string := String new: dataSize.
    string 
        replaceFrom: 1
        to: dataSize
        withByteArray: self data
        startingAt: 1.
    ^string
! !


Eval[

| datagram dsock res |
dsock := Sockets DatagramSocket local: '192.168.1.255' port: 4444.
"dsock := Sockets DatagramSocket local: (Sockets SocketAddress anyLocalAddress) port: 4444."
"dsock := Sockets DatagramSocket local: 'localhost' port: 4444."
dsock bufferSize: 1024.

datagram := dsock next.
datagram inspect.
datagram ifNotNil: [
    |json jsonString|
    datagram data size printNl.
    datagram dataSize printNl.
    "jsonString := datagram data asString." "causes probs when ran from inside SCiTE! (extra null bytes I guess)"
    "jsonString := (datagram data copyFrom: 1 to: datagram dataSize) asString."
    jsonString := datagram dataAsString.
    jsonString size printNl.
    jsonString printNl.
    [   json := JSONReader fromJSON: jsonString.
        json inspect.
        (json at: 'pushId') printNl.
        (json at: 'destDirType') printNl.
    ] on: Exception
      do: [:sig| ('Exception ', sig messageText) printNl. sig return]
]

]
