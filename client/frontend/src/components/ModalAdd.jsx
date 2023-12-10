import { useState } from "react"
import Popup from 'reactjs-popup'
import '../styles/Modal.css'

export default function ModalAdd({ lists, setLists }) {
    const numOfInstances = window.location.host.split(':')[1] - 5173
    const port = parseInt(import.meta.env.VITE_BACKEND_PORT) + numOfInstances
    const backendIP = `http://localhost:${port}`
    const [url, setUrl] = useState("")

    async function addList() {
        const response = await fetch(`${backendIP}/lists/${url}/fetch`, {
          method: 'POST',
          mode: 'cors',
          credentials: 'include',
          body: '',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
          } 
        })

        if (response.ok) {
            console.log("New shopping list added successfully")
            let res = await response.json()
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
                            <input type="text" placeholder="List url" className="rounded-md p-2" name="url" id="url" onChange={(e) => setUrl(e.target.value)}></input>
                        </form>
                    </div>
                </div>
                <div className="actions">
                  <button
                    className="button"
                    onClick={() => {
                      console.log('modal closed ');
                      addList()
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

