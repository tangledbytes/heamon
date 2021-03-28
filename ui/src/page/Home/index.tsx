import Typography from "@material-ui/core/Typography";
import { useEffect, useState } from "react";
import StatusCardGrid from "../../components/StatusGrid";
import StatusBanner from "../../components/StatusBanner";
import { Status } from "../../global/constants";

async function GetStatus() {
  try {
    const res = await fetch("/api/v1/status");
    const json = await res.json();

    return json;
  } catch (err) {
    console.error(err);
    return {};
  }
}

function constructStatus(data: Array<any>): Status {
  if (!Array.isArray(data)) return Status.UNKNOWN;

  for (const el of data) {
    if (el.status === "FAIL") return Status.FAIL;
  }

  for (const el of data) {
    if (el.status === "DEGRADED") return Status.DEGRADED;
  }

  return Status.OK
}

function Home() {
  const [data, setData] = useState([]);

  useEffect(() => {
    GetStatus().then((res) => {
      setData(
        res.report?.map((r: any) => ({
          name: r.service.name,
          status: r.health_status,
        }))
      );
    });
  }, []);

  return (
    <main>
      <StatusBanner status={constructStatus(data)} style={{ marginBottom: "1rem" }} />
      <Typography variant="h4" align="center" gutterBottom>
        Current Status
      </Typography>

      <StatusCardGrid data={data} />
    </main>
  );
}

export default Home;
