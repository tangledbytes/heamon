import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import Typography from "@material-ui/core/Typography";
import Icon from "./Icon";
import { Status } from "../../global/constants"

export interface StatusCardProps {
  name: string;
  status: Status;
}

function StatusCard({ name, status }: StatusCardProps) {
  return (
    <Card style={{ height: "100%" }}>
      <CardContent>
          <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", height: "100%" }} >
            <Typography variant="h6" component="div">{name}</Typography>
            <div style={{ margin: "0 0 0 auto", display: "flex" }}>
                <Icon status={status}/>
            </div>
          </div>
      </CardContent>
    </Card>
  );
}

export default StatusCard;
