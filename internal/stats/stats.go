package stats

type (
	Key struct {
		Timestamp  int64
		Country    string
		Os         string
		Browser    string
		CampaignId uint32
	}

	Value struct {
		Requests    int64
		Impressions int64
	}

	Rows map[Key]Value
)

func (a Value) Assign(b Value) Value {
	res := a
	res.Requests += b.Requests

	res.Impressions += b.Impressions

	return res
}
