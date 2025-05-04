import axios from 'axios';

// Create axios instance with default config
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
});

// HTTP methods
export const get = (url, params = {}) => {
  return api.get(url, { params });
};

export const post = (url, data = {}, config = {}) => {
  return api.post(url, data, config);
};

export const put = (url, data = {}, config = {}) => {
  return api.put(url, data, config);
};

export const del = (url) => {
  return api.delete(url);
};

export default api;
