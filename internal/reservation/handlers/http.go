package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/ports"
)

var (
	ErrInvalidJSON = errors.New("invalid json")
)

type ReservationHandler struct {
	service  ports.ReservationService
	validate *validator.Validate
}

func NewReservationHandler(service ports.ReservationService, validate *validator.Validate) *ReservationHandler {
	return &ReservationHandler{service: service, validate: validate}
}

// Reserve of ReservationHandler
// @Tags reservation
// @Description Creates a reservation for given items if storehouse have required amount
// @Accept json
// @Produce json
// @Param input body domain.ReserveRequest true "destination location and items to reserve"
// @Success 200 {object} ports.ReservationResponseDTO
// @Failure 400 {object} string
// @Router /reserve [post]
func (handler *ReservationHandler) Reserve(c *gin.Context) {
	var dto domain.ReserveRequest

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(ErrInvalidJSON)
		return
	}

	err = handler.validate.Struct(dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reservationResponse, err := handler.service.Reserve(dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reservationResponse)
}

// Release of ReservationHandler
// @Tags reservation
// @Description Releases items for given reservation. If there is no items left, deleted the reservation
// @Accept json
// @Produce json
// @Param input body ports.ReleaseRequestDTO true "reservation ID and items to release"
// @Success 200 {object} ports.ReservationResponseDTO
// @Failure 400 {object} string
// @Router /release [post]
func (handler *ReservationHandler) Release(c *gin.Context) {
	var dto ports.ReleaseRequestDTO

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(ErrInvalidJSON)
		return
	}

	err = handler.validate.Struct(dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reservationResponse, err := handler.service.Release(dto.ReservationID, dto.ItemsToRelease)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reservationResponse)
}

// GetUnreserved of ReservationHandler
// @Tags reservation
// @Description Returns all unreserved items for given storehouse
// @Accept json
// @Produce json
// @Param storehouse-id query string true "storehouse ID"
// @Success 200 {object} ports.GetUnreservedResponseDTO
// @Failure 400 {object} string
// @Router /get-unreserved-items [get]
func (handler *ReservationHandler) GetUnreserved(c *gin.Context) {
	var dto ports.GetUnreservedRequestDTO

	err := c.ShouldBindQuery(&dto)
	if err != nil {
		_ = c.Error(ErrInvalidJSON)
		return
	}

	err = handler.validate.Struct(dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	unreservedItems, err := handler.service.GetUnreserved(dto.StorehouseID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ports.GetUnreservedResponseDTO{StorehouseID: dto.StorehouseID, Items: unreservedItems})
}
