package apis

import (
	"blog/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// AddPost - Recieve uid, title, content, privatePost from PostForm to create a new post
func AddPost(c *gin.Context) {
	privatePost := false

	uid, title, content, privatePostStr := c.PostForm("uid"), c.PostForm("title"), c.PostForm("content"), c.PostForm("privatePost")

	ouid, err := primitive.ObjectIDFromHex(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
			"msg": "Cannot find the OID for Author",
			"uid": uid,
		})
		return
	}

	if privatePostStr == "true" || privatePostStr == "1" {
		privatePost = true
	}

	newPost := models.Post{
		Author:      ouid,
		Title:       title,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		PrivatePost: privatePost,
	}

	insertedID, err := models.AddPost(newPost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     err,
			"msg":     "Cannot add this Post",
			"newPost": newPost,
		})
	}

	newPost.ID = insertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"insertedID": insertedID,
		"newPost":    newPost,
	})
}

// UpdatePostByID - receive pid and post(updateFields) to update the post
func UpdatePostByID(c *gin.Context) {
	pid, postJSON := c.Param("pid"), c.PostForm("post")

	var post map[string]interface{}

	err := json.Unmarshal([]byte(postJSON), &post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":      err,
			"msg":      "Cannot unmarshal the given post",
			"postJSON": postJSON,
		})
		return
	}

	post["updatedAt"] = time.Now()
	upsertedID, err := models.UpdatePostByID(pid, post)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err":  err,
			"msg":  "Cannot Update this post",
			"post": post,
			"pid":  pid,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post":       post,
		"upsertedID": upsertedID,
	})

}

// DeletePostByID - Delete the post by given ID
func DeletePostByID(c *gin.Context) {
	pid := c.Param("pid")

	err := models.DeletePostByID(pid)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err,
			"msg": "Cannot delete this post",
			"pid": pid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pid": pid,
	})
}

// FindPostByID - find the post through pid
func FindPostByID(c *gin.Context) {
	pid := c.Param("pid")
	opid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
			"msg": "Cannot get the ObejctID of pid",
			"pid": pid,
		})
	}
	filterDetail := bson.M{"_id": opid}
	user, ok := c.Get("user")
	if (!ok || user == nil) || (user != nil && user.(*models.User).Role != "admin") {
		filterDetail["privatePost"] = false
	}

	post, err := models.FindOnePost(filterDetail)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err,
			"msg": "Cannot found this post",
			"pid": pid,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pid":  pid,
		"post": post,
	})
}

// FindAllPosts - Return all the posts
func FindAllPosts(c *gin.Context) {

	filterDetail := bson.M{}

	user, ok := c.Get("user")

	if (!ok || user == nil) || (user != nil && user.(*models.User).Role != "admin") {
		filterDetail["privatePost"] = false
	}
	// fmt.Println(ok)
	// fmt.Println(user)
	// fmt.Println(user == nil)
	// fmt.Println((!ok || user == nil))
	// fmt.Println((user != nil && user.(*models.User).Role != "admin"))
	// filterDetail["privatePost"] = false

	// else if user != nil && user.(*models.User).Role != "admin" {
	// 	filterDetail["privatePost"] = false
	// }

	skip, limit, sort := c.Query("skip"), c.Query("limit"), c.Query("sort")

	findOptions := options.Find()
	if strings.TrimSpace(skip) != "" {
		inputSkip, err := strconv.ParseInt(skip, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":  err,
				"msg":  "Cannot setup skip",
				"skip": skip,
			})
			return
		}
		findOptions.SetSkip(inputSkip)
	}

	if strings.TrimSpace(limit) != "" {
		inputLimit, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":   err,
				"msg":   "Cannot setup limit",
				"limit": limit,
			})
			return
		}
		findOptions.SetLimit(inputLimit)
	}

	sortMap := map[string]int{}
	if strings.TrimSpace(sort) != "" {
		if s := strings.Split(sort, "_"); len(s) == 2 {
			sortOrd, err := strconv.Atoi(s[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err":  err,
					"msg":  "Cannot get the sort order",
					"s[1]": s[1],
					"s[0]": s[0],
				})
				return
			}
			// fmt.Printf("s is %+v\n", s)
			// fmt.Printf("s[0] is %+v\n", s[0])
			// fmt.Printf("s[1] is %+v\n", s[1])
			// fmt.Printf("sortOrd is %+v\n", sortOrd)

			sortMap[s[0]] = sortOrd
		} else {
			sortMap[sort] = -1
		}

		findOptions.SetSort(sortMap)
	}

	posts, err := models.FindPosts(filterDetail, findOptions)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err,
			"msg": "Cannot retrieve all the posts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})

}
