package svc

// routes registers the handlers to the specified routes
func (s *QServer) routes() {
	s.router.Handle("/", s.Home()).Methods("GET")
}
