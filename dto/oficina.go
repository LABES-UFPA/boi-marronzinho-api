package dto

type VoucherResponseDTO struct {
	ID          string `json:"id"`
	UsuarioID   string `json:"usuarioId"`
	NomeOficina string `json:"nomeOficina"`
	Descricao   string `json:"descricao"`
	QRCode      []byte `json:"qrcode"`
}
