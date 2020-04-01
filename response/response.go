package response

/* code response
* 200 Success
* 400 Bad Request if failed proccessing data
* 401 Unauthorized
* 422 Validation error
 */

import (
	"encoding/json"
	"net/http"
)

type unauthorized struct {
	Msg string `json:"message"`
}

type responseJSON struct {
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type responseErrorJSON struct {
	Msg   string      `json:"message"`
	Error interface{} `json:"error"`
}

type Error struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
type responseErrorValidate struct {
	Msg    string  `json:"message"`
	Errors []Error `json:"errors"`
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

// Error
func RespondErrorJSON(w http.ResponseWriter, payload interface{}, msg ...string) {

	message := "Error"
	if len(msg) > 0 {
		message = msg[0]
	}

	payload = responseErrorJSON{
		Msg:   message,
		Error: payload,
	}

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(response)
}

// Success
func RespondSuccessJSON(w http.ResponseWriter, payload interface{}, msg ...string) {
	//response, _ := json.Marshal(payload)
	message := "Success"
	if len(msg) > 0 {
		message = msg[0]
	}

	payload = responseJSON{
		Msg:  message,
		Data: payload,
	}

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func RespondErrValidateJSON(w http.ResponseWriter, payload interface{}, msg ...string) {
	//response, _ := json.Marshal(payload)
	message := "InvalidParameter"
	if len(msg) > 0 {
		message = msg[0]
	}

	payload = responseErrorJSON{
		Msg:   message,
		Error: payload,
	}

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	w.Write(response)
}

// Success
func RespondErrorValidateJSON(w http.ResponseWriter, payload []Error, msg ...string) {
	//response, _ := json.Marshal(payload)
	message := "InvalidParameter"
	if len(msg) > 0 {
		message = msg[0]
	}

	res := responseErrorValidate{
		Msg:    message,
		Errors: payload,
	}

	response, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	w.Write(response)
}

// Unauthorized
func RespondUnauthorizedJSON(w http.ResponseWriter, message string) {
	payload := unauthorized{
		Msg: message,
	}

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	w.Write(response)
}

// Unauthorized
func RespondUnverifyJSON(w http.ResponseWriter, message string) {
	payload := unauthorized{
		Msg: message,
	}

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)
	w.Write(response)
}
