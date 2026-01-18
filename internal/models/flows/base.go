package flows

// FlowMode режим выполнения Flow
type FlowMode string

const (
	// ModeDirect прямой вызов SDK операции
	ModeDirect FlowMode = "direct"
	// ModeComplex передача Sub-Agent для сложной логики
	ModeComplex FlowMode = "complex"
)

// BaseFlowInput базовые поля для всех Flow
type BaseFlowInput struct {
	Mode   FlowMode `json:"mode" jsonschema_description:"Режим: direct (простая операция) или complex (сложная логика)"`
	Intent string   `json:"intent,omitempty" jsonschema_description:"Описание намерения (для mode=complex)"`
}
