import "./App.css";
import Home from "./page/Home";
import CssBaseline from "@material-ui/core/CssBaseline";
import { ReactComponent as BGSVG } from "./assets/dashboard.svg";

function App() {
  return (
    <div className="App">
      <CssBaseline />
      <BGSVG style={{ maxHeight: "50vh", width: "100%" }} />
      <div className="Main">
        <Home />
      </div>
    </div>
  );
}

export default App;
