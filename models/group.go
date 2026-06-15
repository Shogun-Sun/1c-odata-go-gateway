package models

type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type ODataGroup struct {
	RefKey      string `json:"Ref_Key"`
	Description string `json:"Description"`
	Quantity    int    `json:"Численность,string"`
}

type GroupCreatePayload struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type ODataGroupCreate struct {
	Description string `json:"Description"`
	Quantity    int    `json:"Численность"`
}

type ODataGroupResponse struct {
	Value []ODataGroup `json:"value"`
}
