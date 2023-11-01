package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CardAutocompleteResponse struct {
	Object      string   `json:"object"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

type CardResponse struct {
	Id            uuid.UUID         `json:"id"`
	OracleId      uuid.UUID         `json:"oracle_id"`
	Name          string            `json:"name"`
	ImageUris     CardImageUris     `json:"image_uris"`
	ManaCost      string            `json:"mana_cost"`
	Cmc           float32           `json:"cmc"`
	TypeLine      string            `json:"type_line"`
	OracleText    string            `json:"oracle_text"`
	Colors        []string          `json:"colors"`
	ColorIdentity []string          `json:"color_identity"`
	Legalities    map[string]string `json:"legalities"`
}

type CardImageUris struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	Png        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

func main() {
	s := gin.Default()
	s.HTMLRender = &TemplRenderer{}

	s.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", Home())
	})

	s.POST("/card/search/", searchCardsAutocomplete)
	s.GET("/card/:cardname", getCardByName)

	s.Run(":8080")
}

func getCardByName(c *gin.Context) {
	query := c.Param("cardname")
	resp, err := http.Get(
		fmt.Sprintf("https://api.scryfall.com/cards/named?exact=%s", query),
	)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	var cardResp CardResponse
	if err := json.NewDecoder(resp.Body).Decode(&cardResp); err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "", CardComponent(cardResp))
}

func searchCardsAutocomplete(c *gin.Context) {
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

	var searchRes CardAutocompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "", SearchResult(searchRes.Data))
}
