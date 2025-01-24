import { useEffect, useState } from "react";
import styles from "./Authentication.module.scss";
const Authentication = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [activated, setActivated] = useState(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const fetchUser = async () => {
      const response = await fetch(`/api/auth/activated`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const activated = await response.json();
      setActivated(activated.activated);
      setLoaded(true);
    };

    fetchUser();
  }, []);
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (activated) {
      const token = await fetch(`/api/auth/login`, {
        method: "POST",
        body: JSON.stringify({ username: username, password: password }),
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (token.status === 200) {
        const tokenData = await token.json();
        localStorage.setItem("token", tokenData.token);
        window.location.reload();
      }
    } else {
      await fetch(`/api/auth/register`, {
        method: "POST",
        body: JSON.stringify({ username: username, password: password }),
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      window.location.href = "/";
    }
  };

  if (!loaded) {
    return null;
  }
  return (
    <>
      <div className={styles.authenticate}>
        <form className={styles.form} onSubmit={handleSubmit}>
          <h1>{activated ? <>Login</> : <>Register</>}</h1>
          <div className={styles.inputContainer}>
            <label>Username: </label>
            <input
              className={styles.text}
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className={styles.inputContainer}>
            <label>Password: </label>
            <input
              className={styles.text}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <input className={styles.submit} type="submit" value="Submit" />
        </form>
      </div>
    </>
  );
};

export default Authentication;
