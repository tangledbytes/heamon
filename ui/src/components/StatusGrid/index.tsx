import Grid from "@material-ui/core/Grid";
import StatusCard, { StatusCardProps } from "../StatusCard";

export interface StatusCardGridProps {
  data: Array<StatusCardProps>;
}

function StatusCardGrid({ data = [] }: StatusCardGridProps) {
  return (
    <Grid container spacing={1}>
      {data.map((d, i) => (
        <Grid key={`status-card-${i}`} item md={4} xs={6}>
          <StatusCard {...d} />
        </Grid>
      ))}
    </Grid>
  );
}

export default StatusCardGrid;
