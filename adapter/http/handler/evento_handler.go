package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	minioClient "boi-marronzinho-api/minio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventoHandler struct {
	EventoUseCase *usecase.EventoUseCase
}

func NewEventoHandler(euc *usecase.EventoUseCase) *EventoHandler {
	return &EventoHandler{EventoUseCase: euc}
}

func (eh *EventoHandler) CriaEvento(c *gin.Context) {
    var eventoDTO domain.Evento

	jsonData := c.PostForm("request")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'request' com os dados do evento é obrigatório."})
		return
	}
	
	if err := json.Unmarshal([]byte(jsonData), &eventoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao parsear os dados JSON do evento."})
		return
	}	

    dataEvento, err := time.Parse(time.RFC3339, eventoDTO.DataEvento.Format(time.RFC3339))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use o formato RFC3339."})
        return
    }
    eventoDTO.DataEvento = dataEvento

    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao receber o arquivo de imagem"})
        return
    }

    fileContent, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo de imagem"})
        return
    }
    defer fileContent.Close()

    imageFileName := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)

    imageUrl, err := minioClient.UploadFile(fileContent, imageFileName, "eventos")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer upload da imagem"})
        return
    }

    eventoDTO.ImagemUrl = imageUrl
    eventoDTO.ID = uuid.New()
    eventoDTO.CriadoEm = time.Now()

    evento, err := eh.EventoUseCase.CriaEvento(&eventoDTO)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, evento)
}




func (eh *EventoHandler) ListaEventos(c *gin.Context) {
	eventos, err := eh.EventoUseCase.ListaEventos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eventos)
}

func (eh *EventoHandler) GetEvento(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do evento é inválido"})
		return
	}

	evento, err := eh.EventoUseCase.GetEvento(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, evento)
}

func (eh *EventoHandler) UpdateEvento(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do evento é inválido"})
		return
	}

	var updateData domain.Evento
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	updatedEvento, err := eh.EventoUseCase.UpdateEvento(id, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEvento)
}

func (eh *EventoHandler) DeleteEvento(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do evento é inválido"})
		return
	}

	err = eh.EventoUseCase.DeleteEvento(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evento deletado com sucesso!"})
}
