package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"github.com/therealofarah/go-calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var validate = validator.New()
var entryCollection *mongo.Collection = OpenCollection(Client,"calories")
func AddEntry(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	var entry models.Entry
	if err:=c.BindJSON(&entry); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	validationErr:= validate.Struct(entry)
	if validationErr!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	entry.ID = primitive.NewObjectID()
	result, insertErr:= entryCollection.InsertOne(ctx, entry)
	if insertErr!=nil{
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}
func GetEntryById(c *gin.Context){
	EntryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryId)

	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	defer cancel()
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id":docID}).Decode(&entry); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entry)
	c.JSON(http.StatusOK,entry)
}
func GetEntries(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	var entries []bson.M
	cursor,err := entryCollection.Find(ctx, bson.M{})
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &entries); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
	}
	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}
func GetEntriesByIngredient(c *gin.Context){
	ingredient :=c.Params.ByName("id")

	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	var entries []bson.M
	cursor,err :=entryCollection.Find(ctx, bson.M{"ingredinets":ingredient})
		if err !=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			fmt.Println(err)
			return
		}
		if err = cursor.All(ctx,&entries);err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			fmt.Println(err)
			return
		}
	c.JSON(http.StatusOK, entries)
	defer cancel()
}
func UpdateEntry(c *gin.Context){
	entryID :=c.Params.ByName("id")
	docID,_:=primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	var entry models.Entry
	if err := c.BindJSON(&entry); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	validationErr:= validate.Struct(entry)
	if validationErr!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	results, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id":docID},
		bson.M{
			"dish":entry.Dish,
			"fat":entry.Fat,
			"protein":entry.Protein,
			"carb":entry.Carb,
			"ingredinets":entry.Ingredients,
			"calories":entry.Calories,
		},
	)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"Error":err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, results.ModifiedCount)
	defer cancel()
}
func UpdateIngredient(c *gin.Context){
	entryID :=c.Params.ByName("id")
	docID,_:=primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	type Ingredent struct{
		Ingredent *string `json:ingredients`
	}
	var ingredent Ingredent
	if err := c.BindJSON(&ingredent); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	result, err:=entryCollection.UpdateOne(ctx, bson.M{"_id":docID},
		bson.D{
			{"$set",bson.D{{"ingredients", ingredent.Ingredent}}},
		},
	)
	if err !=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}
func DeleteEntry(c *gin.Context){
	entryID := c.Params.ByName("id")
	docId, _ :=primitive.ObjectIDFromHex(entryID)
	
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	result, err:= entryCollection.DeleteOne(ctx,bson.M{"_id":docId})

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)

}