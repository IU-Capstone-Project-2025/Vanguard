package shared

type Session struct {
	ID               string `redis:"id"`
	Code             string `redis:"code"`
	State            string `redis:"state"`
	ServerWsEndpoint string `redis:"serverWsEndpoint"`
}
