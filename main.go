package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Generalbelly/q-torial-api/model"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/firestore"
)

func query(ctx context.Context, userKey string, targetUrl string) (*model.QueryResult, error) {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatalf("PROJECT_ID is not set")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}
	defer client.Close()

	var selectedTutorial *model.Tutorial = nil
	var ga *model.Ga = nil

	u, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}

	docs, err := client.Collection("users").
		Doc(userKey).
		Collection("tutorials").
		Where("isActive", "==", true).
		OrderBy("pathPriority", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	matchedTutorials := make([]*model.Tutorial, 0)
	for _, doc := range docs {
		tutorial, err := model.NewTutorial(doc)
		if err != nil {
			return nil, err
		}
		valid, err := model.ValidateUrlPath(tutorial.PathOperator, tutorial.PathValue, u.Path)
		if err != nil {
			return nil, err
		}
		if valid {
			hasSameParameters := true
			for _, param := range tutorial.Parameters {
				if u.Query().Get(param.Key) != param.Value {
					hasSameParameters = false
				}
			}
			tutorialDomain, err := url.Parse(tutorial.Domain)
			if err != nil {
				return nil, err
			}
			if hasSameParameters && (len(tutorial.Domain) == 0 || (len(tutorial.Domain) > 0 && tutorialDomain.Hostname() == u.Hostname())) {
				matchedTutorials = append(matchedTutorials, tutorial)
			}
		}
	}
	if len(matchedTutorials) > 0 {
		selectedTutorial = matchedTutorials[0]
	}
	if selectedTutorial != nil {
		selectedTutorialRef := client.Collection("users").Doc(userKey).Collection("tutorials").Doc(selectedTutorial.ID)
		docs, err := selectedTutorialRef.Collection("steps").OrderBy("order", firestore.Asc).Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
		for _, doc := range docs {
			step, err := model.NewStep(doc)
			if err != nil {
				return nil, err
			}
			selectedTutorial.Steps = append(selectedTutorial.Steps, step)
		}
		if len(selectedTutorial.GaID) > 0 {
			doc, err := client.Collection("users").Doc(userKey).Collection("gas").Doc(selectedTutorial.GaID).Get(ctx)
			if err != nil {
				return nil, err
			}
			ga, err = model.NewGa(doc)
		}
	}
	return &model.QueryResult{
		Tutorial: selectedTutorial,
		Ga:       ga,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Content-Type", "application/json")
	var response *model.QueryResult
	var err error
	var res []byte

	q := r.URL.Query()
	url := q.Get("url")
	userKey := q.Get("key")

	if len(url) == 0 || len(userKey) == 0 {
		err = errors.New("url and key are required parameters")
	} else {
		response, err = query(ctx, userKey, url)
	}

	res, err = json.Marshal(response)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
