import React, { useEffect, useState } from 'react';
import logoImage from './logo192.png';
import './App.css';

function App() {
  const [listOfLists, setlistOfLists] = useState([]);

  const [actualList, setActualList] = useState(null);

  const [selectedItems, setSelectedItems] = useState([]);

  const addNewList = () => {
    const newlistOfLists = [
      ...listOfLists,
      { title: `List ${listOfLists.length + 1}`, items: [] }
    ];
    setlistOfLists(newlistOfLists);
  };

  const addNewItem = () => {
    if (actualList) {
      const newItem = `Item ${actualList.items.length + 1}`;
      const updatedList = { ...actualList, items: [...actualList.items, newItem] };
      const updatedLists = listOfLists.map((list) =>
        list === actualList ? updatedList : list
      );
      setlistOfLists(updatedLists);
      setActualList(updatedList);
    }
  };

  const selectList = (list) => {
    setActualList(list);
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
    const login = await fetch("http://localhost:8080/login", {
      method: 'POST',
      body: JSON.stringify({'username': 'user1'}),
      mode: 'cors',
      credentials: "include",
      headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    }})

    if (login.ok) {
      const lists = await (await fetch("http://localhost:8080/lists", {method: "GET", mode: "cors", credentials: "include"})).json()
      console.log("lists are" + lists)
      return lists
    }

    return []
  }

  return (
    <div className="container">
      <div className="content-left">
        <div className="logo-container">
          <img src={logoImage} alt="Logo image" className="logo" />
          <h1 className="title">List Llama</h1>
        </div>
        <h2 className="lists-title">Your Lists</h2>
        <div className="list-of-lists">
          {listOfLists.map((list, index) => (
            <div key={index}>
              <button onClick={() => selectList(list)}>{list.title}</button>
            </div>
          ))}
        </div>
        <div>
          <button className="button-list" onClick={addNewList}>+ Add List</button>
        </div>
      </div>
      <div className="vertical-line"></div>
      <div className="content-right">
        <div className="right-text">
          {actualList && <button className="button-item" onClick={addNewItem}>+ Add Item</button>}
          {actualList && (
            <>
              <h1 className="list-title">{actualList.title}</h1>
              <div className="horizontal-line"></div>
              <ul className="list-of-items">
                {actualList.items.map((item, index) => (
                  <li
                    key={index}
                    className={selectedItems.includes(index) ? 'strikethrough' : ''}
                    onClick={() => toggleItemSelection(index)}
                  >
                    <input type="checkbox" /> {item}
                  </li>
                ))}
              </ul>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
