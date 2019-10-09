package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

const permissions = 0666
const logPrefix = "github.com/hyperized/bpf/files"
const pattern = "/dev/bpf*"
const emptyPath = ""
const using = "using"
const closing = "closing"
const noAccess = " cannot be accessed, skipping"
const noDescriptor = " doesn't have a valid file descriptor"
const noDevice = "unable to find available bpf device"
const twoLog = "%s: %s"
const threeLog = "%s: %s %s"

type File interface {
	File() (file *os.File)
	Path() (path string)
	FileDescriptor() (fd uintptr)
	Close() (err error)

	getFile() (*os.File, error)
	canBeOpened() (ok bool, err error)
	getFileDescriptor() (fd uintptr)
	canObtainFileDescriptor() (ok bool, err error)
	hasLogging() (ok bool)
}

type bpfDevice struct {
	path           string
	file           *os.File
	fileDescriptor uintptr
	logging        bool
}

// Public interface methods
func (b *bpfDevice) File() (file *os.File) {
	return b.file
}

func (b *bpfDevice) Path() (path string) {
	return b.path
}

func (b *bpfDevice) FileDescriptor() (fd uintptr) {
	return b.fileDescriptor
}

func (b *bpfDevice) Close() (err error) {
	if b.hasLogging() {
		log.Printf("%s: %s %s", logPrefix, closing, b.path)
	}
	return b.file.Close()
}

// Private interface methods
func (b *bpfDevice) getFile() (*os.File, error) {
	return os.OpenFile(b.path, os.O_RDWR, permissions)
}

func (b *bpfDevice) canBeOpened() (ok bool, err error) {
	if file, err := b.getFile(); err == nil {
		ok = true
		b.file = file
		return true, nil
	}
	return false, errors.New(b.path + noAccess)
}

func (b *bpfDevice) getFileDescriptor() (fd uintptr) {
	return b.file.Fd()
}

func (b *bpfDevice) canObtainFileDescriptor() (ok bool, err error) {
	fd := b.getFileDescriptor()
	if int(fd) != -1 {
		b.fileDescriptor = fd
		return true, nil
	}
	return false, errors.New(b.path + noDescriptor)
}

func (b *bpfDevice) hasLogging() (ok bool) {
	return b.logging
}

// public methods
func GetBpfDevice(logging bool) (file File, err error) {
	var paths []string
	file = newBpfDevice()
	paths, err = listBpfDevices()

	for _, path := range paths {
		file = &bpfDevice{
			logging: logging,
			path:    path,
		}
		if ok, err := file.canBeOpened(); !ok {
			myLog(err, file)
			continue
		}
		if ok, err := file.canObtainFileDescriptor(); !ok {
			myLog(err, file)
			continue
		}
		break
	}

	if file.Path() == emptyPath {
		err = errors.New(noDevice)
	}

	if file.hasLogging() {
		log.Printf(threeLog, logPrefix, using, file.Path())
	}

	return file, err
}

// private methods
func newBpfDevice() *bpfDevice {
	return &bpfDevice{}
}

func listBpfDevices() (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	return
}

func myLog(err error, file File) {
	if err != nil && file.hasLogging() {
		log.Printf(twoLog, logPrefix, err)
	}
}
