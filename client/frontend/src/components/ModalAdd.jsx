import { useState } from "react"
import Popup from 'reactjs-popup'
import '../styles/Modal.css'

export default function ModalAdd({ lists, setLists }) {
    const backendIP = "http://localhost:8082"
    const [url, setUrl] = useState("")

    async function addList() {
        const response = await fetch(`${backendIP}/lists/${url}/fetch`);

        if (response.ok) {
            console.log("New shopping list added successfully")
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

