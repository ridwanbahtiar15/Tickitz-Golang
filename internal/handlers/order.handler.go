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
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type HandlerOrder struct {
	repositories.IOrderRepository
}

func InitializeOrderHandler(r repositories.IOrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetOrder(ctx *gin.Context) {
	id, _ := helpers.GetPayload(ctx)
	page, _ := ctx.GetQuery("page")
	result, err := h.RepositoryGetOrderByID(id, page)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Order not found", nil, nil))
		return
	}
	data, errCount := h.RepositoryCountAllOrder(id)
	if errCount != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	currentPage := 1
	if page != "" {
		pageData, _ := strconv.Atoi(page)
		currentPage = pageData
	}
	meta := helpers.GetPagination(ctx, data, currentPage, 4)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Success", result, &meta))
}

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

	success, fail := CoreAPIMidtrans(Order_Id, dataOrder.Price_Amount)
	if fail != nil {
		log.Println(fail.Error())
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Midtrans Error", nil, nil))
		return
	}

	// success, fail := SnapMidtrans(Order_Id, dataOrder.Price_Amount)
	// if fail != nil {
	// 	log.Println(fail.Error())
	// 	ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Midtrans Error", nil, nil))
	// 	return
	// }

	virtualAccount := success.VaNumbers[0].VANumber
	errOrder := h.RepositoryCreateOrder(tx, Order_Id, user_id, &dataOrder, virtualAccount)
	if errOrder != nil {
		log.Println(errOrder.Error())
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in comitting order", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully create order", success.VaNumbers, nil))
}

func (h *HandlerOrder) SubmitPayment(ctx *gin.Context) {
	var webhookData map[string]interface{}
	if err := ctx.BindJSON(&webhookData); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding data order", nil, nil))
		return
	}
	// status := webhookData["transaction_status"].(string)
	orderID := webhookData["order_id"].(string)
	// if status != "capture" {
	// 	ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Payment error", nil, nil))
	// 	return
	// }
	result, err := h.RepositoryOrderSuccess(orderID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var successEdit int64 = 1
	if result < successEdit {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Success Payment", nil, nil))
}

func (h *HandlerOrder) FailedPayment(ctx *gin.Context) {
	var webhookData map[string]interface{}
	if err := ctx.BindJSON(&webhookData); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding data order", nil, nil))
		return
	}
	// status := webhookData["transaction_status"].(string)
	orderID := webhookData["order_id"].(string)
	// if status != "capture" {
	// 	ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Payment error", nil, nil))
	// 	return
	// }
	result, err := h.RepositoryOrderFailed(orderID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var successEdit int64 = 1
	if result < successEdit {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Payment Cancelled", nil, nil))
}

func SnapMidtrans(Order_Id string, price int) (*snap.Response, error) {
	var s = snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	reqSnap := snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  Order_Id,
			GrossAmt: int64(price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}
	snapResp, err := s.CreateTransaction(&reqSnap)
	// data := snapResp.
	if err != nil {
		return nil, err
	}
	return snapResp, nil
}

func CoreAPIMidtrans(Order_Id string, price int) (*coreapi.ChargeResponse, error) {
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  Order_Id,
			GrossAmt: int64(price),
		},
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.Bank("bca"),
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "Ticket_Purchasing",
				Name:  "Ticket",
				Price: int64(price),
				Qty:   1,
			},
		},
	}
	coreApiRes, err := coreapi.ChargeTransaction(chargeReq)
	// data := coreApiRes.
	if err != nil {
		return nil, err
	}
	return coreApiRes, nil
}
