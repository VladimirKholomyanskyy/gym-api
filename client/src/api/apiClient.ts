import axios from 'axios';

const apiClient = axios.create({
    baseURL: 'http://localhost:8080',  // Replace with your backend URL
    headers: {
      'Content-Type': 'application/json',  // Set other default headers if needed
    },
});

export const setAuthToken = (token: string) => {
    console.log(`Token=${token}`)
  apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
};

export default apiClient;
