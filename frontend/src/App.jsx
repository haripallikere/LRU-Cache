import React, { useState, useEffect } from "react";
import axios from "axios";
import "./App.css";

const API_BASE_URL = "http://localhost:8080";
const CACHE_CAPACITY = 5;

function App() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [expiration, setExpiration] = useState(5);
  const [result, setResult] = useState("");
  const [error, setError] = useState("");
  const [cacheItems, setCacheItems] = useState([]);

  useEffect(() => {
    fetchAllItems();
  }, []);

  const fetchAllItems = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/cache`);
      setCacheItems(response.data);
      console.log("Cache items fetched:", response);
    } catch (err) {
      console.error("Failed to fetch cache items:", err);
    }
  };

  const handleGet = async () => {
     console.log(key,'key')
    try {
      const response = await axios.get(`${API_BASE_URL}/cache/${key}`);
      console.log(response,'GET')
      setResult(JSON.stringify(response.data, null, 2));
    } catch (error) {
      console.error(error,'error')
      setResult(JSON.stringify(error.response.data, null, 2));
    }
  };

  const handleSet = async () => {
    console.log("set");
    if (!key ||!value) {
      setError("Key and value are both required");
      return;
    }
    try {
      await axios.post(`${API_BASE_URL}/cache`, {
        key,
        value,
        expiration: Number(expiration),
      });
      setResult("Key set successfully");
      console.log("Key set successfully",key);
    } catch (err) {
      console.error(err);
      setError(err.response?.data?.error || "An error occurred");
      setResult("");
    }
    fetchAllItems();
  };

  const handleDelete = async () => {
    try {
      const response = await axios.delete(`${API_BASE_URL}/cache/${key}`);
      console.log(response,'response')
      setResult("deleted successfully");

    } catch (err) {
      setError(err.response?.data?.error || "An error occurred");
      setResult("");
    }
    fetchAllItems();
  };

  return (
    <div className='App'>
      <div>
        <div>
          <input
            type='text'
            value={key}
            onChange={(e) => setKey(e.target.value)}
            placeholder='Key'
          />
          <input
            type='text'
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder='Value'
          />
          <input
            type='number'
            value={expiration}
            onChange={(e) => setExpiration(e.target.value)}
            placeholder='Expiration (seconds)'
          />
        </div>
        <div className="btn-container">
          <button onClick={handleGet}>Get</button>
          <button onClick={handleSet}>Set</button>
          <button onClick={handleDelete}>Delete</button>

        </div>
 
        {result && <div className="prompt success">Result: {result}</div>}
        {error && <div className="prompt error">Error: {error}</div>}
      
       
        <h2>
          Cache Contents ({cacheItems.length}/{CACHE_CAPACITY})
        </h2>
        <ul>
          {cacheItems.map((item) => (
            <li key={item.key}>
              {item.key}: {JSON.stringify(item.value)} (Expires:{" "}
              {new Date(item.expiresAt).toLocaleString()})
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
