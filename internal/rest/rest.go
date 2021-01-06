package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dannypaul/go-skeleton/internal/exception"
	"github.com/dannypaul/go-skeleton/internal/header"
)

func DecodeReq(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if strings.Split(r.Header.Get(header.ContentType), ";")[0] != "application/json" {
		message := "Content-Type header is not application/json"
		http.Error(w, message, http.StatusUnsupportedMediaType)
		return errors.New(message)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			message := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, message, http.StatusBadRequest)
			break

		case errors.Is(err, io.ErrUnexpectedEOF):
			message := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, message, http.StatusBadRequest)
			break

		case errors.As(err, &unmarshalTypeError):
			message := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, message, http.StatusBadRequest)
			break

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, message, http.StatusBadRequest)
			break

		case errors.Is(err, io.EOF):
			message := "Request body must not be empty"
			http.Error(w, message, http.StatusBadRequest)
			break

		case err.Error() == "http: request body too large":
			message := "Request body must not be larger than 1MB"
			http.Error(w, message, http.StatusBadRequest)
			break
		}

		return err
	}

	return nil
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorRes struct {
	Errors    []Error `json:"errors"`
	RequestId string  `json:"requestId"`
}

func EncodeRes(w http.ResponseWriter, r *http.Request, res interface{}, err error) {
	w.Header().Set(header.ContentType, "application/json")
	if err != nil {
		errList := ErrorRes{
			Errors:    []Error{{Message: exception.Message(err.Error()), Code: err.Error()}},
			RequestId: r.Context().Value("requestId").(string),
		}

		errListJson, _ := json.Marshal(errList)
		log.Print(string(errListJson))

		status := exception.HttpStatus(err.Error())
		w.WriteHeader(status)

		if status == http.StatusInternalServerError {
			errList = ErrorRes{
				Errors:    []Error{{Message: exception.Message(err.Error()), Code: http.StatusText(status)}},
				RequestId: r.Context().Value("requestId").(string),
			}
		}

		json.NewEncoder(w).Encode(errList)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
