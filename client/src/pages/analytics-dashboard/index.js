import { useEffect, useState } from "react";
import axios from "axios";

import Link from "next/link";

function AnalyticsDashboard() {
  const [associationCount, setAssociationCount] = useState(0);
  const [sourceSchemaCount, setSourceSchemaCount] = useState(0);
  const [targetSchemaCount, setTargetSchemaCount] = useState(0);
  const [sourceSchema, setSourceSchema] = useState([]);
  const [targetSchema, setTargetSchema] = useState([]);

  const handleGetAnalyticsData = async () => {
    try {
      console.log("Fetching analytics data..." + process.env.NEXT_PUBLIC_BACKEND_API_KEY);
      const res = await axios.get(`${process.env.NEXT_PUBLIC_BACKEND_API_KEY}/analytics`);
      console.log(res.data.data);

      setAssociationCount(res.data.data.associationCount);
      setSourceSchemaCount(res.data.data.sourceSchemaCount);
      setTargetSchemaCount(res.data.data.targetSchemaCount);
      setSourceSchema(res.data.data.sourceSchema);
      setTargetSchema(res.data.data.targetSchema);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    handleGetAnalyticsData();
  }, []); 
  return (
    <div className="dashboard">
      
      <div className="dashboard-box">
        <h1>Analytics Dashboard</h1>
        <Link href="/">
          <button className="home-button">Home</button>
        </Link>

        <div className="overview-container">
          <div className="overview-box">
            <h3>Association Count</h3>
            <p>{associationCount}</p>
          </div>
          <div className="overview-box">
            <h3>Source Schema Count</h3>
            <p>{sourceSchemaCount}</p>
          </div>
          <div className="overview-box">
            <h3>Target Schema Count</h3>
            <p>{targetSchemaCount}</p>
          </div>
        </div>

        <div className="schema-data">
          <h3>Source Schema Data</h3>
          <pre>{JSON.stringify(sourceSchema, null, 2)}</pre>
          
          <h3>Target Schema Data</h3>
          <pre>{JSON.stringify(targetSchema, null, 2)}</pre>
        </div>

        <button className="fetch-button" onClick={handleGetAnalyticsData}>
          Fetch Analytics Data
        </button>
      </div>
    </div>
  );
}

export default AnalyticsDashboard;
