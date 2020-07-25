package models

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/GhvstCode/Blog-Api/api/utils"
	l "github.com/GhvstCode/Blog-Api/api/utils/logger"
)

type BlogModel struct {
	ID        primitive.ObjectID `bson:"_id, omitempty" json:"id, omitempty"`
	Title     string             `bson:"title" json:"title, omitempty"`     //Ensure title is not empty or greater than 150 characters
	Content   string             `bson:"content" json:"content, omitempty"` //Ensure content is not empty
	Author    string             `bson:"author" json:"author, omitempty"`
	OwnerId   primitive.ObjectID `bson:"ownerId, omitempty" json:"ownerId, omitempty"`
	Published *bool              `bson:"published" json:"published, omitempty"`
	Paid      *bool              `bson:"Paid" json:"Paid, omitempty"`
	Price     int                `bson:"price" json:"price, omitempty"`
}

type ReBlogModel struct {
	ID        string `json:"id, omitempty"`
	Title     string `json:"title, omitempty"`
	Content   string `json:"content, omitempty"`
	Author    string `json:"author, omitempty"`
	OwnerId   string `json:"ownerId, omitempty"`
	Published *bool  `json:"published, omitempty"`
	Paid      *bool  `json:"Paid, omitempty"`
	Price     int    `json:"price, omitempty"`
}

type UpdateBlogModel struct {
	Title     string `json:"title, omitempty"`
	Content   string `json:"content, omitempty"`
	Published *bool  `json:"published, omitempty"`
	Paid      *bool  `json:"Paid, omitempty"`
	Price     int    `json:"price, omitempty"`
}

func getID(id string) (primitive.ObjectID, error) {
	postId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return postId, err
	}

	return postId, nil
}
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Validate(b *BlogModel) *utils.Data {
	if len(b.Title) >= 150 || len(b.Title) == 0 {
		return utils.Response(false, "Title Must be greater than zero & less than 150 characters", http.StatusBadRequest)
	}

	if len(b.Content) == 0 {
		return utils.Response(false, "Content must not be empty", http.StatusBadRequest)
	}

	if len(b.Author) == 0 {
		return utils.Response(false, "Please fill in the author field", http.StatusBadRequest)
	}

	return utils.Response(true, "Validated", http.StatusOK)
}

func (b *BlogModel) Create(Owner string) *utils.Data {
	//Check if a user already has a post with the title they are trying to create/update to avoid possible duplicates
	resp := Validate(b)
	ok := resp.Result
	if !ok {
		return resp
	}

	OwnId, err := primitive.ObjectIDFromHex(Owner)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
	}

	res, err := Blog.InsertOne(context.TODO(), &BlogModel{
		ID:        primitive.NewObjectID(),
		Title:     b.Title,
		Content:   b.Content,
		Author:    b.Author,
		OwnerId:   OwnId,
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

func (b *UpdateBlogModel) UpdatePost(id string) *utils.Data {
	//To-Do Validate Update.
	postId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
	}

	filter := bson.D{{"_id", postId}}
	update := bson.M{"$set": b}
	//opts := options.FindOneAndUpdate().new(true)

	var bm ReBlogModel
	var c BlogModel
	err = Blog.FindOneAndUpdate(context.TODO(), filter, update).Decode(&bm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			l.ErrorLogger.Print(err)
			return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
		}
		l.ErrorLogger.Print(err)
		return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
	}

	err = Blog.FindOne(context.TODO(), filter).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			l.ErrorLogger.Print(err)
			return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
		}
		l.ErrorLogger.Print(err)
		return utils.Response(false, "An Error occurred, Unable to Create Post", http.StatusInternalServerError)
	}

	res := &ReBlogModel{
		ID:        c.ID.Hex(),
		Title:     c.Title,
		Content:   c.Content,
		Author:    c.Author,
		OwnerId:   c.OwnerId.Hex(),
		Published: c.Published,
		Paid:      c.Paid,
		Price:     c.Price,
	}
	response := utils.Response(true, "Updated", http.StatusCreated)
	response.Data = [1]*ReBlogModel{res}
	return response
}

func DeletePost(id string) *utils.Data {
	postId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Delete Post", http.StatusInternalServerError)
	}

	filter := bson.D{{"_id", postId}}
	_, err = Blog.DeleteOne(context.TODO(), filter)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Delete Post", http.StatusInternalServerError)
	}

	response := utils.Response(true, "Deleted", http.StatusOK)
	return response
}

func GetPost(id string, UserID string) *utils.Data {
	var c *UserModel
	var b *BlogModel
	postId, err := getID(id)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Fetch Post", http.StatusInternalServerError)
	}

	userId, err := getID(UserID)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to Fetch Post", http.StatusInternalServerError)
	}

	//filter := bson.D{{"_id", postId}}
	err = Blog.FindOne(context.TODO(), bson.D{{"_id", postId}}).Decode(&b)
	if err != nil {
		l.ErrorLogger.Print(err)
		return utils.Response(false, "An Error occurred, Unable to Fetch Post", http.StatusInternalServerError)
	}

	err = User.FindOne(context.TODO(), bson.D{{"_id", userId}}).Decode(&c)
	if err != nil {
		l.ErrorLogger.Print(err)
		return utils.Response(false, "An Error occurred, Unable to Fetch Post", http.StatusInternalServerError)
	}
	//fmt.Print(b.OwnerId)
	//fmt.Print(userId)
	_, returnBool := Find(c.Subscriptions, id)
	//if  b.OwnerId == userId{
	//	return utils.Response(false, "Post Not found" , http.StatusNotFound)
	//}
	//if !*b.Published && b.OwnerId == userId || !returnBool{
	//	return utils.Response(false, "Post Not found" , http.StatusNotFound)
	//}
	if b.Published == nil || !(*b.Published) && !returnBool || !(b.OwnerId == userId) {
		return utils.Response(false, "Post Not found", http.StatusNotFound)
	}
	res := &ReBlogModel{
		ID:        b.ID.Hex(),
		Title:     b.Title,
		Content:   b.Content,
		Author:    b.Author,
		OwnerId:   b.OwnerId.Hex(),
		Published: b.Published,
		Paid:      b.Paid,
		Price:     b.Price,
	}

	response := utils.Response(true, "Updated", http.StatusCreated)
	response.Data = [1]*ReBlogModel{res}
	return response
}
