# This project was imported from my Gitea student portfolio

## net-cat

### Objectives

This project consists on recreating the **NetCat in a Server-Client Architecture** that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

- NetCat, `nc` system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open TCP connections, send UDP packages, listen on arbitrary TCP and UDP ports and many more.

- To see more information about NetCat inspect the manual `man nc`.

Your project must work in a similar way that the original  NetCat works, in other words, you must create a group chat. The project must have the following features :

- TCP connection between server and multiple clients (relation of 1 to many).
- A name requirement to the client.
- Control connections quantity.
- Clients must be able to send messages to the chat.
- Do not broadcast EMPTY messages from a client.
- Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : `[2020-01-20 15:48:41][client.name]:[client.message]`
- If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
- If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
- If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
- All Clients must receive the messages sent by other Clients.
- If a Client leaves the chat, the rest of the Clients must not disconnect.
- If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: `[USAGE]: ./TCPChat $port`

### Instructions

- Your project must be written in **Go**
- Start TCP server, listen and accept connections
- Your project must have Go-routines
- Your project must have channels or Mutexes
- Maximum 10 connections
- The code must respect the [**good practices**](../good-practices/README.md)
- It is recommended to have **test files** for [unit testing](https://go.dev/doc/tutorial/add-a-test) both the server connection and the client.

- You have to be able to handle the errors from server side and client side

### Allowed Packages

- io
- log
- os
- fmt
- net
- sync
- time
- bufio
- errors
- strings
- reflect

### Usage

```console
$ go run .
Listening on the port :8989
$ go run . 2525
Listening on the port :2525
$ go run . 2525 localhost
[USAGE]: ./TCPChat $port
$
```
Or you can run multiple servers which means you can create multiple chatgroups by typing `go run . {$list of ports}`
```console
go run . 8001 8002 8003
Listening on the ports :8001 :8002 :8003
```
Once server is started open new terminal(s) and join the server by typing `nc localhost {$port}`
```console
nc localhost 8001
```
Type `/changeName` to change client name.