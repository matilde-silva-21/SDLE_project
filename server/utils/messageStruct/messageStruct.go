package messageStruct


import (
	"fmt"
	"encoding/json"	
)


type MessageType string

const (
    Create MessageType = "Create"
    Delete MessageType = "Delete"
    Update MessageType = "Update"
	Add MessageType = "Add" // User adds a list by URL
)


type MessageStruct struct {
	ListURL string
	Username string
	Action MessageType
	Body string
}



func CreateMessage(url, username string, action MessageType, CRDT string) MessageStruct {

	return MessageStruct{ListURL: url, Username: username, Action: action, Body: CRDT}

}


func (message MessageStruct) ToJSON() []byte{
	b, err := json.Marshal(message)
	if(err != nil){
        fmt.Println("Error:", err)
		return []byte{}
	}

	return b
}


func JSONToMessage(body []byte) (MessageStruct,error){
	var message MessageStruct

	err := json.Unmarshal(body, &message)
	if(err != nil){
        //fmt.Println("Error:", err)
		return message, err
	}
	return message, nil
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
	"Body": {
		"{\"Name\":\"My List 1\", \"List\":{\"Map\":{\"apple\":{\"First\":1,\"Second\":3},\"pear\":{\"First\":2,\"Second\":2},\"rice\":{\"First\":3,\"Second\":2}}}, \"State\":{\"Map\":{\"pear\":{\"First\":0,\"Second\":0},\"rice\":{\"First\":2,\"Second\":0}}}}"}
	}

}

mensagem que os servidores recebem - o payload é a mensagem base, os IPs são os servidores com quem falar para fazer quorum

{
	"IP": ["0.0.0.0:0000", "0.0.0.0:0000"],
	"Payload":
	{
		"ListURL": "123",
		"Username": "john.doe",
		"Action": "Create" ou "Delete" ou "Update" ou "Add",
		"Body": {
			"{\"Name\":\"My List 1\", \"List\":{\"Map\":{\"apple\":{\"First\":1,\"Second\":3},\"pear\":{\"First\":2,\"Second\":2},\"rice\":{\"First\":3,\"Second\":2}}}, \"State\":{\"Map\":{\"pear\":{\"First\":0,\"Second\":0},\"rice\":{\"First\":2,\"Second\":0}}}}"}
		}

	}


}


*/