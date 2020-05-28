package labelselector

// LabelSelector is a set of label requirements
type LabelSelector struct {
	Requirements []Requirement
}

type Requirement struct {
	Key       string
	Value     string
	Values    []string
	Operation Operation
}

type Operation string

const (
	OperationIn        Operation = "in"
	OperationNotIn     Operation = "notIn"
	OperationExists    Operation = "exists"
	OperationNotExists Operation = "notExist"
	OperationEquals    Operation = "equals"
	OperationNotEquals Operation = "notEquals"
)
