package types

type ExtBool byte

const (
	ExtBoolNull  ExtBool = 0
	ExtBoolTrue  ExtBool = 1
	ExtBoolFalse ExtBool = 2
)

func (b ExtBool) IsNull() bool {
	if b == ExtBoolNull {
		return true
	}
	return false
}

func (b ExtBool) IsTrue() bool {
	if b == ExtBoolTrue {
		return true
	}
	return false
}

func (b ExtBool) IsFalse() bool {
	if b == ExtBoolFalse {
		return true
	}
	return false
}

func (b ExtBool) ToBool() bool {
	if b == ExtBoolTrue {
		return true
	}
	return false
}
