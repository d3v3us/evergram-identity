package account

const (
	AccountListCacheKey = "account.list"
	AccountByIdCacheKey = "account.by.id.%s"
)

var AccountCacheKeys = []string{AccountListCacheKey, AccountByIdCacheKey}
