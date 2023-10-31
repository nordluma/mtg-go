package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CardSearchResponse struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

func (s CardSearchResponse) TextOutput() string {
	if len(s.Data) > 0 {
		return fmt.Sprintf(
			"Object: %s\nTotalValue: %v\nFirst Result: %v",
			s.Object,
			s.TotalValues,
			s.Data[0],
		)
	}

	return fmt.Sprintf("Object: %s\nTotal values: %v", s.Object, s.TotalValues)
}

func main() {
	s := gin.Default()
	s.HTMLRender = &TemplRenderer{}

	s.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", Home())
	})

	s.POST("/card/search/", func(c *gin.Context) {
		var searchRes CardSearchResponse
		query := c.PostForm("card-name-search")
		resp, err := http.Get(
			fmt.Sprintf(
				"https://api.scryfall.com/cards/autocomplete?q=%s&inlude_extras=false",
				query,
			),
		)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
			log.Println(err)
		}

		c.HTML(http.StatusOK, "", SearchResult(searchRes.Data))
	})

	s.Run(":8080")
}
