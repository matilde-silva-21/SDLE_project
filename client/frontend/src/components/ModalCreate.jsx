import { useState } from "react"
import Popup from 'reactjs-popup'
import '../styles/Modal.css'
import { v4 as uuidv4 } from 'uuid';

export default function ModalCreate({ lists, setLists }) {
    const [name, setName] = useState("")

    async function createList() {
      const url = uuidv4();
      const response = await fetch("http://localhost:8080/lists/create", {
          method: "POST",
          body: JSON.stringify({ "name": name, "url": url }),
          mode: "cors",
          credentials: "include",
          headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
          }
      })

      if (response.ok) {
        let res = await response.json()
        setLists([...lists, res])
      }
  }
    
    return (        
          <Popup
            trigger={<button className="button bg-pink-200 p-2 rounded-md w-full"> Create List </button>}
            modal
            nested
          >
            {close => (
              <div className="modal rounded-md p-2 bg-pink-200">
                <button className="close" onClick={close}>
                  &times;
                </button>
                <div className="header">Create List</div>
                <div className="content">
                    <div className="flex">
                        <form className="flex flex-col gap-2 justify-center align-middle">
                            <label htmlFor="name">Name</label>
                            <input type="text" placeholder="List name" className="rounded-md p-2" name="name" id="name" onChange={(e) => setName(e.target.value)}></input>
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
                    Create List
                  </button>
                </div>
              </div>
            )}
          </Popup>
        )

}

