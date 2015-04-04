// Package go_goInterface is an autogenerated binder stub for package goInterface.
//   gobind -lang=go github.com/phques/androidpush/golib/goInterface
//
// File is generated by gobind. Do not edit.
package go_goInterface

import (
	"github.com/phques/androidpush/golib/goInterface"
	"golang.org/x/mobile/bind/seq"
)

func proxy_Init(out, in *seq.Buffer) {
	// Must be a Go object
	param_param_ref := in.ReadRef()
	param_param := param_param_ref.Get().(*goInterface.InitParam)
	err := goInterface.Init(param_param)
	if err == nil {
		out.WriteUTF16("")
	} else {
		out.WriteUTF16(err.Error())
	}
}

const (
	proxyInitParamDescriptor         = "go.goInterface.InitParam"
	proxyInitParamDevicenameGetCode  = 0x00f
	proxyInitParamDevicenameSetCode  = 0x01f
	proxyInitParamAppFilesDirGetCode = 0x10f
	proxyInitParamAppFilesDirSetCode = 0x11f
	proxyInitParamBooksGetCode       = 0x20f
	proxyInitParamBooksSetCode       = 0x21f
	proxyInitParamDCIMGetCode        = 0x30f
	proxyInitParamDCIMSetCode        = 0x31f
	proxyInitParamDocumentsGetCode   = 0x40f
	proxyInitParamDocumentsSetCode   = 0x41f
	proxyInitParamDownloadsGetCode   = 0x50f
	proxyInitParamDownloadsSetCode   = 0x51f
	proxyInitParamMoviesGetCode      = 0x60f
	proxyInitParamMoviesSetCode      = 0x61f
	proxyInitParamMusicGetCode       = 0x70f
	proxyInitParamMusicSetCode       = 0x71f
	proxyInitParamPicturesGetCode    = 0x80f
	proxyInitParamPicturesSetCode    = 0x81f
)

type proxyInitParam seq.Ref

func proxyInitParamDevicenameSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Devicename = v
}

func proxyInitParamDevicenameGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Devicename
	out.WriteUTF16(v)
}

func proxyInitParamAppFilesDirSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).AppFilesDir = v
}

func proxyInitParamAppFilesDirGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).AppFilesDir
	out.WriteUTF16(v)
}

func proxyInitParamBooksSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Books = v
}

func proxyInitParamBooksGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Books
	out.WriteUTF16(v)
}

func proxyInitParamDCIMSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).DCIM = v
}

func proxyInitParamDCIMGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).DCIM
	out.WriteUTF16(v)
}

func proxyInitParamDocumentsSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Documents = v
}

func proxyInitParamDocumentsGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Documents
	out.WriteUTF16(v)
}

func proxyInitParamDownloadsSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Downloads = v
}

func proxyInitParamDownloadsGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Downloads
	out.WriteUTF16(v)
}

func proxyInitParamMoviesSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Movies = v
}

func proxyInitParamMoviesGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Movies
	out.WriteUTF16(v)
}

func proxyInitParamMusicSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Music = v
}

func proxyInitParamMusicGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Music
	out.WriteUTF16(v)
}

func proxyInitParamPicturesSet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := in.ReadUTF16()
	ref.Get().(*goInterface.InitParam).Pictures = v
}

func proxyInitParamPicturesGet(out, in *seq.Buffer) {
	ref := in.ReadRef()
	v := ref.Get().(*goInterface.InitParam).Pictures
	out.WriteUTF16(v)
}

func init() {
	seq.Register(proxyInitParamDescriptor, proxyInitParamDevicenameSetCode, proxyInitParamDevicenameSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDevicenameGetCode, proxyInitParamDevicenameGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamAppFilesDirSetCode, proxyInitParamAppFilesDirSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamAppFilesDirGetCode, proxyInitParamAppFilesDirGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamBooksSetCode, proxyInitParamBooksSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamBooksGetCode, proxyInitParamBooksGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDCIMSetCode, proxyInitParamDCIMSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDCIMGetCode, proxyInitParamDCIMGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDocumentsSetCode, proxyInitParamDocumentsSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDocumentsGetCode, proxyInitParamDocumentsGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDownloadsSetCode, proxyInitParamDownloadsSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamDownloadsGetCode, proxyInitParamDownloadsGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamMoviesSetCode, proxyInitParamMoviesSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamMoviesGetCode, proxyInitParamMoviesGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamMusicSetCode, proxyInitParamMusicSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamMusicGetCode, proxyInitParamMusicGet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamPicturesSetCode, proxyInitParamPicturesSet)
	seq.Register(proxyInitParamDescriptor, proxyInitParamPicturesGetCode, proxyInitParamPicturesGet)
}

func proxy_NewInitParam(out, in *seq.Buffer) {
	res := goInterface.NewInitParam()
	out.WriteGoRef(res)
}

func proxy_Start(out, in *seq.Buffer) {
	err := goInterface.Start()
	if err == nil {
		out.WriteUTF16("")
	} else {
		out.WriteUTF16(err.Error())
	}
}

func proxy_Stop(out, in *seq.Buffer) {
	err := goInterface.Stop()
	if err == nil {
		out.WriteUTF16("")
	} else {
		out.WriteUTF16(err.Error())
	}
}

func init() {
	seq.Register("goInterface", 1, proxy_Init)
	seq.Register("goInterface", 2, proxy_NewInitParam)
	seq.Register("goInterface", 3, proxy_Start)
	seq.Register("goInterface", 4, proxy_Stop)
}
