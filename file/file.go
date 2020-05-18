package file

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	permissions  = 0666
	logPrefix    = "github.com/hyperized/bpf/files"
	pattern      = "/dev/bpf*"
	emptyPath    = ""
	using        = "using"
	closing      = "closing"
	noAccess     = " cannot be accessed, skipping"
	noDescriptor = " doesn't have a valid file descriptor"
	noDevice     = "unable to find available bpf device"
	twoLog       = "%s: %s"
	threeLog     = "%s: %s %s"
)

var (
	ErrNoAccess     = errors.New(noAccess)
	ErrNoDescriptor = errors.New(noDescriptor)
	ErrNoDevice     = errors.New(noDevice)
)

type File interface {
	File() *os.File
	Path() string
	FileDescriptor() uintptr
	Close() error

	getFile() (*os.File, error)
	canBeOpened() (bool, error)
	getFileDescriptor() uintptr
	canObtainFileDescriptor() (bool, error)
	hasLogging() bool
}

type bpfDevice struct {
	path           string
	file           *os.File
	fileDescriptor uintptr
	logging        bool
}

// Public interface methods.
func (b *bpfDevice) File() *os.File {
	return b.file
}

func (b *bpfDevice) Path() string {
	return b.path
}

func (b *bpfDevice) FileDescriptor() uintptr {
	return b.fileDescriptor
}

func (b *bpfDevice) Close() error {
	if b.hasLogging() {
		log.Printf("%s: %s %s", logPrefix, closing, b.path)
	}

	return b.file.Close()
}

// Private interface methods.
func (b *bpfDevice) getFile() (*os.File, error) {
	return os.OpenFile(b.path, os.O_RDWR, permissions)
}

func (b *bpfDevice) canBeOpened() (bool, error) {
	if file, err := b.getFile(); err == nil {
		b.file = file
		return true, nil
	}

	return false, fmt.Errorf("%s %w", b.path, ErrNoAccess)
}

func (b *bpfDevice) getFileDescriptor() uintptr {
	return b.file.Fd()
}

func (b *bpfDevice) canObtainFileDescriptor() (bool, error) {
	fd := b.getFileDescriptor()
	if int(fd) != -1 {
		b.fileDescriptor = fd
		return true, nil
	}

	return false, fmt.Errorf("%s %w", b.path, ErrNoDescriptor)
}

func (b *bpfDevice) hasLogging() bool {
	return b.logging
}

// Public methods.
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
		err = ErrNoDevice
	}

	if file.hasLogging() {
		log.Printf(threeLog, logPrefix, using, file.Path())
	}

	return file, err
}

// Private methods.
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
