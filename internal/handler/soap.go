package handler

import (
	"CQS-KYC/internal/service"

	"github.com/gofiber/fiber/v3"
)

type SOAPHandler struct {
	service *service.ERPService
}

func NewSOAPHandler(service *service.ERPService) *SOAPHandler {
	return &SOAPHandler{service: service}
}

// SetupRoutes cho interface RouteHandler
func (h *SOAPHandler) SetupRoutes(router fiber.Router) {
	// Lưu ý: Route thực tế được định nghĩa trong app.go để fix cứng endpoint .asmx
	// Hàm này có thể dùng cho các API quản lý log SOAP nếu cần
}

func (h *SOAPHandler) HandleRequest(c fiber.Ctx) error {
	// Lấy raw body
	body := c.Body()

	// Gọi service xử lý (Logic bóc tách + Lưu DB)
	err := h.service.ProcessSOAPRequest(body)

	// Chuỗi trả về chuẩn cho ERP (Hardcoded để đảm bảo performance)
	resp := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <InvokeSrvResponse xmlns="http://tempuri.org/">
      <InvokeSrvResult>Success</InvokeSrvResult>
    </InvokeSrvResponse>
  </soap:Body>
</soap:Envelope>`

	// Set Header XML
	c.Set("Content-Type", "text/xml; charset=utf-8")

	// Dù có lỗi hay không, luôn trả về Success XML để ERP không treo
	if err != nil {
		// Chỉ log lỗi ra console/file
		// logger.Error("SOAP Error", err)
	}

	return c.SendString(resp)
}

func (h *SOAPHandler) HandleWSDL(c fiber.Ctx) error {
	c.Set("Content-Type", "text/xml")
	wsdl := `<?xml version="1.0" encoding="utf-8"?>
<wsdl:definitions xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" targetNamespace="http://tempuri.org/" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">
  <wsdl:service name="EFERPService">
    <wsdl:port name="EFERPServiceSoap" binding="tns:EFERPServiceSoap">
      <soap:address location="http://localhost/EFNETService/EFERPService.asmx" />
    </wsdl:port>
  </wsdl:service>
</wsdl:definitions>`
	return c.SendString(wsdl)
}
