package entity

type Company struct {
	ID             string `json:"id" bson:"_id"`
	Name           string `json:"name" validate:"required,max=15"`
	Description    string `json:"description" validate:"max=3000"`
	Employees      int    `json:"employees" validate:"required"`
	Registered     bool   `json:"registered"`
	Type           string `json:"type" validate:"required,oneof=Corporations NonProfit Cooperative 'Sole Proprietorship'"`
	NameIdentifier string
	OwnerId        string `json:"owner_id"`
}
