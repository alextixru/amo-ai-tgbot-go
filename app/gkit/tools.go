package gkit

// CRM Tools для Genkit
//
// TODO: Реализовать инструменты для работы с amoCRM
//
// Пример правильного определения (по документации Genkit):
//
//	tool := genkit.DefineTool(g, "getLeads", "Получить список сделок",
//	    func(ctx *ai.ToolContext, input struct {
//	        Query string `jsonschema_description:"Поисковый запрос"`
//	    }) (string, error) {
//	        // Реализация
//	        return "...", nil
//	    },
//	)
//
// Важно: первый параметр — *ai.ToolContext, не context.Context!
