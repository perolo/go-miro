package miro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// SimpleCard object represents Miro SimpleCard.
//
// API doc: https://developers.miro.com/reference#card
//go:generate gomodifytags -file $GOFILE -struct Card -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Card -add-tags json -w -transform camelcase
type Card struct {
	Type        string  `json:"type"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Scale       float64 `json:"scale"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Assignee    struct {
		UserID string `json:"userId"`
	} `json:"assignee"`
	Style struct {
		BackgroundColor string `json:"backgroundColor"`
	} `json:"style"`
}

type SimpleCard struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type SimpleCardAssignee struct {
	Assignee struct {
		UserID string `json:"userId"`
	} `json:"assignee"`
}

type WidgetResponseDataType struct {
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
	Type string                   `json:"type"`
	Data []WidgetResponseDataType `json:"data"`
	Size int                      `json:"size"`
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
func (s *WidgetsService) CreateSimpleCard(ctx context.Context, boardid string, b *SimpleCard) (*CreateCardRespType, error) {
	req, err := s.client.NewPostRequest(fmt.Sprintf("boards/%s/%s", boardid, widgetsPath), b)
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

//https://api.miro.com/v1/boards/id/widgets/widgetId
func (s *WidgetsService) UpdateAssigneeCard(ctx context.Context, boardid string, widgetid string, b *SimpleCardAssignee) (*CreateCardRespType, error) {
	req, err := s.client.NewPatchRequest(fmt.Sprintf("boards/%s/%s/%s", boardid, widgetsPath, widgetid), b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

type WidgetMetadataResponseType struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Title     string `json:"title"`
	Issue     string `json:"issue"`
	YourAppID string `json:"appissue"`
	MetaData  string `json:"metadata"`
}

//type WidgetMetadataType map[int]string
type WidgetMetadataType struct {
	Title string
	AppId string
	Issue string
}

func (s *WidgetsService) GetWidgetMetadata(ctx context.Context, boardid string, widgetid string) (*WidgetMetadataResponseType, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("boards/%s/%s/%s", boardid, widgetsPath, widgetid))
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

	wresp := &WidgetMetadataResponseType{}
	/*
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("resp: %v\n", string(body))

	*/
	if err := json.NewDecoder(resp.Body).Decode(wresp); err != nil {
		return nil, err
	}

	return wresp, nil
}

//GET https://api.miro.com/v1/boards/{boardKey}/widgets/{widgetId}
func (s *WidgetsService) UpdateWidgetMetadata(ctx context.Context, boardid string, widgetid string, b *WidgetMetadataType) (*WidgetMetadataResponseType, error) {
	req, err := s.client.NewPatchRequest(fmt.Sprintf("boards/%s/%s/%s", boardid, widgetsPath, widgetid), b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
	board := &WidgetMetadataResponseType{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
}

/*
{"title": "example",
"metadata": {
"3074457352146652951" :{
"issue" : "STP-346" }
}
}
*/

func (this WidgetMetadataType) MarshalJSON() ([]byte, error) {
	str := strings.Replace(this.Title, `"`, `\"`, -1)
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(fmt.Sprintf("\"title\": \"%s\",", str))
	buffer.WriteString(fmt.Sprintf("\"metadata\": {\"%s\" :  { \"issue\" : \"%s\" }", this.AppId, this.Issue))

	buffer.WriteString("}}")
	return buffer.Bytes(), nil
}

func (p *WidgetMetadataResponseType) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "type" {
			p.Type = v.(string)
		}
		if strings.ToLower(k) == "id" {
			p.ID = v.(string)
		}
		if strings.ToLower(k) == "title" {
			p.Title = v.(string)
		}
		if strings.ToLower(k) == "metadata" {
			var _v map[string]interface{}
			_v = v.(map[string]interface{})
			for kk, vv := range _v {
				var _vv map[string]interface{}
				_vv = vv.(map[string]interface{})
				p.YourAppID = kk
				for kkk, vvv := range _vv {
					//fmt.Printf("kk: %s\n", kk)
					if strings.ToLower(kkk) == "issue" {
						p.Issue = vvv.(string)
						//p.Metadata.issue = "val"
					}
				}
			}

		}

	}

	return nil
}
