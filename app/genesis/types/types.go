package types

	
type AhaRow struct {
	Name string `json:"name"`
	Email string `json:"email"`
}
	
type AhaQuery struct {
	Anything string `json:"anything"`
	Example string `json:"example"`
}

	
type SuperTestRes struct {
	Cool bool `json:"cool"`
}
	
type SuperTestQuery struct {
	Example string `json:"example"`
}

	
type Aha3Row struct {
}
	
type Aha3Body struct {
	Size string `json:"size"`
}
		