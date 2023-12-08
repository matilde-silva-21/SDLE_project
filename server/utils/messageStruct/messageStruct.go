package messageStruct


import (
	"fmt"
	"encoding/json"	
)


type MessageType string

const (
    Write MessageType = "Write"
    Read MessageType = "Read"
    Delete MessageType = "Delete"
	Error MessageType = "Error" // If there was an error reading the server-side copies (user doens't send messages with "Error" Action)
)


type MessageStruct struct {
	ListURL string
	Username string
	Action MessageType
	Body string
}

type ServerMessageStruct struct {
	IPs []string
	Payload MessageStruct
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


func ReadServerMessage(body []byte) ([]string, MessageStruct){

	var serverMessageObject ServerMessageStruct
	
	err := json.Unmarshal(body, &serverMessageObject)
	if(err != nil){
		fmt.Println("Error:", err)
		var dummy MessageStruct
		return []string{}, dummy
	}

	return serverMessageObject.IPs, serverMessageObject.Payload
}
/*


mensagem base - o orchestrator recebe sempre a mensagem assim (as mensagens dos clientes e servidores para o orchestrator têm de vir assim):
{
	"ListURL": "123",
	"Username": "john.doe",
	"Action": "Write" ou "Read" ou "Delete" ou "Error",
	"Body": {
		"{\"Name\":\"My List 1\", \"List\":{\"Map\":{\"apple\":{\"First\":1,\"Second\":3},\"pear\":{\"First\":2,\"Second\":2},\"rice\":{\"First\":3,\"Second\":2}}}, \"State\":{\"Map\":{\"pear\":{\"First\":0,\"Second\":0},\"rice\":{\"First\":2,\"Second\":0}}}}"}
	}

}

mensagem que os servidores recebem - o payload é a mensagem base, os IPs são os servidores com quem falar para fazer quorum

{
	"IPs": ["0.0.0.0:0000", "0.0.0.0:0000"],
	"Payload":
	{
		"ListURL": "123",
		"Username": "john.doe",
		"Action": "Write" ou "Read" ou "Delete" ou "Error",
		"Body": {
			"{\"Name\":\"My List 1\", \"List\":{\"Map\":{\"apple\":{\"First\":1,\"Second\":3},\"pear\":{\"First\":2,\"Second\":2},\"rice\":{\"First\":3,\"Second\":2}}}, \"State\":{\"Map\":{\"pear\":{\"First\":0,\"Second\":0},\"rice\":{\"First\":2,\"Second\":0}}}}"}
		}

	}


}


*/