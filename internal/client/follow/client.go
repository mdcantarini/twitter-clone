package follow

type Client interface {
	FetchFollowerIds(userID uint) ([]uint, error)
}
