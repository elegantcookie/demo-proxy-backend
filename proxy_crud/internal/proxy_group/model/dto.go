package model

type CreateProxyGroupDTO struct {
	Name string `json:"name"`
}

func NewCreateProxyDTO(name string) (CreateProxyGroupDTO, error) {
	return CreateProxyGroupDTO{
		Name: name,
	}, nil
}
