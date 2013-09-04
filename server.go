package o2aserver

type Server struct {
	storage               Storage
	appconfig             AppConfig
	tokengenauthorization TokenGenAuthorization
	tokengenaccess        TokenGenAccess
}

func NewServer(storage Storage, appconfig AppConfig, tokengenauthorization TokenGenAuthorization, tokengenaccess TokenGenAccess) *Server {
	return &Server{
		storage:               storage,
		appconfig:             appconfig,
		tokengenauthorization: tokengenauthorization,
		tokengenaccess:        tokengenaccess,
	}
}

func (s *Server) NewAuthorization() *Authorization {
	return NewAuthorization(s.storage, s.tokengenauthorization, s.appconfig)
}
