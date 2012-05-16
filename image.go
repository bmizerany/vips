package vips

/*
#cgo pkg-config: vips-7.28
#include <stdlib.h>
#include <vips/vips.h>
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

func Init(argv0 string) {
	cargv0 := C.CString(argv0)
	defer C.free(unsafe.Pointer(cargv0))
	if errc := int(C.im_init_world(cargv0)); errc != 0 {
		panic("unable to start vips")
	}
}

type Image struct {
	im *C.IMAGE
}

func Open(file string, mode string) (*Image, error) {
	cfile := C.CString(file)
	cmode := C.CString(mode)
	defer func() {
		C.free(unsafe.Pointer(cfile))
		C.free(unsafe.Pointer(cmode))
	}()
	if im := C.im_open(cfile, cmode); im != nil {
		img := &Image{im}
		runtime.SetFinalizer(img, (*Image).Close)
		return img, nil
	}
	return nil, errors.New("unable to open: " + file)
}

func (img *Image) Avg() (float64, error) {
	cavg := C.double(0)
	if C.im_avg(img.im, &cavg) != 0 {
		return 0, errors.New("unable to find average")
	}
	return float64(cavg), nil

}

func (img *Image) Close() {
	if img.im != nil {
		C.im_close(img.im)
		img.im = nil
		runtime.SetFinalizer(img, nil)
	}
}
