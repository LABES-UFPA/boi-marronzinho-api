package usecase

import (
    "boi-marronzinho-api/adapter/repository"
    "boi-marronzinho-api/domain"
    "errors"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type DoacaoUseCase struct {
    doacaoRepo    repository.Repository[domain.Doacoes]
    itemDoacaoRepo repository.Repository[domain.ItemDoacao]
}

func NewDoacaoUseCase(doacaoRepo repository.Repository[domain.Doacoes], itemDoacaoRepo repository.Repository[domain.ItemDoacao]) *DoacaoUseCase {
    return &DoacaoUseCase{doacaoRepo: doacaoRepo, itemDoacaoRepo: itemDoacaoRepo}
}

func (duc *DoacaoUseCase) AdicionaDoacao(doacaoRequest *domain.Doacoes) (*domain.Doacoes, error) {
    itemDoacao, err := duc.itemDoacaoRepo.GetByID(doacaoRequest.ItemDoacaoID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("item de doação não encontrado")
        }
        return nil, err
    }

    boicoinsRecebidos, err := calculaBoicoins(itemDoacao.BoicoinsUnidade, doacaoRequest.Quantidade)
    if err != nil {
        return nil, err
    }

    doacao := &domain.Doacoes{
        ID:                uuid.New(),
        UsuarioID:         doacaoRequest.UsuarioID,
        ItemDoacaoID:      doacaoRequest.ItemDoacaoID,
        BoicoinsRecebidos: boicoinsRecebidos,
        DataDoacao:        time.Now(),
    }

    if _, err := duc.doacaoRepo.Create(doacao); err != nil {
        return nil, err
    }

    return doacao, nil
}

func (duc *DoacaoUseCase) CriarItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
    itemDoacao := &domain.ItemDoacao{
        ID:              uuid.New(),
        Descricao:       itemDoacaoRequest.Descricao,
        UnidadeMedida:   itemDoacaoRequest.UnidadeMedida,
        BoicoinsUnidade: itemDoacaoRequest.BoicoinsUnidade,
    }

    return duc.itemDoacaoRepo.Create(itemDoacao)
}

func (duc *DoacaoUseCase) AtualizaItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
    itemDoacao, err := duc.itemDoacaoRepo.GetByID(itemDoacaoRequest.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("item de doação não encontrado")
        }
        return nil, err
    }

    if itemDoacaoRequest.Descricao != "" {
        itemDoacao.Descricao = itemDoacaoRequest.Descricao
    }
    if itemDoacaoRequest.UnidadeMedida != "" {
        itemDoacao.UnidadeMedida = itemDoacaoRequest.UnidadeMedida
    }
    if itemDoacaoRequest.BoicoinsUnidade != 0 {
        itemDoacao.BoicoinsUnidade = itemDoacaoRequest.BoicoinsUnidade
    }

    return duc.itemDoacaoRepo.Update(itemDoacao)
}

func (duc *DoacaoUseCase) DeletarItemDoacao(id uuid.UUID) error {
    if err := duc.itemDoacaoRepo.Delete(id); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("item de doação não encontrado")
        }
        return err
    }
    return nil
}

func calculaBoicoins(valorUnidade float64, quantidade int64) (float64, error) {
    if valorUnidade < 0 || quantidade < 0 {
        return 0, errors.New("valorUnidade e quantidade devem ser maiores que zero")
    }
    return valorUnidade * float64(quantidade), nil
}
