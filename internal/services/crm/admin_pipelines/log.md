# Log: admin_pipelines tool handler fix

_Дата: 2026-04-04_

---

## Итерация 1 — Перезапись app/gkit/tools/admin_pipelines.go

### Проблемы в старом файле

1. Импортировал `amomodels` и использовал `*amomodels.Pipeline` / `*amomodels.Status` напрямую — маппинг теперь в сервисе
2. Обращался к `input.Data["pipelines"]` и `input.Data["statuses"]` — поле `Data map[string]any` удалено из `AdminPipelinesInput`
3. `ListPipelines(ctx)` — старая сигнатура без `withStatuses bool`
4. `GetPipeline(ctx, input.PipelineID)` — нужно `(ctx, id, name, withStatuses)`
5. `ListStatuses(ctx, input.PipelineID)` — нужно `(ctx, pipelineID, pipelineName)`
6. `GetStatus(ctx, input.PipelineID, input.StatusID)` — нужно `(ctx, pipelineID, pipelineName, statusID, statusName)`
7. `CreatePipelines(ctx, []*amomodels.Pipeline)` — нужно `(ctx, []toolmodels.PipelineData)`
8. `UpdatePipeline(ctx, &p)` — нужно `(ctx, id, name, PipelineData)`
9. `DeletePipeline(ctx, id)` — нужно `(ctx, id, name)`
10. `CreateStatus/CreateStatuses/UpdateStatus/DeleteStatus` — аналогично устаревшие сигнатуры
11. Не передавал `input.PipelineName`, `input.StatusName`, `input.WithStatuses`
12. Описание tool упоминало `data.pipelines`/`data.statuses` — устаревший батч-формат

### Изменения

- Убран импорт `amomodels` — весь маппинг в сервисе
- Все вызовы методов обновлены под новые сигнатуры
- `create`: читает батч из `input.Items` через `json.Unmarshal` в `[]toolmodels.PipelineData`, или одиночный из `input.Pipeline`
- `update`: передаёт `(ctx, input.PipelineID, input.PipelineName, *input.Pipeline)`
- `delete`: передаёт `(ctx, input.PipelineID, input.PipelineName)`
- `list_statuses`: передаёт `(ctx, input.PipelineID, input.PipelineName)`
- `get_status`: передаёт `(ctx, input.PipelineID, input.PipelineName, input.StatusID, input.StatusName)`
- `create_status`: батч из `input.Items` → `[]toolmodels.StatusData` через `CreateStatuses`, одиночный из `input.Status` через `CreateStatus`
- `update_status`: передаёт `(ctx, input.PipelineID, input.PipelineName, input.StatusID, input.StatusName, *input.Status)`
- `delete_status`: передаёт `(ctx, input.PipelineID, input.PipelineName, input.StatusID, input.StatusName)`
- Обновлено описание tool — убраны `data.pipelines`/`data.statuses`, добавлен новый формат
- Убрана избыточная валидация `pipeline_id == 0` из handler — `resolvePipelineID` в сервисе вернёт понятную ошибку
