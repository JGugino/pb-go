package services

type Pocketbase struct {
	Auth       *PBAuth       `json:"auth"`
	Collection *PBCollection `json:"collection"`
	Record     *PBRecord     `json:"record"`
}

func (pb *Pocketbase) Init(url string) error {
	pb.Auth = &PBAuth{
		BaseURL: url,
	}

	pb.Collection = &PBCollection{
		BaseURL: url,
	}

	pb.Record = &PBRecord{
		BaseURL: url,
	}
	return nil
}
