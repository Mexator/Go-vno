package fileserver

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path"

	syscall "golang.org/x/sys/unix"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	"github.com/Mexator/Go-vno/pkg/config"
)

// FileServer implements FileServerServer
type FileServer struct {
	// Absolute path to folder
	storagePath string
	storageDir  *os.File
}

type fileServerConfig struct {
	FilesDir string `json:"files_dir"`
}

// Check if server dir exists and has proper access rights
func initializeServerCatalog(path string) error {
	dir, err := os.Open(path)
	if err == nil {
		info, _ := dir.Stat()
		isNotDir := !info.IsDir()
		validPermissions := (info.Mode().Perm() == 0777)

		if isNotDir || !validPermissions {
			return errors.New("Server catalog can not be initialized: invalid" +
				" permissions")
		}
	}

	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		os.Chmod(path, 0777)
		return err
	}
	return err
}

/*
MakeFileServer creates FileServer reading its configuration from
`config.json` file.
*/
func MakeFileServer(configFilename string) (FileServer, error) {
	conf := new(fileServerConfig)
	err := config.ReadConfig(&conf, configFilename)

	if err == nil {
		err = initializeServerCatalog(conf.FilesDir)

		if err == nil {
			file, _ := os.Open(conf.FilesDir)
			return FileServer{conf.FilesDir, file}, nil
		}
	}
	return FileServer{}, err
}

// Size returns size of file in request, or error
func (server FileServer) Size(ctx context.Context,
	request *api.SizeRequest) (*api.SizeResponse, error) {

	filePath := path.Join(server.storagePath, request.Inode)
	fileInfo, err := os.Stat(filePath)

	if err == nil {
		log.Print(api.SizeResponse{Size: uint64(fileInfo.Size())})
		return &api.SizeResponse{Size: uint64(fileInfo.Size())}, nil
	}

	return nil, err
}

// Read reads a file on server
func (server FileServer) Read(ctx context.Context,
	request *api.ReadRequest) (*api.ReadResponse, error) {

	filePath := path.Join(server.storagePath, request.Inode)
	file, err := os.Open(filePath)

	if err == nil {
		defer file.Close()
		buffer := make([]byte, 0, request.Size)
		var len int
		len, err = file.ReadAt(buffer, int64(request.Offset))
		if len > 0 && err == io.EOF {
			return &api.ReadResponse{Content: buffer}, nil
		}
	}
	return &api.ReadResponse{Content: nil}, err
}

func (server FileServer) Write(ctx context.Context,
	request *api.WriteRequest) (*api.WriteResponse, error) {

	filePath := path.Join(server.storagePath, request.Inode)
	file, err := os.Open(filePath)

	if err == nil {
		defer file.Close()
		_, err = file.WriteAt(request.Content, int64(request.Offset))
		return &api.WriteResponse{}, nil
	}
	return nil, err
}

// Create creates file or reports an error
func (server FileServer) Create(ctx context.Context,
	request *api.CreateRequest) (*api.CreateResponse, error) {

	filePath := path.Join(server.storagePath, request.Inode)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(filePath)
		defer file.Close()

		if err == nil {
			return &api.CreateResponse{}, nil
		}
	}
	return nil, os.ErrExist
}

// Remove removes file or reports an error
func (server FileServer) Remove(ctx context.Context,
	request *api.RemoveRequest) (*api.RemoveResponse, error) {

	filePath := path.Join(server.storagePath, request.Inode)
	err := os.Remove(filePath)
	if err == nil {
		return &api.RemoveResponse{}, nil
	}
	return nil, err
}

// ReportDiskSpace generates a brief report about used space on disk
func (server FileServer) ReportDiskSpace(ctx context.Context,
	_ *api.Empty) (*api.DiskSpaceResponse, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(server.storagePath, &stat)
	if err == nil {
		return &api.DiskSpaceResponse{
				FreeBlocks:     int64(stat.Bavail),
				BusyBlocks:     int64(stat.Blocks - stat.Bavail),
				BlockSizeBytes: stat.Bsize},
			nil
	}
	return nil, err
}
