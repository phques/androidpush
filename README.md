AndroidPush project 
Copyright 2014 Philippe Quesnel
is Licensed under the Academic Free License version 3.0  
see License.txt or License.html 


----


#### AndroidPush 

PC & Android apps to 'push' (share) files between Android devices and PC,  
through *local* network (ie- normally, wi-fi)  
Using known 'public' directories on Android (ie Downloads, Music, Pictures etc)

Most of the functionality will be written in Go (golang).  
This makes it useable on Linux,Windows,Android.

Typical work-flow:
* Select files (eg some pictures) to 'push' on a device, 'Share with AndroidPush'
* it will determine the root/base dir (in this case the one corresponding to 'Pictures')
* it will locate other running instances of AndroidPush & ask user to select if more than one
* a 'push request' will be sent to other instance, which will check for common root/base dir  
  check for overwritting files etc.
* then it will pull (ah! andoirdPull then !! ;-p) the files from the originating androidPush

Android (Java) app will access the go functionality through a .so lib (see [golang.org/x/mobile](https://github.com/golang/mobile/))

Also using marcoPoloPQ simple service discovery [github.com/phques/mppq](https://github.com/phques/mppq)

-----

2013-07-08 : note that this idea came from the fact that I could not use my Galaxy S3 with MTP (USB cable) on Linux !  
the current Linux version I am using now supports this.  
So this project lost some of it's original purpose ;-)

Philippe Quesnel  

