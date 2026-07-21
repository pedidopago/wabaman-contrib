package fbgraph

type LocationObject struct {
	// Endereço da localização. (Opcional)
	Address string `json:"address,omitzero"`
	// Latitude da localização em graus decimais. (Obrigatório)
	Latitude string `json:"latitude"`
	// Longitude da localização em graus decimais. (Obrigatório)
	Longitude string `json:"longitude"`
	// Nome da localização. (Opcional)
	Name string `json:"name,omitzero"`
}
