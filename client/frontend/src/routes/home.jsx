import React, { useEffect, useState } from 'react';
import logoImage from '../images/logo192.png';
import ModalCreate from '../components/ModalCreate';
import ModalAdd from '../components/ModalAdd';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUpload, faDownload, faCopy } from '@fortawesome/free-solid-svg-icons';
import { AMQPClient } from '@cloudamqp/amqp-client' // @ts-ignore


export default function HomePage() {
  const backendIP = "http://localhost:8082"

  const [listOfLists, setlistOfLists] = useState([]);

  const [actualList, setActualList] = useState(null);

  const [item, setItem] = useState("")

  const [quantity, setQuantity] = useState(0)
  const [newItemQuantity, setNewItemQuantity] = useState(0);

  const addNewItem = async (list) => {
    const res = await fetch(`${backendIP}/lists/${list.url}/add`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify({"name": item, "done": false, "quantity": parseInt(newItemQuantity, 10), "list": list}),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    });

    const itemObj = await res.json()
    console.log(itemObj)

    setActualList({
      ...list,
      items: [...(list.items ?? []), itemObj]
    })
    console.log(actualList)
    setItem("");
    setQuantity(0);
  };

  const deleteItem = async (item) => {
    await fetch(`${backendIP}/lists/${actualList.url}/remove`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify({"name": item.name }),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    });

    setActualList({
      ...actualList,
      items: actualList.items.filter((i) => item !== i)
    })
  }

  const selectList = async (list) => {
    console.log(list)
    const items = await (await fetch(`${backendIP}/lists/${list.url}`, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })).json()

    setActualList({
      ...list,
      items: items
    })
  };

  const deleteList = async (list) => {
    await fetch(`${backendIP}/lists/remove`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify({
        "name": list.name,
        "url": list.url
      }),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    });

    setlistOfLists(listOfLists.filter((l) => l !== list))
    setActualList(null)
  };

  useEffect(() => {
    getLists()
  }, [])

  // useEffect(() => {
  //   setActualList(listOfLists[listOfLists.length-1])
  // }, [listOfLists])

  const getLists = async () => {
    
    let lists = await fetch(`${backendIP}/lists`, {method: "GET", mode: "cors", credentials: "include"})
    if (lists.status === 401) {
      document.location = "/login"
      return
    }
    
    lists = await lists.json()

    if (lists != null) 
      setlistOfLists(lists)
  }

  const handleCheckboxChange = (index) => {
    const updatedItems = [...actualList.items];
    updatedItems[index].done = !updatedItems[index].done;
    setActualList({
      ...actualList,
      items: updatedItems,
    });
  };

  const handlePush = async (list) => {
    console.log(list)
    
    await fetch(`${backendIP}/lists/${list.url}/upload`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: '',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    })
    
  };

  const handlePull = async (list) => {
    const res = await fetch(`http://localhost:8082/lists/${list.url}/fetch`, {
       method: 'POST',
       mode: 'cors',
       credentials: 'include',
       body: '',
       headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
     })

     //const updatedList = await res.json()

	  //setActualList(updatedList);
	  window.location.reload();
  };

  const handleCopyUrl = (list) => {
    const listUrl = `${list.url}`;
    navigator.clipboard.writeText(listUrl)
      .then(() => alert('URL copied to clipboard'))
      .catch((err) => console.error('Failed to copy URL', err));
  };

  const updateItem = async (item, updatedQuantity) => {
    const res = await fetch(`${backendIP}/lists/${actualList.url}/update`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify({
        "name": item.name,
        "quantity": parseInt(updatedQuantity, 10),
      }),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    });
  
    const updatedItem = await res.json();
    
    setActualList({
      ...actualList,
      items: actualList.items.map((i) => (i.name === updatedItem.name ? updatedItem : i)),
    });
  };

  return (
    <div className='h-screen'>
      <div className='grid grid-cols-[25%_auto] grid-rows-[15%_auto] grid-flow-row h-full'>
      <div className="col-span-2 col-start-2 row-start-1 border-l border-black h-full"></div>
        <div className='row-span-1 col-span-1 col-start-1 row-start-1'>
            <div className='flex flex-row mt-2'>
              <img src={logoImage} alt="Logo image" className="w-12 h-12 ml-3" />
              <h1 className="text-2xl font-semibold ml-2 p-3">List Llama</h1>
            </div>
        </div>
        <div className='col-span-1 col-start-1 row-start-2 mb-2 ml-3 mr-3'>
          <div className="flex flex-col justify-evenly h-full">
            <h2 className="flex font-semibold">My Lists</h2>
            <div className='flex flex-col justify-between h-full'>
              <div className="ml-1 flex flex-col gap-1 mt-1">
                {
                  listOfLists.length === 0 ? 
                    <div>
                      You have no shopping lists yet
                    </div> : 
                    listOfLists.map((list, index) => (
                      <div key={index} className='flex flex-row justify-between bg-pink-50 p-2.5 rounded-md'>
                        <button className='flex items-center' onClick={() => selectList(list)}>{list.name}</button>
                        <button className='flex p-2 bg-pink-200 rounded-md' onClick={() => deleteList(list)}>Delete</button>
                      </div>
                ))}
              </div>
              <div className='flex mb-3 justify-center flex-col gap-2'>
                <ModalAdd lists={listOfLists} setLists={setlistOfLists} />
                <ModalCreate lists={listOfLists} setLists={setlistOfLists} />
              </div>
            </div>
          </div>
        </div>
        <div className="col-span-2 col-start-2 row-start-2 border-l border-black h-full"></div>
        <div className='col-start-2 col-span-2 row-start-1 mt-10'>
          <div className='flex justify-center'>
              <div className="flex flex-col justify-center gap-2 mx-4">
                {actualList && (
                  <>
                    <div className='grid grid-cols-4 gap-2 grid-flow-col items-center'>
                      <div className='flex flex-col justify-center items-center'>
                        <div className='col-start-1 col-span-1'>
                          <button className='flex p-2 bg-pink-200 rounded-md align-center' onClick={() => handleCopyUrl(actualList)}><FontAwesomeIcon icon={faCopy} /></button>
                        </div>
                      </div>
                      <h1 className="font-semibold col-start-2 col-span-2 text-center mb-5 mt-5 text-xl">{actualList.name}</h1>
                      <div className='flex flex-row justify-end gap-1 col-start-4 col-span-1'>
                        <button className='flex p-2 bg-pink-200 rounded-md' onClick={() => handlePush(actualList)}>
                          <FontAwesomeIcon icon={faUpload} />
                        </button>
                        <button className='flex p-2 bg-pink-200 rounded-md' onClick={() => handlePull(actualList)}>
                          <FontAwesomeIcon icon={faDownload} />
                        </button>
                      </div>
                    </div>
                   
                    <div className='grid grid-cols-4 gap-2 grid-flow-col'>
                      <span className='grid font-bold justify-center'>Bought</span>
                      <span className='grid font-bold justify-center'>Item Name</span>
                      <span className='grid font-bold justify-center'>Quantity</span>
                      <span className='grid font-bold justify-center'>Action</span>
                    </div>
                    { actualList.items ? 
                      actualList.items.map((item, index) => {
                        console.log(item);
                        return (
                          <div className={`flex flex-row justify-between ${item.done ? 'line-through' : ''} grid grid-flow-col grid-cols-4 gap-5`} key={index}>
                            <div className={`grid row-start-${index + 1} justify-center`}>
                              <input type="checkbox" checked={item.done} onChange={() => handleCheckboxChange(index)}/>
                            </div>
                            <div className={`${item.done ? 'text-gray-700' : ''} grid row-start-${index + 1} justify-center`}>{item.name}</div>
                            <div className={`grid row-start-${index + 1}`}>
                            <input
                              className={`${item.done ? 'bg-pink-100 text-gray-700' : 'bg-pink-200'} justify-center p-1 rounded-md text-center`}
                              type='number'
                              value={quantity !== 0 ? quantity : item.quantity}
                              onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                  updateItem(item, quantity);
                                }
                              }}
                              onChange={(e) => setQuantity(e.target.value)}
                            />
                            </div>
                            <div className={`grid row-start-${index + 1}`}>
                              <button className={`${item.done ? 'bg-pink-100 text-gray-700' : 'bg-pink-100'} p-1 rounded-md`} onClick={() => deleteItem(item)}>Delete</button>
                            </div>
                          </div>
                        );
                      })
                    : <></> }
                  </>
                )}
                {
                  actualList && 
                    <div className='grid grid-cols-4 gap-20 mt-7'>
                      <div className='grid col-start-1 col-auto'></div>
                      <div className='grid col-start-2 justify-center col-span-1'><input className='rounded-md p-1 justify-center text-center bg-pink-300' type='text' id='itemName' value={item} placeholder='name' onChange={(e) => setItem(e.target.value)}></input></div>
                      <div className='grid col-start-3 justify-center col-span-1'><input className='rounded-md bg-pink-300 p-1 text-center max-w-s' type='number' id='itemQuantity' value={newItemQuantity} onChange={(e) => setNewItemQuantity(e.target.value)}></input></div>
                      <div className='grid col-start-4 col-span-1'><button className=" bg-pink-300 p-1 rounded-md" onClick={() => addNewItem(actualList)}>Add Item</button></div>
                    </div>
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    );
}
