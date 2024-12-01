import { useEffect, useState } from "react";
import axios from "axios";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Link from "next/link";

function Home() {
  const [active, setActive] = useState(false);
  const [SecretKey, setSecretKey] = useState("");
  const [InvalidSecret, setInvalidSecret] = useState(false)
  const [apiKey, setApiKey] = useState("");
  const [fetching, setFetching] = useState(false);

  const notify = () =>
    toast("✅ API Key Copied", {
      position: "top-right",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "dark",
    });


  useEffect(() => {
    if (apiKey !== "") {
      notify();
      navigator.clipboard.writeText(apiKey).catch((error) => console.error("Failed to copy:", error));
    }
  }, [apiKey]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFetching(true);
    const data = {
      secret: SecretKey,
    };

    try {
      const res = await axios.post(
        `${process.env.NEXT_PUBLIC_BACKEND_API_KEY}/accessKey`,
        data
      );
      setApiKey(res.data.data);
      setActive(true);
    } catch (err) {
      console.log(err);
      setInvalidSecret(true);
    } finally {
      setFetching(false);
    }
  };

  return (
    <div className="master">
      <ToastContainer />
      <Link href="/">
        <button>Home</button>
      </Link>
      <div className="login-box">
        <form className="formItems">
          <div className="label">
            <label id="text">SecretKey :</label>
            <input
              className="input"
              type="text"
              name="secretKey"
              value={SecretKey}
              onChange={(e) => {
                setSecretKey(e.target.value);
              }}
            />
          </div>

          
          <button
            disabled={active}
            id="button"
            onClick={(e) => handleSubmit(e)}
          >
            {InvalidSecret ? "Invalid Secret Provided": fetching && !active
              ? "Generating Key..."
              : active
              ? "User Authenticated"
              : "Generate Key"}
          </button>
        </form>
      </div>
      
    </div>
  );
}
export default Home;
