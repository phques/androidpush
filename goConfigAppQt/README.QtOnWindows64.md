Using Go-QML with Go tip on Windows 64bits,
-> Qt seems to be available only as 32bits  
(apparently 64bits builds are available, cf [go-qml 64 bits][refGoQml64])

## Install Go

**Install binaries for Go >= 1.4 (req to build from source)**  
http://golang.org/doc/install
 
We want to install Go from source ('tip', dev version).  
As a reference, have a look at [Install Go from Source][refFromSrc] and [Set up C tools][refWiki]  

**Install TDM-GCC**  
We _do_ want to install [TDM-GCC](refTDM), we will use it to build Go packages, including CGO stuff,
both 64 and 32bits.  
So select the MinGW-w64/TDM64 version during installation.  
I suggest to select to *not* add to path in the options, it will conflict with Qt's MinGW later.

**Get Go sources ('tip'/dev)**  
git clone https://go.googlesource.com/go gotip  
This will install Go sources into a 'gotip' directory

**Build Go** (64 bits)  
_You need an installed version of Go >= 1.4_  
_Antiviruses can interfere with this, for eg. BitDefender needs to be turned off temporarely_  
We also need to add TDM-GCC to the path to build Go (CGO stuff):  
path=C:\TDM-GCC-64\bin;%path%  
cd gotip\src  
set GOROOT_BOOTSTRAP=c:\go1.4.2  
set CGO_ENABLED=1  
make.bat (or all.bat)  


Add GOROOT = c:\gotip in Windows environment variables  
Also add Go to the path (c:\gotip\bin)  
Close your command prompt, open a new one,  
check that our new go is found: 'go version' :  
>go version devel +75c0566 Sat Apr 11 19:36:19 2015 +0000 windows/amd64

**Build 32bits Go libs**  
path=C:\TDM-GCC-64\bin;%path%  
set GOARCH=386  
set CGO_ENABLED=1  
go install -v std  

_carefull not to copy/paste any spaces at the end of these lines !!!!_

_**Remember to re-enable your antivirus if you turned it off**_

----

### Install Qt
**Download Qt mingw x86**  
[Qt 5 mingw install][refQt]

**Install**  
_choose to **also install the mingw** that comes with it_  
(in Select Components tab: Qt | Tools  | MinGW 4.9.1, in this case)

Make certain that you have both Qt & MinGW (from Qt) in your path  
for Qt: `C:\Qt\Qt5.4.1\5.4\mingw491_32\bin`  
for MinGW: `C:\Qt\Qt5.4.1\Tools\mingw491_32\bin`  

To check that Windows sees MinGW gcc (new command prompt to see path changes) :  
type 'gcc --version', you should see something like:  
> gcc (i686-posix-dwarf-rev2, Built by MinGW-W64 project) 4.9.1  
Copyright (C) 2014 Free Software Foundation, Inc.  
This is free software; see the source for copying conditions.  There is NO  
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  

From now on, we are using the MinGW compiler installed with Qt.

**Install pkg-config lite**  
Go-qml uses pkg-config to get C flags to compile, so we need to install it.  
get [pkg-config lite][refPkgCfgLite]  
decompress into C:\Qt\Qt5.4.1\Tools\mingw491_32  
check that it is found in PATH: 'pkg-config --version'
> 0.28

Add a new Windows enviiroment variable : PKG_CONFIG_PATH  
set it to point to C:\Qt\Qt5.4.1\5.4\mingw491_32\lib\pkgconfig

----
### Install Go-QML

**Get Go-QML**  
(make certain you have a proper GOPATH set 1st! new command prompt if you just added PKG_CONFIG_PATH)  
go get -d gopkg.in/qml.v1  

now build it (we need to build in 32bits, since Qt libs are 32bits):  
set GOARCH=386  
set CGO_ENABLED=1  
set CGO_CFLAGS=-IC:\Qt\Qt5.4.1\5.4\mingw491_32\include  
set CGO_LDFLAGS=-LC:\Qt\Qt5.4.1\5.4\mingw491_32\lib  
go install -v gopkg.in/qml.v1

(again carefull about trailing whitespace here !)

Voila !
You can now build / run the go-qml exmaples or your own go-qml app.  
Note that you need to set the env variables GOARCH as above to build.


[refFromSrc]: http://tip.golang.org/doc/install/source
[refWiki]: https://github.com/golang/go/wiki/InstallFromSource#install-c-tools
[refTDM]: http://tdm-gcc.tdragon.net/
[refQt]: http://download.qt.io/official_releases/qt/5.4/5.4.1/qt-opensource-windows-x86-mingw491_opengl-5.4.1.exe
[refPkgCfgLite]: http://sourceforge.net/projects/pkgconfiglite/files
[refGoQml64]: https://groups.google.com/forum/#!topic/go-qml/S5Vho-XtQyo

