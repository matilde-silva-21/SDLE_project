import React, { useEffect, useState } from 'react';
import logoImage from '../images/logo192.png';
import Modal from '../components/Modal';

export default function HomePage() {
  const [listOfLists, setlistOfLists] = useState([]);

  const [actualList, setActualList] = useState(null);

  const [selectedItems, setSelectedItems] = useState([]);

  const [item, setItem] = useState("")

  const addNewItem = async (list) => {
    const res = await fetch(`http://localhost:8080/lists/${list.url}/add`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify({"name": item, "done": false, "list": list}),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      } 
    });

    const itemObj = await res.json()

    setActualList({
      ...list,
      items: [...(list.items ?? []), itemObj]
    })
  };

  const deleteItem = async (item) => {
    await fetch(`http://localhost:8080/lists/${actualList.url}/remove`, {
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
    const items = await (await fetch(`http://localhost:8080/lists/${list.url}`, {
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
    await fetch(`http://localhost:8080/lists/remove`, {
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

  const toggleItemSelection = (index) => {
    if (selectedItems.includes(index)) {
      setSelectedItems(selectedItems.filter((i) => i !== index));
    } else {
      setSelectedItems([...selectedItems, index]);
    }
  };

  useEffect(() => {
    getLists()
  }, [])

  const getLists = async () => {
    
    let lists = await fetch("http://localhost:8080/lists", {method: "GET", mode: "cors", credentials: "include"})
    if (lists.status === 401) {
      document.location = "/login"
      return
    }
    
    lists = await lists.json()

    if (lists != null) 
      setlistOfLists(lists)
  }

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
                      <div key={index} className='flex flex-row justify-between bg-pink-50 p-2 rounded-md'>
                        <button className='flex' onClick={() => selectList(list)}>{list.name}</button>
                        <button className='flex p-2 bg-pink-300 rounded-md' onClick={() => deleteList(list)}>Delete</button>
                      </div>
                ))}
              </div>
              <div className='flex mb-3 justify-center'>
                <Modal lists={listOfLists} setLists={setlistOfLists} className="button-list"/>
              </div>
            </div>
          </div>
        </div>
        <div className="col-span-2 col-start-2 row-start-2 border-l border-black h-full"></div>
        <div className='col-start-2 col-span-2 row-start-1 mt-10'>
          <div className='flex justify-center'>
              <div className="flex flex-col justify-center gap-2">
                {actualList && (
                  <>
                    <h1 className="font-semibold flex justify-center mb-10 text-xl">{actualList.name}</h1>
                    <ul className='flex flex-col gap-2'>
                      {
                        actualList.items ? 
                          actualList.items.map((item, index) => (
                            <li
                              key={index}
                              className={`flex flex-row justify-between ${selectedItems.includes(index) ? 'line-through' : ''}`}
                              onClick={() => toggleItemSelection(index)}
                            >
                              <input type="checkbox" value={item.done}/> {item.name}
                              <button className='flex bg-pink-200 p-1 rounded-md' onClick={() => deleteItem(item)}>Delete</button>
                            </li>
                          ))
                        : <></>}
                    </ul>
                  </>
                )}
                {
                  actualList && 
                  <div className='flex flex-row gap-1'>
                    <input className='flex rounded-md p-1' type='text' id='itemName' value={item} placeholder='name' onChange={(e) => setItem(e.target.value)}></input>
                    <button className="flex bg-pink-200 p-1 rounded-md" onClick={() => addNewItem(actualList)}>Add Item</button>
                  </div>
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    );
}
