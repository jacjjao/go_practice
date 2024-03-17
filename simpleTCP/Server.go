package main

type Server struct {
	joinQueue chan *Client
	clients   []*Client
}

func (server *Server) Run() {
	for {
		select {
		case client := <-server.joinQueue:
			server.clients = append(server.clients, client)
			client.Run()
		default:
			// do nothing
		}
		for _, client := range server.clients {
			select {
			case msg := <-client.msgQueue:
				for _, wok := range server.clients {
					if wok.id == msg.sender.id {
						continue
					}
					wok.incomeQueue <- &msg.body
				}
			default:
				// do nothing
			}
		}
		remainClients := make([]*Client, 0, len(server.clients))
		for _, client := range server.clients {
			select {
			case <-client.leftSignal:
				// do nothing
			default:
				remainClients = append(remainClients, client)
			}
		}
		server.clients = remainClients
	}
}
