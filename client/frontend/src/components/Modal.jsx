import { useState } from "react"
import Popup from 'reactjs-popup'
import '../styles/Modal.css'

export default function Modal({ lists, setLists }) {
    const [url, setUrl] = useState("")
    const [name, setName] = useState("")

    async function createList() {
        const response = await fetch("http://localhost:8080/lists/create", {
            method: "POST",
            body: JSON.stringify({"url": url, "name": name}),
            mode: "cors",
            credentials: "include",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }
        })

        if (response.ok) {
            console.log("New shopping list created successfully")
            let res = await response.json()
            console.log(res)
            setLists([...lists, res])
        }
    }
    
    return (        
          <Popup
            trigger={<button className="button bg-pink-200 p-2 rounded-md w-full"> Add List </button>}
            modal
            nested
          >
            {close => (
              <div className="modal rounded-md p-2 bg-pink-200">
                <button className="close" onClick={close}>
                  &times;
                </button>
                <div className="header">Add List</div>
                <div className="content">
                    <div className="flex">
                        <form className="flex flex-col gap-2 justify-center align-middle">
                            <label htmlFor="url">Url</label>
                            <input type="text" value={url} placeholder="List url" className="rounded-md p-2" name="url" id="url" onChange={(e) => setUrl(e.target.value)}></input>
                            <label htmlFor="name">Name</label>
                            <input type="text" value={name} placeholder="List name" className="rounded-md p-2" name="name" id="name" onChange={(e) => setName(e.target.value)}></input>
                        </form>
                    </div>
                </div>
                <div className="actions">
                  <button
                    className="button"
                    onClick={() => {
                      console.log('modal closed ');
                      createList()
                      close();
                    }}
                  >
                    Add List
                  </button>
                </div>
              </div>
            )}
          </Popup>
        )

}

