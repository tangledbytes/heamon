import Typography from "@material-ui/core/Typography";
import { useEffect, useState } from "react";
import StatusCardGrid from "../../components/StatusGrid";
import StatusBanner from "../../components/StatusBanner";

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

function allFunctional(data: Array<any>): boolean {
  if (!Array.isArray(data)) return true;

  for (const el of data) {
    if (el.status !== "OK") return false;
  }

  return true;
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

  const isFunctional = allFunctional(data);

  return (
    <main>
      <StatusBanner isFunctional={isFunctional} style={{ marginBottom: "1rem" }} />
      <Typography variant="h4" align="center" gutterBottom>
        Current Status
      </Typography>

      <StatusCardGrid data={data} />
    </main>
  );
}

export default Home;
