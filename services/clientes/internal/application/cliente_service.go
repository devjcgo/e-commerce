package application

import (
	"context"
	"ecommerce/clientes/internal/domain"
)

// ClienteService é a implementação dos nossos casos de uso de cliente.
type ClienteService struct {
	repo domain.ClienteRepository
}

// NewClienteService é o construtor do nosso serviço de aplicação.
func NewClienteService(repo domain.ClienteRepository) *ClienteService {
	return &ClienteService{
		repo: repo,
	}
}

// EnderecoInput é um DTO para os dados de endereço vindos da requisição.
type EnderecoInput struct {
	Rua    string `json:"rua"`
	Cidade string `json:"cidade"`
	Estado string `json:"estado"`
	CEP    string `json:"cep"`
}

// ClienteInput é o DTO principal para a criação de um novo cliente.
type ClienteInput struct {
	Nome      string          `json:"nome"`
	Email     string          `json:"email"`
	Enderecos []EnderecoInput `json:"enderecos"`
}

// CriarCliente é o caso de uso para criar um novo cliente.
// Ele orquestra a conversão de DTOs para o domínio e a persistência.
func (s *ClienteService) CriarCliente(ctx context.Context, input ClienteInput) (*domain.Cliente, error) {
	// 1. Converte os DTOs de EnderecoInput para o tipo do domínio.
	var enderecosDominio []*domain.Endereco
	for _, endInput := range input.Enderecos {
		enderecosDominio = append(enderecosDominio, &domain.Endereco{
			Rua:    endInput.Rua,
			Cidade: endInput.Cidade,
			Estado: endInput.Estado,
			CEP:    endInput.CEP,
		})
	}

	// 2. Cria a entidade principal do domínio.
	// Em um cenário mais complexo, aqui poderíamos chamar um construtor
	// como domain.NewCliente() que validaria as regras de negócio.
	novoCliente := &domain.Cliente{
		Nome:      input.Nome,
		Email:     input.Email,
		Enderecos: enderecosDominio,
	}

	// 3. Chama o repositório para salvar o novo cliente no banco de dados.
	if err := s.repo.Save(ctx, novoCliente); err != nil {
		return nil, err
	}

	// 4. Retorna o cliente criado (agora com ID e datas preenchidas pelo repositório).
	return novoCliente, nil
}

// ListarClientes é o caso de uso para buscar todos os clientes.
func (s *ClienteService) ListarClientes(ctx context.Context) ([]*domain.Cliente, error) {
	// A lógica aqui é simples: apenas repassamos a chamada para a camada de repositório.
	return s.repo.FindAll(ctx)
}
