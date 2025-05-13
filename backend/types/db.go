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

func (at *AccessType) Scan(value any) error {
	if value == nil {
		*at = AccessType(0)
		return nil
	}
	if num, ok := value.(int); ok {
		*at = AccessType(num)
		return nil
	}

	return fmt.Errorf("failed to scan AccessType: %v", value)
}

func (ar *AccessRole) Scan(value any) error {
	if value == nil {
		*ar = AccessRole(0)
		return nil
	}
	if num, ok := value.(int); ok {
		*ar = AccessRole(num)
		return nil
	}

	return fmt.Errorf("failed to scan AccessRule: %v", value)
}

func (ot *ObjectType) Scan(value any) error {
	if value == nil {
		*ot = ObjectType("")
		return nil
	}
	if str, ok := value.(string); ok {
		*ot = ObjectType(str)
		return nil
	}
	return fmt.Errorf("failed to scan NullString: %v", value)
}
