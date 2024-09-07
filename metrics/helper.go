package metrics

import (
	"fmt"
	"time"
)

type APIRequest struct {
	Method           string
	Path             string
	StatusCode       int
	Success          bool
	Duration         time.Duration
	BusinessCategory string
}

func (r APIRequest) ToLabels() []Label {
	return []Label{
		{Name: "Method", Value: r.Method},
		{Name: "Path", Value: r.Path},
		{Name: "StatusCode", Value: fmt.Sprintf("%d", r.StatusCode)},
		{Name: "Success", Value: fmt.Sprintf("%t", r.Success)},
		{Name: "BusinessCategory", Value: r.BusinessCategory},
	}
}

func RecordAPIRequest(req APIRequest) {
	labels := req.ToLabels()
	EmitCounter("api.requests", 1, labels...)
	EmitTimer("api.duration", req.Duration, labels...)
}

type DBOperation struct {
	Database         string
	Table            string
	OperationType    string
	BusinessCategory string
	Success          bool
	Duration         time.Duration
}

func (op DBOperation) ToLabels() []Label {
	return []Label{
		{Name: "Database", Value: op.Database},
		{Name: "Table", Value: op.Table},
		{Name: "OperationType", Value: op.OperationType},
		{Name: "BusinessCategory", Value: op.BusinessCategory},
		{Name: "Success", Value: fmt.Sprintf("%t", op.Success)},
	}
}

func RecordDBOperation(op DBOperation) {
	labels := op.ToLabels()
	EmitCounter("db.operations", 1, labels...)
	EmitTimer("db.duration", op.Duration, labels...)
}

type LogicOperation struct {
	PrimaryCategory   string
	SecondaryCategory string
	TertiaryCategory  string
	Success           bool
	Duration          time.Duration
}

func (op LogicOperation) ToLabels() []Label {
	return []Label{
		{Name: "PrimaryCategory", Value: op.PrimaryCategory},
		{Name: "SecondaryCategory", Value: op.SecondaryCategory},
		{Name: "TertiaryCategory", Value: op.TertiaryCategory},
		{Name: "Success", Value: fmt.Sprintf("%t", op.Success)},
	}
}

func RecordLogicOperation(op LogicOperation) {
	labels := op.ToLabels()
	EmitCounter("logic.operations", 1, labels...)
	EmitTimer("logic.duration", op.Duration, labels...)
}
