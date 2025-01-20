import styles from "./App.module.scss";
import HeaderComponent from "./components/header/Header";
import SideBar from "./components/sideBar/SideBar";
import MediaComponent from "./components/media/Media";
import MassEditor from "./components/massEditor/massEditor";
import Queue from "./components/queue/Queue";
import History from "./components/history/History";
import Settings from "./components/settings/Settings";
import Profiles from "./components/profiles/Profiles";
import General from "./components/general/General";
import Status from "./components/status/Status";
import { useContext, useEffect, useState } from "react";
import {
  BrowserRouter as Router,
  Route,
  Routes,
  useNavigate,
  useParams,
  useLocation,
} from "react-router-dom";
import Series from "./components/series/Series";
import { getTheme } from "./styles/themes";
import Events from "./components/events/Events";
import Authenticaton from "./components/authentication/Authentication";
import { ThemeContext } from "./contexts/themeContext";
import Movie from "./components/movie/Movie";

function App() {
  const [selectedItem, setSelectedItem] = useState<any>(null);
  const [selectedOption, setSelectedOption] = useState<any>(null);
  const theme = useContext(ThemeContext);
  const t = getTheme(theme);
  if (t) {
    Object.entries(t).forEach(([key, value]) => {
      document.documentElement.style.setProperty(`--${key}`, value);
    });
  }
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loaded, setLoaded] = useState(false);
  useEffect(() => {
    const fetchToken = async () => {
      const token = localStorage.getItem("token");
      if (token) {
        const response = await fetch(`/api/auth/logintoken`, {
          method: "POST",
          headers: {
            Authorization: token,
          },
        });

        const data = await response;
        setIsLoggedIn(data.ok);
      }
      setLoaded(true);
    };

    fetchToken();
  }, []);

  if (!loaded) {
    return null;
  }
  return !isLoggedIn ? (
    <Authenticaton />
  ) : t ? (
    <Router>
      <div className={styles.app}>
        <HeaderComponent />
        <div className={styles.content}>
          <div className={styles.sideBar}>
            <SideBar
              selectedOption={selectedOption}
              setSelectedOption={setSelectedOption}
              selectedItem={selectedItem}
              setSelectedItem={setSelectedItem}
            />
          </div>
          <Page
            setSelectedOption={setSelectedOption}
            setSelectedItem={setSelectedItem}
          />
        </div>
      </div>
    </Router>
  ) : null;

  function Page({ setSelectedOption, setSelectedItem }: any) {
    const location: any = useLocation();
    const pathname: string = location.pathname;

    useEffect(() => {
      const sidebar: any = {
        "/": [0, -1],
        "/mass-editor": [0, 0],
        "/activity": [1, -1],
        "/activity/queue": [1, 0],
        "/activity/history": [1, 1],

        "/settings": [2, -1],
        "/settings/profiles": [2, 0],
        "/settings/general": [2, 1],
        "/system": [3, -1],
        "/system/status": [3, 0],
        "/system/events": [3, 1],
      };
      setSelectedOption(pathname in sidebar ? sidebar[pathname][0] : 0);
      setSelectedItem(pathname in sidebar ? sidebar[pathname][1] : -1);
    }, [pathname, setSelectedOption, setSelectedItem]);

    return (
      <Routes>
        <Route path="/" element={<MediaComponent />} />
        <Route path="/series/:seriesName" element={<SeriesSelect />} />
        <Route path="/movies/:movieName" element={<MovieSelect />} />
        <Route path="/mass-editor" element={<MassEditor />} />

        <Route path="/activity" element={<Activity />} />
        <Route path="/activity/queue" element={<Queue />} />
        <Route path="/activity/history" element={<History />} />

        <Route path="/settings" element={<Settings />} />
        <Route path="/settings/profiles" element={<Profiles />} />
        <Route path="/settings/general" element={<General />} />
        <Route path="/system" element={<Status />} />

        <Route path="/system/status" element={<Status />} />
        <Route path="/system/events" element={<Events />} />

        {<Route path="*" element={<NotFound />} />}
      </Routes>
    );
  }

  function Activity() {
    const navigate = useNavigate();

    useEffect(() => {
      navigate("/activity/queue");
    }, [navigate]);
    return null;
  }
  function SeriesSelect() {
    const { seriesName } = useParams();
    return <Series seriesName={seriesName} />;
  }
  function MovieSelect() {
    const { movieName } = useParams();
    return <Movie movieName={movieName} />;
  }

  function NotFound() {
    const navigate = useNavigate();
    useEffect(() => {
      navigate("/");
    }, [navigate]);
    return null;
  }
}

export default App;
