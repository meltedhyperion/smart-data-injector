import { useEffect, useState } from "react";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import axios from "axios";
import Link from "next/link";
import { useRouter } from "next/router";
const status = {
  notStarted: "Not Completed",
  started: '<b style="color: yellow;">Started</b>',
  completed: '<b style="color: green;">Completed</b>',
};

function Injector() {
  const [file_uploaded, setFileUploaded] = useState(status.notStarted);
  const [s3_trigger_started, setS3TriggerStarted] = useState(status.notStarted);
  const [parsed_to_json, setParsedToJson] = useState(status.notStarted);
  const [source_schema_metadata, setSourceSchemaMetadata] = useState(status.notStarted);
  const [target_schema_metadata, setTargetSchemaMetadata] = useState(status.notStarted);
  const [association_mapping_generated, setAssociationMappingGenerated] = useState(status.notStarted);
  const [data_injection_completed, setDataInjectionCompleted] = useState(status.notStarted);

  const [file, setFile] = useState();
  const [apiKey, setApiKey] = useState("");
  const [uploadEnabled, setUploadEnabled] = useState(false);
  const [logId, setLogId] = useState(null);
  const [startFetchingLogs, setStartFetchingLogs] = useState(false);

  const router = useRouter();

  const renderStatus = (htmlString) => {
    return <span dangerouslySetInnerHTML={{ __html: htmlString }} />;
  };

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const handleApiKeyChange = (event) => {
    setApiKey(event.target.value);
  };

  const startSimulation = async () => {
    if (!apiKey || !file) {
      errorNotify("Please provide both API key and file.");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await axios.post(`${process.env.NEXT_PUBLIC_BACKEND_API_KEY}/upload`, formData, {
        headers: {
          Authorization: apiKey,
          "Content-Type": "multipart/form-data",
        },
      });
      const upload_response = response.data;
      setLogId(upload_response.data.logId);
      setStartFetchingLogs(true);
      notify("🚀 Starting Jobs...");
    } catch (error) {
      console.error("Error:", error.response);
      errorNotify(error.response?.data?.message || "File upload failed.");
    }
  };

  const fetchLogStatus = async () => {
    try {
      if (!logId) return;
      const response = await axios.get(`${process.env.NEXT_PUBLIC_BACKEND_API_KEY}/logs/${logId}`, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (response.data.status === 200) {
        const data = response.data.data;

        setFileUploaded(data.file_upload.status ? status.completed : status.notStarted);
        setS3TriggerStarted(data.s3_trigger_completed.status ? status.completed : status.notStarted);
        setParsedToJson(data.parsed_to_json.status ? status.completed : status.notStarted);
        setSourceSchemaMetadata(data.source_schema_metadata.status ? status.completed : status.notStarted);
        setTargetSchemaMetadata(data.target_schema_metadata.status ? status.completed : status.notStarted);
        setAssociationMappingGenerated(data.association_mapping_generated.status ? status.completed : status.notStarted);
        setDataInjectionCompleted(data.data_injection_completed.status ? status.completed : status.notStarted);
      } else {
        errorNotify("Failed to fetch log status.");
      }
    } catch (error) {
      console.error("Error fetching log status:", error);
    }
  };

  useEffect(() => {
    if (!startFetchingLogs) return;

    const intervalId = setInterval(() => {
      fetchLogStatus();
    }, 2000);

    return () => clearInterval(intervalId);
  }, [startFetchingLogs]);

  const notify = (message) =>
    toast(message, {
      position: "top-right",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "dark",
    });

  const errorNotify = (message) =>
    toast.error(`❌ ${message}`, {
      position: "top-right",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "dark",
    });

  return (
    <div className="master">
      <ToastContainer />
      <Link href="/">
        <button>Home</button>
      </Link>
      <div className="login-box">
        <label htmlFor="apiKey" id="text">API Key:</label>
        <input type="text" id="apiKey" onChange={handleApiKeyChange} value={apiKey} />
        <label htmlFor="file" id="text">Upload File:</label>
        <input type="file" id="file" onChange={handleFileChange} accept=".csv, .json" />
        <button onClick={startSimulation} disabled={!apiKey || !file}>Upload and Start ETL</button>
        <ol className="orderList">
          <li className="listItem"><b>S3 Upload Status:</b> {renderStatus(file_uploaded)}</li>
          <li className="listItem"><b>Trigger Job:</b> {renderStatus(s3_trigger_started)}</li>
          <li className="listItem"><b>Upload to Json Status:</b> {renderStatus(parsed_to_json)}</li>
          <li className="listItem"><b>Source Schema Metadata extraction:</b> {renderStatus(source_schema_metadata)}</li>
          <li className="listItem"><b>Target Schema Metadata extraction:</b> {renderStatus(target_schema_metadata)}</li>
          <li className="listItem"><b>Association Mapping Generation:</b> {renderStatus(association_mapping_generated)}</li>
          <li className="listItem"><b>Data Injection to DB:</b> {renderStatus(data_injection_completed)}</li>
        </ol>
        {data_injection_completed === status.completed && (
          <button
            onClick={() => router.push("/analytics-dashboard")} 
          >
            Proceed to Analytics
          </button>
        )}
      </div>
    </div>
  );
}

export default Injector;
