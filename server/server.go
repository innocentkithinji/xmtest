package server

type Server interface {
	Serve(port string)
}
