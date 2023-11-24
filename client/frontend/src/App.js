import React, { useState } from 'react';
import logoImage from './logo192.png';
import './App.css';

function App() {
  const [listOfLists, setlistOfLists] = useState([]);

  const [actualList, setActualList] = useState(null);

  const [selectedItems, setSelectedItems] = useState([]);

  const addNewList = () => {
    // if (listOfLists.length < 10) {
      const newlistOfLists = [
        ...listOfLists,
        { title: `List ${listOfLists.length + 1}`, items: [] }
      ];
      setlistOfLists(newlistOfLists);
    // } else {
      // alert("You can only have 10 lists!");
    // }
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

  return (
    <div className="container">
      <div className="content-left">
        <div className="logo-container">
          <img src={logoImage} alt="Logo image" className="logo" />
          <h1 className="title">List Llama</h1>
        </div>
        <div className="lists">
          <h2 className="lists-title">Your Lists</h2>
          {listOfLists.length > 0 ? (
            <div className="list-of-lists">
              {listOfLists.map((list, index) => (
                <div key={index}>
                  <button className="list-button" onClick={() => selectList(list)}>
                    {list.title}
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className="empty-message">You don't have any lists yet, create one below!</p>
          )}
          <button className="button-list" onClick={addNewList}>
            + Add List
          </button>
        </div>
      </div>
      <div className="vertical-line"></div>
      <div className="content-right">
        <div className="right-text">
          {actualList && (
            <>
              <button className="button-item" onClick={addNewItem}>
                + Add Item
              </button>
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
