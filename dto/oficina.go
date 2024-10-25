package dto

import (
	"time"

	"github.com/google/uuid"
)

type VoucherResponseDTO struct {
	ID          string `json:"id"`
	UsuarioID   string `json:"usuarioId"`
	NomeOficina string `json:"nomeOficina"`
	Descricao   string `json:"descricao"`
	Validado    bool   `json:"validado"`
	QRCode      []byte `json:"qrcode"`
	QRCodeURL   string `json:"qrcodeUrl"`
}

type ValidaVoucherDTO struct {
	IdVoucher string `json:"idVoucher"`
}

type VoucherValidation struct {
	IDVoucher string `json:"idVoucher"`
	Validado  bool   `json:"validado"`
}

type Teste struct {
	ID string `json:"id"`
}

type OficinaResponse struct {
	ID                  uuid.UUID `json:"id"`
	Nome                string    `json:"nome"`
	Descricao           string    `json:"descricao"`
	PrecoBoicoins       float64   `json:"precoBoicoins"`
	PrecoReal           float64   `json:"precoReal"`
	DataEvento          time.Time `json:"dataEvento"`
	LimiteParticipantes *int      `json:"limiteParticipantes"`
	ParticipantesAtual  int       `json:"participantesAtual"`
	LinkEndereco        string    `json:"linkEndereco"`
	Imagem              string    `json:"imagem"` // Alterado para string
}
