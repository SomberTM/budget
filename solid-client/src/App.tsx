import type { Component } from "solid-js";
import { Api } from "./api/";

const App: Component = () => {
  Api.users
    .getCurrentUser()
    .then((response) =>
      console.log(
        "Api response",
        response.success
          ? response.data.user_name
          : response.detailedError?.message
      )
    );
  return <div></div>;
};

export default App;
