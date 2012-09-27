// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
'marco polo': 
my own mini zero-configuration networking (just two 'hardcoded' udp ports ;-p). 

PC calls 'marco', android answers 'polo' on UDP, 
this way the PC can find the Android IP.

this version only works with polo launched first:
Android/Polo waits for udp datagram 'marco' on port 4444
PC/Polo bradcasts 'marco' on UDP to port 4444
polo recvs 'marco' and answers 'polo' to PC using the address found in te datagram (ie we know who called us)
marco recvs 'polo' from Android, thus we know the IP address of the Android device.

next version will work no matter in which order the 2 apps are launched: 

PC broadcasts 'marco' to 4444 (on 'anyPort')
PC asynch recv on 4445
PC asynch recv on 'anyPort' (from send/broadcast above)
=> recv from either a 'polo' response, or from a "I'm here" 'polo' message

Android broadcast 'polo' to 4445 
Android recvs on 4444 and answers to 'polo' to sender address

