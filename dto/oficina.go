package dto

type VoucherResponseDTO struct {
	ID          string `json:"id"`
	UsuarioID   string `json:"usuarioId"`
	NomeOficina string `json:"nomeOficina"`
	Descricao   string `json:"descricao"`
	QRCode      []byte `json:"qrcode"`
}


type ValidaVoucherDTO struct {
	IdVoucher string `json:"idVoucher"`
}


type VoucherValidation struct {
	IDVoucher string `json:"idVoucher"`
	Validado bool   `json:"validado"`
}