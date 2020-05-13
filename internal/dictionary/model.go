package dictionary

type Citizenship struct {
	tableName struct{} `pg:"maasapi.citizenship"`
	ID        int      `pg:"id,pk" json:"id"`
	Name      string   `pg:"smallname" json:"name"`
	Fullname  string   `pg:"fullname" json:"fullName"`
	Alpha2    string   `pg:"alpha2" json:"alpha2"`
	Alpha3    string   `pg:"alpha3" json:"alpha3"`
	Lang      string   `pg:"lng" json:"lng"`
}
