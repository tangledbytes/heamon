import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import Icon from "../StatusCard/Icon";
import { Status, COLORGREEN, COLORRED, COLORYELLOW, COLORGRAY } from "../../global/constants";

function colorForStatus(status : Status): string {
  if (status === Status.OK) return COLORGREEN;
  if (status === Status.FAIL) return COLORRED;
  if (status === Status.DEGRADED) return COLORYELLOW;

  return COLORGRAY
}

function Content({ status = Status.UNKNOWN }: { status?: Status }) {
  let content: Record<Status, string> = {
    OK: "All Systems Operational",
    FAIL: "Some Systems are Failing!",
    DEGRADED: "Some Systems are in Degraded State",
    UNKNOWN: "All Systems Operational"
  }

  return (
    <>
      <Icon
        status={status}
        style={{ marginRight: "1rem", color: "#fff" }}
      />
      <Typography style={{ color: "#fff" }} variant="h6" component="div">
        {content[status]}
      </Typography>
    </>
  );
}

function StatusBanner({
  status = Status.UNKNOWN,
  style = {},
}: {
  status?: Status;
  style?: React.CSSProperties;
}) {

  return (
    <Paper
      elevation={3}
      style={{
        backgroundColor: colorForStatus(status),
        padding: "1rem",
        ...style
      }}
    >
      <div style={{ display: "flex", alignItems: "center" }}>
        <Content status={status} />
      </div>
    </Paper>
  );
}

export default StatusBanner;
