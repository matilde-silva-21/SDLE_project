import React, { useEffect, useState } from 'react';
import logoImage from '../images/logo192.png';
import '../styles/App.css';
import Modal from '../components/Modal';

export default function HomePage() {
  const [listOfLists, setlistOfLists] = useState([]);

  const [actualList, setActualList] = useState([]);

  const [selectedItems, setSelectedItems] = useState([]);

  const [modalEnabled, setModalEnabled] = useState(false)

  const addNewList = () => {
    console.log("click")
    setModalEnabled(true)
    /*const newlistOfLists = [
      ...listOfLists,
      { title: `List ${listOfLists.length + 1}`, items: [] }
    ];
    setlistOfLists(newlistOfLists);*/
  };

  const addNewItem = () => {
    /*if (actualList) {
      const newItem = `Item ${actualList.items.length + 1}`;
      const updatedList = { ...actualList, items: [...actualList.items, newItem] };
      const updatedLists = listOfLists.map((list) =>
        list === actualList ? updatedList : list
      );
      setlistOfLists(updatedLists);
      setActualList(updatedList);
    }*/
  };

  const selectList = async (list) => {
    console.log(list.url)
    const items = await (await fetch(`http://localhost:8080/lists/${list.url}`, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })).json()

    console.log(items)

    setActualList(items)
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
      <div className='grid grid-cols-[23%_auto] grid-rows-[15%_auto] grid-flow-row h-full'>
        <div className='row-span-1 col-span-1 col-start-1 row-start-1'>
            <div className='flex flex-row mt-2'>
              <img src={logoImage} alt="Logo image" className="w-12 h-12 ml-3" />
              <h1 className="text-2xl font-semibold ml-2 p-3">List Llama</h1>
            </div>
        </div>
        <div className='col-span-1 col-start-1 row-start-2 mb-2 ml-3'>
          <div className="flex flex-col justify-evenly h-full">
            <h2 className="flex font-semibold">My Lists</h2>
            <div className='flex flex-col justify-between h-full'>
              <div className="ml-3">
                {
                  listOfLists.length === 0 ? 
                    <div>
                      You have no shopping lists yet
                    </div> : 
                    listOfLists.map((list, index) => (
                      <div key={index}>
                        <button onClick={() => selectList(list)}>{list.name}</button>
                      </div>
                ))}
              </div>
              <div className='flex mb-3 justify-center'>
                <Modal className="button-list"/>
              </div>
            </div>
          </div>
        </div>
        <div className='col-start-2 col-span-2 row-start-2'>
              <div className="flex flex-row justify-center">
                {actualList && (
                  <>
                    <h1 className="font-semibold">{actualList.title}</h1>
                    <ul>
                      {actualList.map((item, index) => (
                        <li
                          key={index}
                          className={selectedItems.includes(index) ? 'line-through' : ''}
                          onClick={() => toggleItemSelection(index)}
                        >
                          <input type="checkbox" value={item.done}/> {item.name}
                        </li>
                      ))}
                    </ul>
                  </>
                )}
                {actualList && <button className="flex" onClick={addNewItem}>+ Add Item</button>}
              </div>
          </div>
        </div>
      </div>
    );
}
