import axios from "axios";
import { User } from "oidc-client-ts";

const apiClient = axios.create({
  baseURL: "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

function getUser() {
  const oidcStorage = localStorage.getItem(
    `oidc.user:http://localhost:8070/realms/gainz:react-client`
  );
  if (!oidcStorage) {
    return null;
  }

  return User.fromStorageString(oidcStorage);
}
export const setupAxiosInterceptors = () => {
  apiClient.interceptors.request.use(
    (config) => {
      const user = getUser();
      const token = user?.access_token;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
        console.log("token valid");
      } else {
        console.log("token not valid");
      }
      return config;
    },
    (error) => Promise.reject(error)
  );
};

export default apiClient;
