package model

import (
	"cloud.google.com/go/firestore"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Tutorial struct {
	ID           string      `firestore:"id" json:"id"`
	Name         string      `firestore:"name" json:"name"`
	Description  string      `firestore:"description" json:"description"`
	Domain       string      `firestore:"domain" json:"domain"`
	PathOperator string      `firestore:"pathOperator" json:"pathOperator"`
	PathValue    string      `firestore:"pathValue" json:"pathValue"`
	Parameters   []Parameter `firestore:"parameters" json:"parameters"`
	Settings     Settings    `firestore:"settings" json:"settings"`
	IsActive     bool        `firestore:"isActive" json:"isActive"`
	GaID         string      `firestore:"gaId" json:"gaId"`
	Steps        []*Step     `json:"steps"`
}

type Settings struct {
	Once bool `firestore:"once" json:"once"`
}

type Parameter struct {
	Key   string `firestore:"key" json:"key"`
	Value string `firestore:"value" json:"value"`
}

func NewTutorial(doc *firestore.DocumentSnapshot) (*Tutorial, error) {
	var tutorial *Tutorial
	if err := doc.DataTo(&tutorial); err != nil {
		return nil, err
	}
	tutorial.ID = doc.Ref.ID
	return tutorial, nil
}

type Step struct {
	ID      string `firestore:"id" json:"id"`
	Type    string `firestore:"type" json:"type"`
	Trigger struct {
		Target      string `firestore:"target" json:"target"`
		Event       string `firestore:"event" json:"event"`
		WaitingTime int    `firestore:"waitingTime" json:"waitingTime"`
	} `firestore:"trigger" json:"trigger"`
	HighlightTarget string                 `firestore:"highlightTarget" json:"highlightTarget"`
	Config          map[string]interface{} `firestore:"config" json:"config"`
	Order           int                    `firestore:"order" json:"order"`
	PathOperator    string                 `firestore:"pathOperator" json:"pathOperator"` // ALL, EQUALS, STARTS_WITH, ENDS_WITH, CONTAINS, REGEX, NOT_EQUALS
	PathValue       string                 `firestore:"pathValue" json:"pathValue"`
	Parameters      []Parameter            `firestore:"parameters" json:"parameters"`
}

func NewStep(doc *firestore.DocumentSnapshot) (*Step, error) {
	var step *Step
	if err := doc.DataTo(&step); err != nil {
		return nil, err
	}
	step.ID = doc.Ref.ID
	return step, nil
}

type Ga struct {
	ID           string `firestore:"id" json:"id"`
	Email        string `firestore:"email" json:"-"`
	RefreshToken string `firestore:"refreshToken" json:"-"`
	AccountID    string `firestore:"accountId" json:"accountId"`
	AccountName  string `firestore:"accountName" json:"accountName"`
	PropertyID   string `firestore:"propertyId" json:"propertyId"`
	PropertyName string `firestore:"propertyName" json:"propertyName"`
	ViewID       string `firestore:"viewId" json:"viewId"`
	ViewName     string `firestore:"viewName" json:"viewName"`
}

func NewGa(doc *firestore.DocumentSnapshot) (*Ga, error) {
	var ga *Ga
	if err := doc.DataTo(&ga); err != nil {
		return nil, err
	}
	ga.ID = doc.Ref.ID
	return ga, nil
}

type QueryResult struct {
	Tutorial *Tutorial `json:"tutorial,omitempty"`
	Ga       *Ga       `json:"ga,omitempty"`
}

const PathAll = "ALL"
const PathEquals = "EQUALS"
const PathStartsWith = "STARTS_WITH"
const PathRegex = "REGEX"

func ValidateUrlPath(pathOperator string, pathValue string, urlPath string) (bool, error) {
	pv := pathValue
	if pathValue[len(pathValue)-1:] == "/" {
		pv = pathValue[:len(pv)-1]
	}
	fmt.Println(pv)
	up := urlPath
	if urlPath[len(urlPath)-1:] == "/" {
		up = urlPath[:len(urlPath)-1]
	}
	switch pathOperator {
	case PathEquals:
		return pv == up, nil
	case PathStartsWith:
		return strings.HasPrefix(up, pv), nil
	case PathRegex:
		valid, err := regexp.MatchString(pv, up)
		if err != nil {
			return false, err
		}
		return valid, nil
	case PathAll:
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("path operator is unknown, %s", pathOperator))
}
