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


func JSONToMessage(body []byte) MessageStruct{
	var message MessageStruct

	err := json.Unmarshal(body, &message)
	if(err != nil){
        fmt.Println("Error:", err)
	}
	return message
}

func (message MessageStruct) BuildMessageForServer(IPaddresses []string) []byte{

	payload := message.ToJSON()

	jsonArray, _ := json.Marshal(IPaddresses)

	IPstrings := fmt.Sprintf("\"IPs\": %s", jsonArray)

	return []byte(fmt.Sprintf("{%s, \"Payload\":%s}", IPstrings, payload))
}

/*


mensagem base - o orchestrator recebe sempre a mensagem assim (as mensagens dos clientes e servidores para o orchestrator têm de vir assim):
{
	"ListURL": "123",
	"Username": "john.doe",
	"Action": "Create" ou "Delete" ou "Update",
	"Body": {TBD}

}

mensagem que os servidores recebem - o payload é a mensagem base, os IPs são os servidores com quem falar para fazer quorum

{
	"IP": ["0.0.0.0:0000", "0.0.0.0:0000"],
	"Payload":
	{
		"ListURL": "123",
		"Username": "john.doe",
		"Action": "Create" ou "Delete" ou "Update",
		"Body": {TBD}

	}


}


*/