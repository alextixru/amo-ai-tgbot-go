# SDK Method Coverage in Genkit Tools

## AccountService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| GetCurrent | ❌ | - | Нет в references.go |
| AvailableWith | ❌ | - | Нет в references.go |

## CallsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Add | ❌ | - | Нет tool |
| AddOne | ❌ | - | Нет tool |

## CatalogElementsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| GetOne | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Update | ❌ | - | Нет tool |

## CatalogsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| GetOne | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Update | ❌ | - | Нет tool |

## ChatTemplatesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Update | ❌ | - | Нет tool |
| Delete | ❌ | - | Нет tool |
| DeleteMany | ❌ | - | Нет tool |
| GetOne | ❌ | - | Нет tool |
| SendOnReview | ❌ | - | Нет tool |
| UpdateReviewStatus | ❌ | - | Нет tool |

## CompaniesService → crm_companies
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | search | |
| GetOne | ✅ | get | |
| Create | ✅ | create | |
| Update | ✅ | update | |
| Link | ❌ | - | Только через crm_manage_links? |
| Unlink | ❌ | - | Только через crm_manage_links? |

## ContactsService → crm_contacts
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | search | |
| GetOne | ✅ | get | |
| Create | ✅ | create | |
| Update | ✅ | update | |
| Link | ❌ | - | Только через crm_manage_links? |
| Unlink | ❌ | - | Только через crm_manage_links? |
| GetChats | ❌ | - | |
| LinkChats | ❌ | - | |

## CurrenciesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |

## CustomFieldGroupsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| GetOne | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Update | ❌ | - | Нет tool |
| Delete | ❌ | - | Нет tool |

## CustomFieldsService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | custom_fields | |
| GetOne | ✅ | custom_fields | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |

## CustomerBonusPointsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| EarnPoints | ❌ | - | Нет tool |
| RedeemPoints | ❌ | - | Нет tool |

## CustomerStatusesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| GetOne | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Update | ❌ | - | Нет tool |
| Delete | ❌ | - | Нет tool |

## CustomerTransactionsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | Нет tool |
| Create | ❌ | - | Нет tool |
| Delete | ❌ | - | Нет tool |

## CustomersService → crm_customers
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | search | |
| GetOne | ✅ | get | |
| Create | ✅ | create | |
| Update | ✅ | update | |
| Delete | ❌ | - | |
| Link | ❌ | - | |

## EntityFilesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Link | ❌ | - | |
| Unlink | ❌ | - | |

## EntitySubscriptionsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Subscribe | ❌ | - | |
| Unsubscribe | ❌ | - | |

## EventTypesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |

## EventsService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | events | |
| GetOne | ✅ | events | |

## FilesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Delete | ❌ | - | |
| UploadOne | ❌ | - | |

## LeadsService → crm_leads
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | search | |
| GetWithPagination | ❌ | - | |
| GetOne | ✅ | get | |
| Create | ✅ | create | |
| CreateOne | ❌ | - | Есть Create |
| Update | ✅ | update | |
| UpdateOne | ❌ | - | Есть Update |
| Delete | ❌ | - | |
| Link | ❌ | - | Через crm_manage_links? |
| Unlink | ❌ | - | Через crm_manage_links? |
| AddComplex | ❌ | - | |
| AddOneComplex | ❌ | - | |

## LinksService → crm_manage_links
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Link | ✅ | link | |
| Unlink | ✅ | unlink | |

## LossReasonsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |

## NotesService → crm_notes
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | list | |
| GetByParent | ✅ | list | |
| GetOne | ❌ | - | |
| Create | ✅ | create | |
| Update | ❌ | - | |

## PipelinesService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | pipelines | |
| GetOne | ✅ | pipelines | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |
| GetStatuses | ❌ | - | |
| GetStatus | ❌ | - | |
| CreateStatus | ❌ | - | |
| UpdateStatus | ❌ | - | |
| DeleteStatus | ❌ | - | |

## ProductsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |

## RolesService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | roles | |
| GetOne | ✅ | roles | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |

## SegmentsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Create | ❌ | - | |
| Delete | ❌ | - | |

## ShortLinksService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Create | ❌ | - | |
| Delete | ❌ | - | |

## SourcesService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Create | ❌ | - | |
| Update | ❌ | - | |
| Delete | ❌ | - | |

## TagsService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | tags | |
| Create | ❌ | - | |
| Delete | ❌ | - | |

## TalksService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Close | ❌ | - | |

## TasksService → crm_tasks
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | list | |
| GetOne | ❌ | - | |
| Create | ✅ | create | |
| Update | ❌ | - | |
| Complete | ✅ | complete | |

## UnsortedService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Create | ❌ | - | |
| Accept | ❌ | - | |
| Decline | ❌ | - | |
| Link | ❌ | - | |
| Summary | ❌ | - | |

## UsersService → crm_get_reference
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ✅ | users | |
| GetOne | ✅ | users | |
| Create | ❌ | - | |
| AddToGroup | ❌ | - | |
| GetRoles | ❌ | - | Есть RolesService.Get |
| GetRole | ❌ | - | Есть RolesService.GetOne |

## WebhooksService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| Subscribe | ❌ | - | |
| Unsubscribe | ❌ | - | |

## WebsiteButtonsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| CreateAsync | ❌ | - | |
| UpdateAsync | ❌ | - | |
| AddOnlineChatAsync | ❌ | - | |

## WidgetsService → ❓
| Метод | Реализован | Action | Комментарий |
|-------|------------|--------|-------------|
| Get | ❌ | - | |
| GetOne | ❌ | - | |
| Install | ❌ | - | |
| Uninstall | ❌ | - | |
