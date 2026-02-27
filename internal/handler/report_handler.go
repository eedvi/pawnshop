package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
)

// ReportHandler handles report endpoints
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetDashboard retrieves dashboard statistics
func (h *ReportHandler) GetDashboard(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)

	stats, err := h.reportService.GetDashboardStats(c.Context(), int64(branchID))
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, stats)
}

// GetLoanReport retrieves loan report
func (h *ReportHandler) GetLoanReport(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	report, err := h.reportService.GetLoanReport(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, report)
}

// GetPaymentReport retrieves payment report
func (h *ReportHandler) GetPaymentReport(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	report, err := h.reportService.GetPaymentReport(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, report)
}

// GetSalesReport retrieves sales report
func (h *ReportHandler) GetSalesReport(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	report, err := h.reportService.GetSalesReport(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, report)
}

// GetOverdueReport retrieves overdue loans report
func (h *ReportHandler) GetOverdueReport(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)

	report, err := h.reportService.GetOverdueReport(c.Context(), int64(branchID))
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, report)
}

// ExportDailyReport exports daily report as PDF
func (h *ReportHandler) ExportDailyReport(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateStr := c.Query("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return response.BadRequest(c, "Invalid date format")
	}

	pdfData, err := h.reportService.GenerateDailyReportPDF(c.Context(), int64(branchID), date)
	if err != nil {
		return response.InternalError(c, "Failed to generate report")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=daily_report_"+dateStr+".pdf")
	return c.Send(pdfData)
}

// ExportLoanContract exports loan contract as PDF
func (h *ReportHandler) ExportLoanContract(c *fiber.Ctx) error {
	loanID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID format")
	}

	pdfData, err := h.reportService.GenerateLoanContractPDF(c.Context(), loanID)
	if err != nil {
		return response.InternalError(c, "Failed to generate contract")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=loan_contract_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// ExportPaymentReceipt exports payment receipt as PDF
func (h *ReportHandler) ExportPaymentReceipt(c *fiber.Ctx) error {
	paymentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	pdfData, err := h.reportService.GeneratePaymentReceiptPDF(c.Context(), paymentID)
	if err != nil {
		return response.InternalError(c, "Failed to generate receipt")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=payment_receipt_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// ExportSaleReceipt exports sale receipt as PDF
func (h *ReportHandler) ExportSaleReceipt(c *fiber.Ctx) error {
	saleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid sale ID format")
	}

	pdfData, err := h.reportService.GenerateSaleReceiptPDF(c.Context(), saleID)
	if err != nil {
		return response.InternalError(c, "Failed to generate receipt")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=sale_receipt_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// ExportLoanReportPDF exports loan report as PDF
func (h *ReportHandler) ExportLoanReportPDF(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	pdfData, err := h.reportService.ExportLoanReportPDF(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, "Failed to generate loan report PDF")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=reporte_prestamos_"+dateFrom+"_"+dateTo+".pdf")
	return c.Send(pdfData)
}

// ExportPaymentReportPDF exports payment report as PDF
func (h *ReportHandler) ExportPaymentReportPDF(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	pdfData, err := h.reportService.ExportPaymentReportPDF(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, "Failed to generate payment report PDF")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=reporte_pagos_"+dateFrom+"_"+dateTo+".pdf")
	return c.Send(pdfData)
}

// ExportSalesReportPDF exports sales report as PDF
func (h *ReportHandler) ExportSalesReportPDF(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	pdfData, err := h.reportService.ExportSalesReportPDF(c.Context(), int64(branchID), dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, "Failed to generate sales report PDF")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=reporte_ventas_"+dateFrom+"_"+dateTo+".pdf")
	return c.Send(pdfData)
}

// ExportOverdueReportPDF exports overdue report as PDF
func (h *ReportHandler) ExportOverdueReportPDF(c *fiber.Ctx) error {
	branchID := c.QueryInt("branch_id", 0)

	pdfData, err := h.reportService.ExportOverdueReportPDF(c.Context(), int64(branchID))
	if err != nil {
		return response.InternalError(c, "Failed to generate overdue report PDF")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=reporte_vencidos_"+time.Now().Format("2006-01-02")+".pdf")
	return c.Send(pdfData)
}

// ============================================================================
// Thermal Ticket Endpoints
// ============================================================================

// ExportLoanTicket exports a loan thermal ticket as PDF
func (h *ReportHandler) ExportLoanTicket(c *fiber.Ctx) error {
	loanID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID format")
	}

	paperSize := c.Query("paper_size", "80mm") // Default to 80mm

	pdfData, err := h.reportService.GenerateLoanTicket(c.Context(), loanID, paperSize)
	if err != nil {
		return response.InternalError(c, "Failed to generate loan ticket")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=ticket_prestamo_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// ExportPaymentTicket exports a payment thermal ticket as PDF
func (h *ReportHandler) ExportPaymentTicket(c *fiber.Ctx) error {
	paymentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	pdfData, err := h.reportService.GeneratePaymentTicket(c.Context(), paymentID)
	if err != nil {
		return response.InternalError(c, "Failed to generate payment ticket")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=ticket_pago_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// ExportSaleTicket exports a sale thermal ticket as PDF
func (h *ReportHandler) ExportSaleTicket(c *fiber.Ctx) error {
	saleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid sale ID format")
	}

	pdfData, err := h.reportService.GenerateSaleTicket(c.Context(), saleID)
	if err != nil {
		return response.InternalError(c, "Failed to generate sale ticket")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=ticket_venta_"+c.Params("id")+".pdf")
	return c.Send(pdfData)
}

// RegisterRoutes registers report routes
func (h *ReportHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	reports := app.Group("/reports")
	reports.Use(authMiddleware.Authenticate())

	// Dashboard
	reports.Get("/dashboard", authMiddleware.RequirePermission("reports.read"), h.GetDashboard)

	// Reports (JSON)
	reports.Get("/loans", authMiddleware.RequirePermission("reports.read"), h.GetLoanReport)
	reports.Get("/payments", authMiddleware.RequirePermission("reports.read"), h.GetPaymentReport)
	reports.Get("/sales", authMiddleware.RequirePermission("reports.read"), h.GetSalesReport)
	reports.Get("/overdue", authMiddleware.RequirePermission("reports.read"), h.GetOverdueReport)

	// Report exports (PDF)
	reports.Get("/loans/export", authMiddleware.RequirePermission("reports.export"), h.ExportLoanReportPDF)
	reports.Get("/payments/export", authMiddleware.RequirePermission("reports.export"), h.ExportPaymentReportPDF)
	reports.Get("/sales/export", authMiddleware.RequirePermission("reports.export"), h.ExportSalesReportPDF)
	reports.Get("/overdue/export", authMiddleware.RequirePermission("reports.export"), h.ExportOverdueReportPDF)

	// Individual document exports (PDF - full size)
	reports.Get("/export/daily", authMiddleware.RequirePermission("reports.export"), h.ExportDailyReport)
	reports.Get("/export/loan/:id/contract", authMiddleware.RequirePermission("reports.export"), h.ExportLoanContract)
	reports.Get("/export/payment/:id/receipt", authMiddleware.RequirePermission("reports.export"), h.ExportPaymentReceipt)
	reports.Get("/export/sale/:id/receipt", authMiddleware.RequirePermission("reports.export"), h.ExportSaleReceipt)

	// Thermal ticket exports (PDF - thermal paper size)
	reports.Get("/ticket/loan/:id", authMiddleware.RequirePermission("reports.export"), h.ExportLoanTicket)
	reports.Get("/ticket/payment/:id", authMiddleware.RequirePermission("reports.export"), h.ExportPaymentTicket)
	reports.Get("/ticket/sale/:id", authMiddleware.RequirePermission("reports.export"), h.ExportSaleTicket)
}
