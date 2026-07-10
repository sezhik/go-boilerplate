package example

type ListInput struct {
	Limit  int
	Offset int
}

type Item struct {
	ID   int64
	Name string
}

type ListOutput struct {
	Items  []Item
	Limit  int
	Offset int
	Total  int
}
