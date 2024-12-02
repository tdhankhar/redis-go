package server

import (
	"log"
	"syscall"
)

type eventLoop struct {
	serverFd int
	kqueueFd int
}

func initEventLoop(serverFd int) *eventLoop {
	kqueueFd, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	serverEvent := syscall.Kevent_t{
		Ident: uint64(serverFd),
		Filter: syscall.EVFILT_READ,
		Flags: syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data: 0,
		Udata: nil,
	}
	serverEventRegistered, err := syscall.Kevent(kqueueFd, []syscall.Kevent_t{serverEvent}, nil, nil)
	if err != nil || serverEventRegistered == -1 {
		panic(err)
	}

	return &eventLoop{
		serverFd: serverFd,
		kqueueFd: kqueueFd,
	}
}

func (e *eventLoop) Run() {
	clients := 0
	events := make([]syscall.Kevent_t, maxClients)
	defer syscall.Close(e.kqueueFd)

	for {
		nevents, err := syscall.Kevent(e.kqueueFd, nil, events, nil)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			eventFd := int(events[i].Ident)
			if events[i].Flags & syscall.EV_EOF != 0 {
				syscall.Close(eventFd)
				clients -= 1
				log.Println("client disconnected")
				log.Println("concurrent clients", clients)
		    } else if eventFd == e.serverFd {
				socketFd, _, err := syscall.Accept(e.serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}

				clients += 1
				log.Println("client connected")
				log.Println("concurrent clients", clients)
				syscall.SetNonblock(e.serverFd, true)

				socketEvent := syscall.Kevent_t{
					Ident: uint64(socketFd),
					Filter: syscall.EVFILT_READ,
					Flags: syscall.EV_ADD,
					Fflags: 0,
					Data: 0,
					Udata: nil,
				}
				socketEventRegistered, err := syscall.Kevent(e.kqueueFd, []syscall.Kevent_t{socketEvent}, nil, nil)
				if err != nil || socketEventRegistered == -1 {
					log.Println("err", err)
					continue
				}
			} else {
				clientConnection := fdConn{ fd: eventFd }
				cmd, err := readCommand(clientConnection)
				if err != nil {
					log.Println("err", err)
					continue
				}
				respond(clientConnection, cmd)
			}
		}
	}
}
