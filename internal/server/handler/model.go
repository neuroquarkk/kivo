package handler

type PutReq struct {
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

type ExpireReq struct {
	TTL string `json:"ttl"`
}
