import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import Icon from "../StatusCard/Icon";
import { Status, COLORGREEN, COLORRED } from "../../global/constants";

function Content({ isFunctional = true }: { isFunctional?: boolean }) {
  if (isFunctional) {
    return (
      <>
        <Icon
          status={Status.OK}
          style={{ marginRight: "1rem", color: "#fff" }}
        />
        <Typography style={{ color: "#fff" }} variant="h6" component="div">
          All Systems Operational
        </Typography>
      </>
    );
  }

  return (
    <>
      <Icon
        status={Status.FAIL}
        style={{ marginRight: "1rem", color: "#fff" }}
      />
      <Typography style={{ color: "#fff" }} variant="h6" component="div">
        Some Systems are Failing!
      </Typography>
    </>
  );
}

function StatusBanner({
  isFunctional = true,
  style = {},
}: {
  isFunctional?: boolean;
  style?: React.CSSProperties;
}) {
  return (
    <Paper
      elevation={3}
      style={{
        backgroundColor: isFunctional ? COLORGREEN : COLORRED,
        padding: "1rem",
        ...style
      }}
    >
      <div style={{ display: "flex", alignItems: "center" }}>
        <Content isFunctional={isFunctional} />
      </div>
    </Paper>
  );
}

export default StatusBanner;
