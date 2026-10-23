package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/julienschmidt/httprouter"
)

type NewHandler struct {
	Client *dynamodb.Client
}

func (h *NewHandler) GetProductByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	result, err := h.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("products"),
		Key: map[string]types.AttributeValue{
			"product_id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	if err != nil {
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}
	if result.Item == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	product := Product{
		ProductId: id,
	}

	if v, ok := result.Item["name"].(*types.AttributeValueMemberS); ok {
		product.Name = v.Value
	}
	if v, ok := result.Item["category"].(*types.AttributeValueMemberS); ok {
		product.Catogary = v.Value
	}
	if v, ok := result.Item["price"].(*types.AttributeValueMemberN); ok {
		product.Price, _ = strconv.ParseFloat(v.Value, 64)
	}

	if v, ok := result.Item["quantity"].(*types.AttributeValueMemberN); ok {
		product.Quantity, _ = strconv.Atoi(v.Value)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(product)
}

func (h *NewHandler) GetAllProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result, err := h.Client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("products"),
	})

	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	products := []Product{}

	for _, item := range result.Items {
		product := Product{}

		if v, ok := item["name"].(*types.AttributeValueMemberS); ok {
			product.Name = v.Value
		}
		if v, ok := item["category"].(*types.AttributeValueMemberS); ok {
			product.Catogary = v.Value
		}
		if v, ok := item["price"].(*types.AttributeValueMemberN); ok {
			product.Price, _ = strconv.ParseFloat(v.Value, 64)
		}

		if v, ok := item["quantity"].(*types.AttributeValueMemberN); ok {
			product.Quantity, _ = strconv.Atoi(v.Value)
		}

		products = append(products, product)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(products)
}
func (h *NewHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err = h.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("products"),
		Key: map[string]types.AttributeValue{
			"product_id": &types.AttributeValueMemberS{
				Value: product.ProductId,
			},
		},

		UpdateExpression: aws.String(
			"SET #n= :name, category= :category,price = :price , quantity= :quantity",
		),
		ExpressionAttributeNames: map[string]string{
			"#n": "name",
		},

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":     &types.AttributeValueMemberS{Value: product.Name},
			":catogary": &types.AttributeValueMemberS{Value: product.Catogary},
			":price":    &types.AttributeValueMemberS{Value: strconv.FormatFloat(product.Price, 'f', -1, 64)},
			":quantity": &types.AttributeValueMemberS{Value: strconv.Itoa(product.Quantity)},
		},
	})

	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product update suceesfully",
	})
}

func (h *NewHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	_, err := h.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("products"),
		Key: map[string]types.AttributeValue{
			"product_id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "product delete suceesfully",
	})

}
