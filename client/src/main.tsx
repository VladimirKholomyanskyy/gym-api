import { Provider } from "@/components/ui/provider";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { AuthProvider } from "react-oidc-context";
import { User, WebStorageStateStore } from "oidc-client-ts";

const oidcConfig = {
  authority: "http://localhost:8070/realms/gainz",
  client_id: "react-client",
  redirect_uri: window.location.origin,
  post_logout_redirect_uri: "http://localhost:5173",
  onSigninCallback: (_user: User | void): void => {
    window.history.replaceState(
      {},
      document.title,
      window.location.pathname // Clean up the URL by removing query parameters
    );
  },
  userStore: new WebStorageStateStore({ store: window.localStorage }),
};

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider {...oidcConfig}>
      <Provider>
        <App />
      </Provider>
    </AuthProvider>
  </StrictMode>
);
