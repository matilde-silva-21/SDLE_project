package messageStruct


import (
	"fmt"
	"encoding/json"	
)


type MessageType string

const (
    Create MessageType = "create"
    Delete MessageType = "delete"
    Update MessageType = "update"
)

type CRDTString struct {
	List string
	State string
}

type MessageStruct struct {
	ListURL string
	Username string
	Action MessageType
	Body CRDTString
}


func CreateMessage(url, username string, action MessageType, list, state string) MessageStruct {

	body := CRDTString{List: list, State: state}
	return MessageStruct{ListURL: url, Username: username, Action: action, Body: body}

}


func (message MessageStruct) ToJSON() []byte{
	b, err := json.Marshal(message)
	if(err != nil){
        fmt.Println("Error:", err)
		return []byte{}
	}

	return b
}


func JSONToMessage(body []byte, message *MessageStruct){
	err := json.Unmarshal(body, *message)
	if(err != nil){
        fmt.Println("Error:", err)
		return
	}
}