package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

type QuestionRequest struct {
	Question string `json:"question"`
}

type AIResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func main() {
	r := gin.Default()
	r.POST("/ask", func(c *gin.Context) {
		var req QuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		answer, err := getAIAnswer(req.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI answer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"question": req.Question,
			"answer":   answer,
		})
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func getAIAnswer(question string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"question": question,
	})
	resp, err := http.Post("http://localhost:8000/get-answer", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", err
	}
	return aiResp.Answer, nil
}
