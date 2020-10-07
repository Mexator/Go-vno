package fileserver

import (
	"context"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	syscall "golang.org/x/sys/unix"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	"google.golang.org/grpc/peer"
)

// FileServer implements FileServerServer
type FileServer struct {
	// Absolute path to folder
	storagePath string
	// For authentication
	nsUrl      string
	storageDir *os.File
}

const (
	serverCatalogPerms os.FileMode = 0777
)

var (
	ErrNotNS = errors.New("Not a name server")
)

// Check if server catalog exists and has proper access rights
func initializeServerCatalog(path string) error {
	dir, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "Can not open server catalog:")
	}

	if os.IsNotExist(err) {
		log.Println("Server catalog not found. Creating one in " + path)
		err := os.Mkdir(path, serverCatalogPerms)
		if err != nil {
			return errors.Wrap(err, "Can not create server catalog")
		}
		// Because of fucking umask
		os.Chmod(path, serverCatalogPerms)
	}

	dir, err = os.Open(path)
	info, err := dir.Stat()
	if err != nil {
		return err
	}
	isPermissionsInvalid := (info.Mode().Perm() != serverCatalogPerms)

	if !info.IsDir() || isPermissionsInvalid {
		fmt := "`%s' has invalid permissions or is not a directory"
		return errors.Errorf(fmt, path)
	}

	return err
}

/*
 * MakeFileServer creates FileServer reading its configuration from
 * `config.json` file.
 */
func MakeFileServer(storagePath, nsUrl string) (FileServer, error) {
	err := initializeServerCatalog(storagePath)
	if err != nil {
		return FileServer{}, errors.Wrap(err,
			"Server catalog can not be initialized")
	}

	file, _ := os.Open(storagePath)
	nsUrl = nsUrl[:strings.LastIndex(nsUrl, ":")]
	return FileServer{storagePath, nsUrl, file}, nil
}

func (server *FileServer) isNS(ctx context.Context) bool {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return false
	}
	addr := peer.Addr.String()
	addr = addr[:strings.LastIndex(addr, ":")]
	return addr == server.nsUrl
}

// Size returns size of file in request, or error
func (server *FileServer) Size(
	ctx context.Context,
	request *api.SizeRequest,
) (*api.SizeResponse, error) {
	filePath := path.Join(server.storagePath, request.Inode)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &api.SizeResponse{Size: uint64(fileInfo.Size())}, nil
}

// Read reads a file on server
func (server *FileServer) Read(
	ctx context.Context,
	request *api.ReadRequest,
) (*api.ReadResponse, error) {
	filePath := path.Join(server.storagePath, request.Inode)

	file, err := os.Open(filePath)
	if err != nil {
		return &api.ReadResponse{Content: nil}, errors.New("Fragment not exist")
	}

	defer file.Close()

	buffer := make([]byte, request.Size)

	len, err := file.ReadAt(buffer, int64(request.Offset))
	if len > 0 && (err == nil || err == io.EOF) {
		return &api.ReadResponse{Content: buffer[:len]}, nil
	}
	return &api.ReadResponse{Content: nil}, errors.New("Error reading fragment")

}

func (server *FileServer) Write(
	ctx context.Context,
	request *api.WriteRequest,
) (*api.WriteResponse, error) {
	filePath := path.Join(server.storagePath, request.Inode)
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return nil, errors.New("Can not open fragment for writing")
	}

	defer file.Close()
	_, err = file.WriteAt(request.Content, int64(request.Offset))
	return &api.WriteResponse{}, err
}

// Create creates file or reports an error
func (server *FileServer) Create(
	ctx context.Context,
	request *api.CreateRequest,
) (*api.CreateResponse, error) {
	if !server.isNS(ctx) {
		return nil, errors.Wrap(ErrNotNS, "Failed to create")
	}
	filePath := path.Join(server.storagePath, request.Inode)

	_, err := os.Stat(filePath)
	// Check that fragment is not exists
	if os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(filePath)
		defer file.Close()

		if err != nil {
			return nil, errors.Wrap(err, "Can not create fragment")
		}
		return &api.CreateResponse{}, nil
	}
	return nil, errors.New("Fragment exists")
}

// Remove removes file or reports an error
func (server *FileServer) Remove(
	ctx context.Context,
	request *api.RemoveRequest,
) (*api.RemoveResponse, error) {
	if !server.isNS(ctx) {
		return nil, errors.Wrap(ErrNotNS, "Failed to remove")
	}
	filePath := path.Join(server.storagePath, request.Inode)
	err := os.Remove(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Can not remove fragment")
	}
	return &api.RemoveResponse{}, nil
}

// ReportDiskSpace generates a brief report about used space on disk
func (server *FileServer) ReportDiskSpace(
	ctx context.Context,
	_ *api.Empty,
) (*api.DiskSpaceResponse, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(server.storagePath, &stat)
	if err != nil {
		return nil, err
	}
	resp := &api.DiskSpaceResponse{
		FreeBlocks:     int64(stat.Bavail),
		BusyBlocks:     int64(stat.Blocks - stat.Bavail),
		BlockSizeBytes: stat.Bsize,
	}
	return resp, nil
}
