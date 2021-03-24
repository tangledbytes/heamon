package monitor

import (
	"fmt"

	"github.com/utkarsh-pro/heamon/models"
	"github.com/utkarsh-pro/heamon/pkg/eventbus"
)

type Status struct {
	*models.Status

	subscriberIds map[string]int64
}

// Refresh will publish the refresh event and will refresh the list
func (st *Status) Refresh(svcs []models.Service) {
	// Unsubscribe to the events
	for _, svc := range st.Report {
		topic := fmt.Sprintf("%s.%s", StatusUpdate, svc.Name)
		eventbus.Bus.Unsubscribe(topic, st.subscriberIds[topic])
	}

	st.Status = models.NewStatus(svcs)
}
