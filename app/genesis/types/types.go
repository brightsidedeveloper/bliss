package types

	
type AhaRow struct {
	Name string `json:"name"`
	Email string `json:"email"`
}
	
type AhaQuery struct {
	Example string `json:"example"`
	Anything string `json:"anything"`
}

	
type SuperTestRes struct {
	Cool bool `json:"cool"`
}
	
type SuperTestQuery struct {
	Example string `json:"example"`
}

	
type Aha3Row struct {
	Count int `json:"count"`
}
	
type Aha3Body struct {
	Size string `json:"size"`
}
		