package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // other data types
	CSRFToken string                 // CSRF (Cross-Site Request Forgery) token: security token
	// msgs sent to client
	Flash   string
	Warning string
	Error   string
}
