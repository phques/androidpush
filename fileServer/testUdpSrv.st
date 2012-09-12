"pharo smalltalk test, UDP socket receive"
| sock res buf nbRead delay cpt |
sock := Socket newUDP setPort: 4444.
delay := Delay forMilliseconds: 200.
cpt := 1000 / (delay delayDuration) * 15.
buf := ByteArray new: 1024.
[   "Transcript show:  cpt; flush."
    FileStream stdout print: cpt; crlf.
    res := sock receiveUDPDataInto: buf.
    "Transcript cr; flush."
    nbRead := res at: 1.
    cpt := cpt - 1.
    (nbRead isZero) & (cpt > 0)
] whileTrue: [delay wait].
res inspect.
buf inspect.
sock closeAndDestroy

