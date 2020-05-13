package notebook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/parnurzeal/gorequest"
)

var (
	notebookClient *gorequest.SuperAgent
)

const (
	DeactivateQueryString = "?documentId=%d"
)

func init() {
	notebookClient = common.DefaultRequest.Clone().Timeout(config.C.Push.RequestTimeout)
}

func GetDocuments(ctx common.Context) ([]Customer, error) {
	resp, err := common.UapiAuthorizedRequest(ctx, notebookClient, nil, http.MethodGet,
		config.C.PersonalArea.Urls.DocumentsV2, nil)
	if err != nil {
		return nil, err
	}

	customers := []Customer{}
	err = json.Unmarshal(resp.Data, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func DeactivateDocument(ctx common.Context, documentID int64) (*models.RawResponse, error) {
	url := config.C.PersonalArea.Urls.DocumentsDeactivateV2 + fmt.Sprintf(
		DeactivateQueryString, documentID)

	resp, err := common.UapiAuthorizedRequest(ctx, notebookClient, nil, http.MethodPost,
		url, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CreateDocuments(ctx common.Context, customer Customer) (*Customer, error) {
	url := config.C.PersonalArea.Urls.CustomersV2

	for i := range customer.CustomerDocuments {
		customer.CustomerDocuments[i].IsActive = true
	}

	resp, err := common.UapiAuthorizedRequest(ctx, notebookClient, customer, http.MethodPost,
		url, nil)
	if err != nil {
		return nil, err
	}

	newCustomer := &Customer{}
	err = json.Unmarshal(resp.Data, newCustomer)
	if err != nil {
		return nil, err
	}

	if len(newCustomer.CustomerDocuments) == 0 {
		return nil, fmt.Errorf("Нет документа в ответе от Uapi: %+v", newCustomer)
	}

	newDoc := newCustomer.CustomerDocuments[0]

	allDocs, err := GetDocuments(ctx)
	if err != nil {
		return nil, err
	}

	for _, doc := range allDocs {
		if len(doc.CustomerDocuments) != 0 && doc.CustomerDocuments[0].ID == newDoc.ID {
			return newCustomer, nil
		}
	}

	return nil, fmt.Errorf("Uapi вернул документ, но не записал его в базу")
}
