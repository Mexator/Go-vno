# GoVNO FS (GO written Virtual Network Operating FileSystem)
## Authors: Anton Brisilin @Mexator and Ivan Rybin @i1i1 

This is a distributed file system project we done during the Distributed 
Systems course.

## Architecture


There are several actors in our system:
- Clients
- File Servers
- Name Server
- DNS - used by fileservers to resolve addresses of nameserver at startup. 
Any ready solution can be used.

As a mean of communication between our custom services we decided to use 
[gRPC](https://grpc.io/) protocol. It is simple, convenient to the developers
and flexible. The following diagram illustrates all actor and possible RPC 
calls that can be done during the work of our system.

<center>
    <img src='./imgs/diagram.svg' width=500>
</center>

At the following sections each component will be discussed in greater detail.
### **Client**

To implement the client we decided to use [FUSE](https://en.wikipedia.org/wiki/Filesystem_in_Userspace) 
driver. This driver allows developers create custom filesystems running in 
userspace.  

Examples: sshfs, GlusterFS and our project 😄.

Diagram of how it works:  
<center>
<img src='https://upload.wikimedia.org/wikipedia/commons/thumb/0/08/FUSE_structure.svg/1280px-FUSE_structure.svg.png' width=300>
</center>

As you can see, each time a user wants to access the directory where our 
custom filesystem is mounted, it does a syscall to the kernel, which is 
redirected by FUSE to the local server. The local server is responsible for 
answering for syscalls such as `open()`, `read()` `stat()` and others. Our 
implementation does this with RPC calls to the remote Name- and File 
Servers.

Such an approach has one major advantage: the work of our DFS is **almost 
transparent** from a user's point of view. All they see is just a regular 
folder with regular files. They don't need to explicitly download or upload 
files, open some 
special terminal sessions or using a GUI. Everything that needs to be done at 
client's side is to launch our daemon server.

---
However, the requirements says that we should Dockerize our client. This 
disables almost all advantages we got from using `fuse` and gives us 
additional complexity, for example when launching the client container (it should have certain privileges to be launched correctly)

---

### **NameServer**
Mostly, the clients talk to the Name server which is used for indexing and 
replicating files. It stores whole directory tree and mapping of file 
fragments to their location on file servers (FS shortly).

There are several RPC calls that name server support:
- **Create(path, isDir)**  
    Creates a file or directory
- **Remove(path)**  
    Removes a file or directory
- **ReadDir(path)**  
    Returns info about all files stored in a given directory
- **Rename(fromPath, toPath)**  
    Renames file or directory. Used for `move` and `copy` operations
- **LookupStorageNodes(path)**  
    Returns addresses of storage nodes where given file fragment can be 
    found together with info required for reading/writing the fragment
- **ConnectFileServer(port)**  
    Connects a file server to the system. After a check whether sender 
    supports file server calls, it immediately can be used to store files 
    onto it.

The files fragments are replicated on several file servers.  
At startup Name Server waits connections from fileservers and clients. 
However, until there are at least two file servers attached, it will be 
impossible to create files, because there are not enough nodes to store or 
replicate files on.

### **File servers**
File servers are basically a key value storage with file fragments. 
Its RPC calls are:
- **Size(fragment)**  
    Return size of the fragment
- **Read(fragment, len, offset)**  
    Return len bytes of the fragment, starting from offset
- **Write(fragment, len, offset)**  
    Write len bytes to the fragment, starting from offset
- **Create**
    Create new fragment
- **Remove**
    Remove fragment

For now, we haven't done fragmentation, but it can be added in future. The 
'fragments' we refer to at this README are actually the whole files.

At startup File Servers try to establish connection with Name Server.
After they are connected, they send a lifecheck message to Name Server every 
second. If Name Server haven't receive lifecheck from FileServer, NS thinks 
that FS is down and replicates all files which were there to other File 
Servers.

## Technologies
- [FUSE](https://en.wikipedia.org/wiki/Filesystem_in_Userspace) - described at
client section
- Golang - all our servers and client are written on it
- [protobuf](https://wikipedia.org/wiki/Protocol_Buffers) and [gRPC](https://grpc.io/) - 
to define communication protocols between all actors.
- [make](https://en.wikipedia.org/wiki/Make_(software)) - to automate 
building 
and test deployment of the software. 
- [Docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) - same as make 

## Project layout

- Protobuf files are in [api/](./api) folder.
- [cmd/](./cmd) folder is used for all cli tools in project.
- [pkg/](./pkg) is for go libraries we use in [cmd/](./cmd)
- [pkg/api/](./pkgapi/) is for protobuf generated files
- [docker/](./docker) is for dockerfiles and docker-compose

## Local run

In order to start everything run

    $ make -j3 local_run
	
you will end up in container where you should run

    # ./client nameserver:3000 ./mnt/
	
after that you may go to `/mnt/` directory and use filesystem as every other
filesystem you used to know.

## Deployed system
We have deployed our system to AWS, with 1 name server and 3 file servers, 
each one on different instance, inside Docker container. You can test the 
filesystem by calling the following commands (**NOTE**: you should have 
`fuse` kernel module installed on your machine):
```shell
$ docker run -d --device /dev/fuse --cap-add=SYS_ADMIN --security-opt=apparmor:unconfined mexator/client
$ dc exec -it CONTAINER_NAME sh
```
This will get you to the client container, where you should run 
```shell
$ nohup ./client  ec2-18-191-130-90.us-east-2.compute.amazonaws.com:3000 /mnt &
```
After this, the DFS will be mounted to `/mnt` folder, inside the container. 
You can use it just like any other Unix directory.

## Provable contribution of all team members
There are two of us, and as you can see, we had a lot of pull request, several
issues and branches. I think, it can be used as a proof that each of worked 
hardly at this project.  
@Mexator is the one who created file server, worked on deploying, testing and 
Dockerizing of all parts. Nameserver was created together by @i1i1 and 
@Mexator (mostly @i1i1)  
@i1i1 created fuse client and designed architecture for the DFS.

## Docker images
- **Client** https://hub.docker.com/repository/docker/mexator/client
- **Fileserver** https://hub.docker.com/repository/docker/mexator/fileserver
- **Nameserver** https://hub.docker.com/repository/docker/mexator/nameserver