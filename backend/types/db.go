package types

import "fmt"

func (ns *NullString) Scan(value any) error {
	if value == nil {
		*ns = NullString("")
		return nil
	}
	if str, ok := value.(string); ok {
		*ns = NullString(str)
		return nil
	}
	return fmt.Errorf("failed to scan NullString: %v", value)
}
