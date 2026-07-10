package example

import modelexample "go-boilerplate/internal/model/example"

type Presenter struct{}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) List(output modelexample.ListOutput) ListResponse {
	items := make([]ItemResponse, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, ItemResponse{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return ListResponse{
		Items:  items,
		Limit:  output.Limit,
		Offset: output.Offset,
		Total:  output.Total,
	}
}
