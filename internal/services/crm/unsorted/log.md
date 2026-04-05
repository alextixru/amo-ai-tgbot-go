# Лог изменений: сервис unsorted

---

## 2026-04-04 — Рефакторинг tool handler

### Файл: `app/gkit/tools/unsorted.go`

**Проблема:** Handler не компилировался из-за несовместимости со обновлёнными входными моделями и интерфейсом сервиса.

**Изменения:**

1. Удалён алиас `"search"` из switch — оставлен только `"list"`.

2. `case "create"`: убрана ручная сборка `[]*models.Unsorted` с обращением к старым полям `item.PipelineID (int)` и `int64(item.CreatedAt)`. Теперь `input.CreateData.Items` (`[]gkitmodels.UnsortedCreateItem`) передаётся напрямую в `CreateUnsorted(ctx, category, items)`. Маппинг полей и резолвинг `PipelineName` → ID выполняются в сервисном слое.

3. `case "accept"`: удалено обращение к несуществующим полям `input.AcceptParams.UserID` и `input.AcceptParams.StatusID`. Теперь `input.AcceptParams` (`*gkitmodels.UnsortedAcceptParams`) передаётся напрямую в `AcceptUnsorted(ctx, uid, params)`. Сервис выполняет резолвинг `UserName`, `PipelineName`, `StatusName` → ID.

4. `case "decline"`: удалось некорректное обращение к `input.AcceptParams.UserID`. Теперь `input.DeclineParams` (`*gkitmodels.UnsortedDeclineParams`) передаётся напрямую в `DeclineUnsorted(ctx, uid, params)`.

5. Удалён импорт `"github.com/alextixru/amocrm-sdk-go/core/models"` — он был нужен только для ручной сборки `[]*models.Unsorted`, которая перенесена в сервис.

**Результат:** Handler стал тонким маршрутизатором без знания о внутренних SDK-моделях. Пакет компилируется.
