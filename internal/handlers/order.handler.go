package handlers

import (
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type HandlerOrder struct {
	repositories.IOrderRepository
}

func InitializeOrderHandler(r repositories.IOrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetOrder(ctx *gin.Context) {}

func (h *HandlerOrder) GetDetailSchedule(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("schedule_id"))
	result, err := h.RepositoryGetScheduleDetail(ID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Schedule not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Success", result, nil))
}

func (h *HandlerOrder) CreateTransaction(ctx *gin.Context) {
	var dataOrder models.OrderDetailModel
	if err := ctx.ShouldBind(&dataOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding data order", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&dataOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	tx, errTx := h.Begin()
	if errTx != nil {
		log.Println(errTx.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error begin tx", nil, nil))
		return
	}
	user_id, _ := helpers.GetPayload(ctx)
	Order_Id := uuid.NewString()
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
	// chargeReq := &coreapi.ChargeReq{
	// 	PaymentType: coreapi.PaymentTypeBankTransfer,
	// 	TransactionDetails: midtrans.TransactionDetails{
	// 		OrderID:  Order_Id,
	// 		GrossAmt: int64(dataOrder.Price_Amount),
	// 	},
	// 	BankTransfer: &coreapi.BankTransferDetails{
	// 		Bank: midtrans.Bank(dataOrder.Payment),
	// 	},
	// 	Items: &[]midtrans.ItemDetails{
	// 		{
	// 			ID:    "Ticket_Purchasing",
	// 			Price: int64(dataOrder.Price_Amount),
	// 			Qty:   1,
	// 		},
	// 	},
	// }
	// coreApiRes, err := coreapi.ChargeTransaction(chargeReq)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error charging transaction", nil, nil))
	// 	return
	// }
	var s = snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	reqSnap := snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  Order_Id,
			GrossAmt: int64(dataOrder.Price_Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}
	snapResp, err := s.CreateTransaction(&reqSnap)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error charging transaction", nil, nil))
		return
	}
	errOrder := h.RepositoryCreateOrder(tx, Order_Id, user_id, &dataOrder, snapResp.RedirectURL)
	if errOrder != nil {
		log.Println(errOrder.Error())
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	// updateSeats, errUpdate := h.RepositoryUpdateSeatSchedule(tx, &dataOrder)
	// if errUpdate != nil {
	// 	log.Println(errUpdate.Error())
	// 	tx.Rollback()
	// 	ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in updatign seats", nil, nil))
	// 	return
	// }
	// var dataExecuted int64 = 1
	// if updateSeats < dataExecuted {
	// 	tx.Rollback()
	// 	ctx.JSON(http.StatusNotFound, helpers.NewResponse("Schedule id not found", nil, nil))
	// 	return
	// }
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in comitting order", nil, nil))
		return
	}
	ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Successfully create order", snapResp.RedirectURL, nil))
}

func (h *HandlerOrder) SubmitPayment(ctx *gin.Context) {}

func (h *HandlerOrder) FailedPayment(ctx *gin.Context) {}
