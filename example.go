

import (
	"github.com/gin-gonic/gin"
)

type note struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var notes = []note{
	{Title: "Homework Plan", Body: "Do Linear Algebra on Saturday and everything else on Friday"},
	{Title: "Sunday Reminder", Body: "Go to church during the day!"},
	{Title: "Make daily checklist", Body: ""},
}

func getNotes(context *gin.Context) {
	context.IndentedJSON(200, notes)
}

func main() {
	// Initialize server
	server := gin.Default()

	// Create an endpoint	localhost:8080/nich-neuron
	server.GET("/nich-neuron", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Welcome to Nich Neuron. We hope you enjoy your stay.",
		})
	})

	// Second endpoint	localhost:8080/nich-neuron/2
	server.GET("/nich-neuron/2", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Second NN route.",
		})
	})

	// Looks for variable in request	localhost:8080/nich-neuron/3?user=USERNAME
	server.GET("/nich-neuron/3", func(context *gin.Context) {

		username := context.DefaultQuery("user", "guest")
		username = "Hello, " + username
		context.JSON(200, gin.H{
			"message": username,
		})
	})

	server.GET("/nich-neuron/notes", getNotes)

	server.Run(":8080")

}
