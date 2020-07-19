package models

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/GhvstCode/Blog-Api/utils"
	l "github.com/GhvstCode/Blog-Api/utils/logger"
)

//We want to check before we show any article
//if the article is published
//if the article paid, if it is and the person reading it is not the Owner, or has this article ID in its Subscription array then disallow from reading the article.
type BlogModel struct {
	 ID        primitive.ObjectID `bson:"_id, omitempty" json:"id, omitempty"`
	 Title     string 			  `bson:"title" json:"title, omitempty"`//Ensure title is not empty or greater than 150 characters
	 Content   string 		      `bson:"content" json:"content, omitempty"`//Ensure content is not empty
	 Author    string			   `bson:"author" json:"author, omitempty"`
	 OwnerId     primitive.ObjectID `bson:"ownerId, omitempty" json:"ownerId, omitempty"`
	 Published bool					`bson:"published" json:"published, omitempty"`
	 Paid      bool					`bson:"Paid" json:"Paid, omitempty"`
	 Price     float64				`bson:"price" json:"price, omitempty"`
}

type ReBlogModel struct {
	ID        string   			   `json:"id, omitempty"`
	Title     string 			   `json:"title, omitempty"`
	Content   string 		       `json:"content, omitempty"`
	Author    string			   `json:"name, omitempty"`
	OwnerId   string 				`json:"ownerId, omitempty"`
	Published bool					`json:"published, omitempty"`
	Paid      bool					`json:"Paid, omitempty"`
	Price     float64				`json:"price, omitempty"`
}

func Validate(b *BlogModel) *utils.Data{
	if len(b.Title) >= 150 || len(b.Title) == 0 {
		return utils.Response(false,"Title Must be greater than zero & less than 150 characters", http.StatusBadRequest)
	}

	if len(b.Content) == 0 {
		return utils.Response(false,"Content must not be empty", http.StatusBadRequest)
	}

	if len(b.Author) == 0 {
		return utils.Response(false,"Please fill in the author field", http.StatusBadRequest)
	}

	return utils.Response(true,"Validated", http.StatusOK)
}

func (b *BlogModel)Create(Owner string) *utils.Data {
	resp := Validate(b)
	ok := resp.Result; if !ok {
		return resp
	}

	Ownid, err := primitive.ObjectIDFromHex(Owner)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Create Post" , http.StatusInternalServerError)
	}

	res, err := Blog.InsertOne(context.TODO(), &BlogModel{
		ID:        primitive.NewObjectID(),
		Title:     b.Title,
		Content:   b.Content,
		Author:    b.Author,
		OwnerId:   Ownid,
		Published: b.Published,
		Paid:      b.Paid,
		Price:     b.Price,
	})

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to create Post", http.StatusInternalServerError)
	}

	var UID string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		UID = oid.Hex()
	}

	r := &ReBlogModel{
		ID:        UID,
		Title:     b.Title,
		Content:   b.Content,
		Author:    b.Author,
		OwnerId:   Owner,
		Published: b.Published,
		Paid:      b.Paid,
		Price:     b.Price,
	}

	response := utils.Response(true, "created", http.StatusCreated)
	response.Data = [1]*ReBlogModel{r}
	return response
}