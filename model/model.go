package model

// Gateway: base model class for the gateway object
type Gateway struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IpAddress string `json:"ip_address"`
}

// Route: base model class for the route object
type Route struct {
	ID        string `json:"id"`
	Prefix    string `json:"prefix"`
	GatewayId int64  `json:"gateway_id"`
}

type CustomRoute struct {
	Prefix    string `json:"prefix"`
	GatewayId int64  `json:"gateway_id"`
}

type CustomRouteResponse struct {
	ID      int64   `json:"id"`
	Prefix  string  `json:"prefix"`
	Gateway Gateway `json:"gateway"`
}

type ErrorGateway struct {
	Message string `json:"message"`
	Param   string `json:"param"`
}
