package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-uploader/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-uploader/pkg/httpserver"
)

// @Summary Login
// @Description Login the customer by its CPF
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body entity.CustomerForm true "customer form"
// @Success 200 {object} entity.Token
// @Failure 404 "Customer not found"
// @Router /auth/login [post]
func LoginCustomerHandler(loginCustomerUseCase usecase.LoginCustomerUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customerForm entity.CustomerForm

		err := httpserver.DecodeJSONBody(w, r, &customerForm)

		if err != nil {
			log.Print("decoding customer form body", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		token, err := loginCustomerUseCase.Execute(r.Context(), customerForm.CPF)

		if err != nil {
			log.Print("login user", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, token)
	}
}

// @Summary Create new customer
// @Description Create new customer. This process is not required to make an order
// @Tags Customer
// @Accept json
// @Produce json
// @Param product body entity.Customer true "customer"
// @Success 200 {object} entity.CustomerResponse
// @Failure 400 "Customer has required fields"
// @Failure 409 "This Customer is already added"
// @Router /auth/signup [post]
func CreateUserHandler(createCustomer usecase.CreateCustomerUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer entity.Customer

		err := httpserver.DecodeJSONBody(w, r, &customer)

		if err != nil {
			log.Print("decoding customer body", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		response, err := createCustomer.Execute(r.Context(), customer)

		if err != nil {
			log.Print("create customer", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Update customer
// @Description Update customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Param product body entity.Customer true "customer"
// @Success 204
// @Failure 400 "Customer has required fields"
// @Failure 404 "Customer not found"
// @Router /api/customers/{id} [put]
func UpdateCustomerHandler(updateCustomer usecase.UpdateCustomerUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerIdStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update customer", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		customerId, err := strconv.Atoi(customerIdStr)

		if err != nil {
			log.Print("update customer", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		var customer entity.Customer
		err = httpserver.DecodeJSONBody(w, r, &customer)

		if err != nil {
			log.Print("decoding customer body for update", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		customer.ID = uint(customerId)
		err = updateCustomer.Execute(r.Context(), customer)

		if err != nil {
			log.Print("update customer", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}
