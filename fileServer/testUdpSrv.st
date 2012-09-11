"gnu smalltalk test, UDP socket receive"
Socket newUDP self setPort: 4444.

| res buf nbRead delay cpt |
delay := Delay forMilliseconds: 200.
cpt := 50 * 15.
buf := ByteArray new: 1024.
[	Transcript show:  cpt.
	res := self receiveUDPDataInto: buf.
	nbRead := res at: 1.
	cpt := cpt - 1.
	(nbRead isZero) & (cpt > 0)
] whileTrue: [delay wait]
