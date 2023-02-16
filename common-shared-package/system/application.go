package system

type Controller struct {
}

type UserContext struct {
	UserId      string `json:"id" bson:"_id"`
	Name        string `json:"name" binding:"required,min=2,max=200" bson:"name"`
	EmailId     string `json:"emailId" binding:"required,email,min=5,max=200" bson:"emailId"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
	AccessToken string `json:"accessToken"`
}
