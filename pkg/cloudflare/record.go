package cloudflare

const (
	A     Type = "A"
	AAAA  Type = "AAAA"
	MX    Type = "MX"
	CNAME Type = "CNAME"
	TXT   Type = "TXT"
)

type (
	// Type of the DNS record
	Type string
	// An Error in the response
	Error struct {
		Code       int     `json:"code"`
		Message    string  `json:"message"`
		ErrorChain []Error `json:"error_chain"`
	}
	// Response
	Response struct {
		ResultInfo *struct {
			Page       int `json:"page"`
			TotalPages int `json:"total_pages"`
		} `json:"result_info"`
		Success  bool          `json:"success"`
		Errors   []Error       `json:"errors"`
		Messages []interface{} `json:"messages"`
	}

	// DNSResponse is the response returned from the `GET zones/:zone_identifier/dns_records` endpoint
	DNSResponse struct {
		Response
		Result *[]Record `json:"result"`
	}
	// DNSUpdateRequest updates a DNS entry calling the `PUT zones/:zone_identifier/dns_records/:identifier` endpoint
	DNSUpdateRequest struct {
		Name    string `json:"name"`
		Type    Type   `json:"type"`
		Content string `json:"content"`
		Proxied bool   `json:"proxied"`
		TTL     int    `json:"ttl"`
	}
	// A Record represents a DNS record
	Record struct {
		ID      string `json:"id"`
		Type    Type   `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
	}
)
