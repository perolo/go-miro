package miro

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	widgetsPath = "widgets"
)

// WidgetsService handles communication to Miro Widgets API.
//
// API doc: https://developers.miro.com/reference#board-object
type WidgetsService service

type Widget interface {
}

// Sticker object represents Miro Sticker.
//
// API doc: https://developers.miro.com/reference#sticker
//go:generate gomodifytags -file $GOFILE -struct Sticker -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Sticker -add-tags json -w -transform camelcase
type Sticker struct {
}

// Shape object represents Miro Shape.
//
// API doc: https://developers.miro.com/reference#shape
//go:generate gomodifytags -file $GOFILE -struct Shape -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Shape -add-tags json -w -transform camelcase
type Shape struct {
}

// Text object represents Miro Text.
//
// API doc: https://developers.miro.com/reference#text
//go:generate gomodifytags -file $GOFILE -struct Text -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Text -add-tags json -w -transform camelcase
type Text struct {
}

// Line object represents Miro Line.
//
// API doc: https://developers.miro.com/reference#line
//go:generate gomodifytags -file $GOFILE -struct Line -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Line -add-tags json -w -transform camelcase
type Line struct {
}

// Card object represents Miro Card.
//
// API doc: https://developers.miro.com/reference#card
//go:generate gomodifytags -file $GOFILE -struct Card -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Card -add-tags json -w -transform camelcase
type Card struct {
	Type        string  `json:"type"`
//	X           float64 `json:"x"`
//	Y           float64 `json:"y"`
//	Scale       float64 `json:"scale"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
/*	Date        string  `json:"date"`
	Assignee    struct {
		UserID string `json:"userId"`
	} `json:"assignee"`
	Style struct {
		BackgroundColor string `json:"backgroundColor"`
	} `json:"style"`*/
}
type WidgetResponseDataType struct  {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Date        string      `json:"date"`
	Card        interface{} `json:"card"`
	X           float64     `json:"x"`
	Rotation    float64     `json:"rotation"`
	Assignee    struct {
		UserID string `json:"userId"`
	} `json:"assignee"`
	Y     float64 `json:"y"`
	Scale float64 `json:"scale"`
	Style struct {
		BackgroundColor string `json:"backgroundColor"`
	} `json:"style"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy struct {
		Type string `json:"type"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"createdBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy struct {
		Type string `json:"type"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"modifiedBy"`
}


type WidgetResponseType struct {
	Type string `json:"type"`
	Data []WidgetResponseDataType  `json:"data"`
	Size int `json:"size"`
}

//https://api.miro.com/v1/boards/id/widgets/
func (s *WidgetsService) ListAllWidgets(ctx context.Context, id string) (*WidgetResponseType, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("boards/%s/%s", id, widgetsPath))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	wresp := &WidgetResponseType{}
	/*
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp: %v\n", string(body))

	 */
	if err := json.NewDecoder(resp.Body).Decode(wresp); err != nil {
		return nil, err
	}

	return wresp, nil
}

type CreateCardRespType struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Date        string      `json:"date"`
	Card        interface{} `json:"card"`
	X           float64     `json:"x"`
	Rotation    float64     `json:"rotation"`
	Assignee    struct {
		UserID string `json:"userId"`
	} `json:"assignee"`
	Y     float64 `json:"y"`
	Scale float64 `json:"scale"`
	Style struct {
		BackgroundColor string `json:"backgroundColor"`
	} `json:"style"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy struct {
		Type string `json:"type"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"createdBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy struct {
		Type string `json:"type"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"modifiedBy"`
	Capabilities struct {
		Editable bool `json:"editable"`
	} `json:"capabilities"`
}

//https://api.miro.com/v1/boards/id/widgets
func (s *WidgetsService) CreateCard(ctx context.Context, id string,b *Card) (*CreateCardRespType, error) {
	req, err := s.client.NewPostRequest(fmt.Sprintf("boards/%s/%s", id, widgetsPath), b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}
	/*
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp: %v\n", string(body))

	 */
	board := &CreateCardRespType{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
}
