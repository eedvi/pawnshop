package service

import (
	"context"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/pdf"
	"pawnshop/internal/repository"
)

// ReportService handles report generation
type ReportService struct {
	loanRepo     repository.LoanRepository
	paymentRepo  repository.PaymentRepository
	saleRepo     repository.SaleRepository
	customerRepo repository.CustomerRepository
	itemRepo     repository.ItemRepository
	pdfGenerator *pdf.Generator
}

// NewReportService creates a new ReportService
func NewReportService(
	loanRepo repository.LoanRepository,
	paymentRepo repository.PaymentRepository,
	saleRepo repository.SaleRepository,
	customerRepo repository.CustomerRepository,
	itemRepo repository.ItemRepository,
	pdfGenerator *pdf.Generator,
) *ReportService {
	return &ReportService{
		loanRepo:     loanRepo,
		paymentRepo:  paymentRepo,
		saleRepo:     saleRepo,
		customerRepo: customerRepo,
		itemRepo:     itemRepo,
		pdfGenerator: pdfGenerator,
	}
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	// Loan stats
	ActiveLoans      int     `json:"active_loans"`
	OverdueLoans     int     `json:"overdue_loans"`
	TotalLoanAmount  float64 `json:"total_loan_amount"`
	TotalOutstanding float64 `json:"total_outstanding"`

	// Today's activity
	NewLoansToday    int     `json:"new_loans_today"`
	PaymentsToday    int     `json:"payments_today"`
	PaymentAmount    float64 `json:"payment_amount_today"`
	SalesToday       int     `json:"sales_today"`
	SalesAmountToday float64 `json:"sales_amount_today"`

	// Inventory
	ItemsInPawn   int     `json:"items_in_pawn"`
	ItemsForSale  int     `json:"items_for_sale"`
	InventoryValue float64 `json:"inventory_value"`

	// Customers
	TotalCustomers int `json:"total_customers"`
	NewCustomers   int `json:"new_customers_today"`
}

// GetDashboardStats retrieves dashboard statistics
func (s *ReportService) GetDashboardStats(ctx context.Context, branchID int64) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Get loan statistics
	activeStatus := domain.LoanStatusActive
	loanParams := repository.LoanListParams{
		BranchID: branchID,
		Status:   &activeStatus,
		PaginationParams: repository.PaginationParams{
			PerPage: 1,
		},
	}

	activeLoans, err := s.loanRepo.List(ctx, loanParams)
	if err == nil {
		stats.ActiveLoans = activeLoans.Total
	}

	overdueStatus := domain.LoanStatusOverdue
	loanParams.Status = &overdueStatus
	overdueLoans, err := s.loanRepo.List(ctx, loanParams)
	if err == nil {
		stats.OverdueLoans = overdueLoans.Total
	}

	// Get all loans to calculate totals
	loanParams.Status = nil
	loanParams.PaginationParams.PerPage = 10000
	allLoans, err := s.loanRepo.List(ctx, loanParams)
	if err == nil {
		for _, loan := range allLoans.Data {
			if loan.Status == domain.LoanStatusActive || loan.Status == domain.LoanStatusOverdue {
				stats.TotalLoanAmount += loan.LoanAmount
				stats.TotalOutstanding += loan.RemainingBalance()
			}
		}
	}

	// Get today's payments
	today := time.Now().Format("2006-01-02")
	paymentParams := repository.PaymentListParams{
		BranchID: branchID,
		DateFrom: &today,
		DateTo:   &today,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	payments, err := s.paymentRepo.List(ctx, paymentParams)
	if err == nil {
		stats.PaymentsToday = len(payments.Data)
		for _, p := range payments.Data {
			if p.Status == domain.PaymentStatusCompleted {
				stats.PaymentAmount += p.Amount
			}
		}
	}

	// Get today's sales
	saleParams := repository.SaleListParams{
		BranchID: branchID,
		DateFrom: &today,
		DateTo:   &today,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	sales, err := s.saleRepo.List(ctx, saleParams)
	if err == nil {
		stats.SalesToday = len(sales.Data)
		for _, sale := range sales.Data {
			if sale.Status == domain.SaleStatusCompleted {
				stats.SalesAmountToday += sale.FinalPrice
			}
		}
	}

	// Get inventory stats
	itemParams := repository.ItemListParams{
		BranchID: branchID,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	items, err := s.itemRepo.List(ctx, itemParams)
	if err == nil {
		for _, item := range items.Data {
			if item.Status == domain.ItemStatusPawned {
				stats.ItemsInPawn++
				stats.InventoryValue += item.AppraisedValue
			} else if item.Status == domain.ItemStatusForSale {
				stats.ItemsForSale++
				stats.InventoryValue += item.AppraisedValue
			}
		}
	}

	// Get customer stats
	customerParams := repository.CustomerListParams{
		BranchID: branchID,
		PaginationParams: repository.PaginationParams{
			PerPage: 1,
		},
	}

	customers, err := s.customerRepo.List(ctx, customerParams)
	if err == nil {
		stats.TotalCustomers = customers.Total
	}

	return stats, nil
}

// LoanReport represents loan report data
type LoanReport struct {
	TotalLoans       int               `json:"total_loans"`
	TotalAmount      float64           `json:"total_amount"`
	TotalInterest    float64           `json:"total_interest"`
	TotalOutstanding float64           `json:"total_outstanding"`
	ByStatus         map[string]int    `json:"by_status"`
	ByStatusAmount   map[string]float64 `json:"by_status_amount"`
	RecentLoans      []domain.Loan     `json:"recent_loans,omitempty"`
}

// GetLoanReport generates a loan report
func (s *ReportService) GetLoanReport(ctx context.Context, branchID int64, dateFrom, dateTo string) (*LoanReport, error) {
	report := &LoanReport{
		ByStatus:       make(map[string]int),
		ByStatusAmount: make(map[string]float64),
	}

	params := repository.LoanListParams{
		BranchID: branchID,
		DueAfter: &dateFrom,
		DueBefore: &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
			OrderBy: "created_at",
			Order:   "desc",
		},
	}

	result, err := s.loanRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	report.TotalLoans = result.Total

	for _, loan := range result.Data {
		report.TotalAmount += loan.LoanAmount
		report.TotalInterest += loan.InterestAmount
		if loan.Status == domain.LoanStatusActive || loan.Status == domain.LoanStatusOverdue {
			report.TotalOutstanding += loan.RemainingBalance()
		}

		status := string(loan.Status)
		report.ByStatus[status]++
		report.ByStatusAmount[status] += loan.LoanAmount
	}

	// Get 10 most recent loans
	if len(result.Data) > 10 {
		report.RecentLoans = result.Data[:10]
	} else {
		report.RecentLoans = result.Data
	}

	return report, nil
}

// PaymentReport represents payment report data
type PaymentReport struct {
	TotalPayments    int                    `json:"total_payments"`
	TotalAmount      float64                `json:"total_amount"`
	TotalPrincipal   float64                `json:"total_principal"`
	TotalInterest    float64                `json:"total_interest"`
	TotalLateFees    float64                `json:"total_late_fees"`
	ByMethod         map[string]int         `json:"by_method"`
	ByMethodAmount   map[string]float64     `json:"by_method_amount"`
	RecentPayments   []domain.Payment       `json:"recent_payments,omitempty"`
}

// GetPaymentReport generates a payment report
func (s *ReportService) GetPaymentReport(ctx context.Context, branchID int64, dateFrom, dateTo string) (*PaymentReport, error) {
	report := &PaymentReport{
		ByMethod:       make(map[string]int),
		ByMethodAmount: make(map[string]float64),
	}

	params := repository.PaymentListParams{
		BranchID: branchID,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
			OrderBy: "payment_date",
			Order:   "desc",
		},
	}

	result, err := s.paymentRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	report.TotalPayments = result.Total

	for _, payment := range result.Data {
		if payment.Status != domain.PaymentStatusCompleted {
			continue
		}

		report.TotalAmount += payment.Amount
		report.TotalPrincipal += payment.PrincipalAmount
		report.TotalInterest += payment.InterestAmount
		report.TotalLateFees += payment.LateFeeAmount

		method := string(payment.PaymentMethod)
		report.ByMethod[method]++
		report.ByMethodAmount[method] += payment.Amount
	}

	// Get 10 most recent payments
	if len(result.Data) > 10 {
		report.RecentPayments = result.Data[:10]
	} else {
		report.RecentPayments = result.Data
	}

	return report, nil
}

// SalesReport represents sales report data
type SalesReport struct {
	TotalSales     int                `json:"total_sales"`
	TotalAmount    float64            `json:"total_amount"`
	TotalDiscounts float64            `json:"total_discounts"`
	NetAmount      float64            `json:"net_amount"`
	ByStatus       map[string]int     `json:"by_status"`
	ByMethod       map[string]int     `json:"by_method"`
	ByMethodAmount map[string]float64 `json:"by_method_amount"`
	RecentSales    []domain.Sale      `json:"recent_sales,omitempty"`
}

// GetSalesReport generates a sales report
func (s *ReportService) GetSalesReport(ctx context.Context, branchID int64, dateFrom, dateTo string) (*SalesReport, error) {
	report := &SalesReport{
		ByStatus:       make(map[string]int),
		ByMethod:       make(map[string]int),
		ByMethodAmount: make(map[string]float64),
	}

	params := repository.SaleListParams{
		BranchID: branchID,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
			OrderBy: "sale_date",
			Order:   "desc",
		},
	}

	result, err := s.saleRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	report.TotalSales = result.Total

	for _, sale := range result.Data {
		report.ByStatus[string(sale.Status)]++

		if sale.Status == domain.SaleStatusCompleted {
			report.TotalAmount += sale.SalePrice
			report.TotalDiscounts += sale.DiscountAmount
			report.NetAmount += sale.FinalPrice

			method := string(sale.PaymentMethod)
			report.ByMethod[method]++
			report.ByMethodAmount[method] += sale.FinalPrice
		}
	}

	// Get 10 most recent sales
	if len(result.Data) > 10 {
		report.RecentSales = result.Data[:10]
	} else {
		report.RecentSales = result.Data
	}

	return report, nil
}

// OverdueReport represents overdue loans report
type OverdueReport struct {
	TotalOverdue     int           `json:"total_overdue"`
	TotalAmount      float64       `json:"total_amount"`
	TotalLateFees    float64       `json:"total_late_fees"`
	OverdueLoans     []domain.Loan `json:"overdue_loans"`
	ApproachingDue   []domain.Loan `json:"approaching_due"`
	AboutToDefault   []domain.Loan `json:"about_to_default"`
}

// GetOverdueReport generates an overdue loans report
func (s *ReportService) GetOverdueReport(ctx context.Context, branchID int64) (*OverdueReport, error) {
	report := &OverdueReport{
		OverdueLoans:   []domain.Loan{},
		ApproachingDue: []domain.Loan{},
		AboutToDefault: []domain.Loan{},
	}

	// Get overdue loans
	overdueLoans, err := s.loanRepo.GetOverdueLoans(ctx, branchID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	for _, loan := range overdueLoans {
		if loan.Status == domain.LoanStatusOverdue {
			report.TotalOverdue++
			report.TotalAmount += loan.RemainingBalance()
			report.TotalLateFees += loan.LateFeeAmount
			report.OverdueLoans = append(report.OverdueLoans, *loan)

			// Check if about to default
			gracePeriodEnd := loan.DueDate.AddDate(0, 0, loan.GracePeriodDays)
			daysUntilDefault := int(gracePeriodEnd.Sub(now).Hours() / 24)
			if daysUntilDefault <= 3 && daysUntilDefault > 0 {
				report.AboutToDefault = append(report.AboutToDefault, *loan)
			}
		}
	}

	// Get loans approaching due date (within 7 days)
	activeStatus := domain.LoanStatusActive
	params := repository.LoanListParams{
		BranchID: branchID,
		Status:   &activeStatus,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	activeLoans, err := s.loanRepo.List(ctx, params)
	if err == nil {
		for _, loan := range activeLoans.Data {
			daysUntilDue := int(loan.DueDate.Sub(now).Hours() / 24)
			if daysUntilDue >= 0 && daysUntilDue <= 7 {
				report.ApproachingDue = append(report.ApproachingDue, loan)
			}
		}
	}

	return report, nil
}

// GenerateDailyReportPDF generates a daily report PDF
func (s *ReportService) GenerateDailyReportPDF(ctx context.Context, branchID int64, date time.Time) ([]byte, error) {
	dateStr := date.Format("2006-01-02")

	// Gather data for the report
	dailyReport := &pdf.DailyReport{
		Date:     date,
		BranchID: branchID,
	}

	// Get loans created on this date
	loanParams := repository.LoanListParams{
		BranchID: branchID,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	loans, _ := s.loanRepo.List(ctx, loanParams)
	for _, loan := range loans.Data {
		if loan.CreatedAt.Format("2006-01-02") == dateStr {
			dailyReport.NewLoansCount++
			dailyReport.NewLoansAmount += loan.LoanAmount
		}
		if loan.Status == domain.LoanStatusRenewed && loan.UpdatedAt.Format("2006-01-02") == dateStr {
			dailyReport.RenewalsCount++
		}
	}

	// Get overdue loans
	overdueLoans, _ := s.loanRepo.GetOverdueLoans(ctx, branchID)
	for _, loan := range overdueLoans {
		if loan.Status == domain.LoanStatusOverdue {
			dailyReport.OverdueCount++
			dailyReport.OverdueAmount += loan.RemainingBalance()
		}
	}

	// Get payments
	paymentParams := repository.PaymentListParams{
		BranchID: branchID,
		DateFrom: &dateStr,
		DateTo:   &dateStr,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	payments, _ := s.paymentRepo.List(ctx, paymentParams)
	for _, p := range payments.Data {
		if p.Status == domain.PaymentStatusCompleted {
			dailyReport.PaymentsCount++
			dailyReport.PaymentsAmount += p.Amount
		}
	}

	// Get sales
	saleParams := repository.SaleListParams{
		BranchID: branchID,
		DateFrom: &dateStr,
		DateTo:   &dateStr,
		PaginationParams: repository.PaginationParams{
			PerPage: 10000,
		},
	}

	sales, _ := s.saleRepo.List(ctx, saleParams)
	for _, sale := range sales.Data {
		if sale.Status == domain.SaleStatusCompleted {
			dailyReport.SalesCount++
			dailyReport.SalesAmount += sale.FinalPrice
		}
	}

	// Calculate income/expenses
	dailyReport.TotalIncome = dailyReport.PaymentsAmount + dailyReport.SalesAmount
	dailyReport.TotalExpenses = dailyReport.NewLoansAmount

	return s.pdfGenerator.GenerateDailyReport(dailyReport)
}

// GenerateLoanContractPDF generates a loan contract PDF
func (s *ReportService) GenerateLoanContractPDF(ctx context.Context, loanID int64) ([]byte, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(ctx, loan.CustomerID)
	if err != nil {
		return nil, err
	}

	item, err := s.itemRepo.GetByID(ctx, loan.ItemID)
	if err != nil {
		return nil, err
	}

	return s.pdfGenerator.GenerateLoanContract(loan, customer, item)
}

// GeneratePaymentReceiptPDF generates a payment receipt PDF
func (s *ReportService) GeneratePaymentReceiptPDF(ctx context.Context, paymentID int64) ([]byte, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	loan, err := s.loanRepo.GetByID(ctx, payment.LoanID)
	if err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(ctx, payment.CustomerID)
	if err != nil {
		return nil, err
	}

	return s.pdfGenerator.GeneratePaymentReceipt(payment, loan, customer)
}

// GenerateSaleReceiptPDF generates a sale receipt PDF
func (s *ReportService) GenerateSaleReceiptPDF(ctx context.Context, saleID int64) ([]byte, error) {
	sale, err := s.saleRepo.GetByID(ctx, saleID)
	if err != nil {
		return nil, err
	}

	item, err := s.itemRepo.GetByID(ctx, sale.ItemID)
	if err != nil {
		return nil, err
	}

	var customer *domain.Customer
	if sale.CustomerID != nil {
		customer, _ = s.customerRepo.GetByID(ctx, *sale.CustomerID)
	}

	return s.pdfGenerator.GenerateSaleReceipt(sale, item, customer)
}

// ExportLoanReportPDF exports the loan report as PDF
func (s *ReportService) ExportLoanReportPDF(ctx context.Context, branchID int64, dateFrom, dateTo string) ([]byte, error) {
	report, err := s.GetLoanReport(ctx, branchID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Get loans with customer info for details
	params := repository.LoanListParams{
		BranchID:  branchID,
		DueAfter:  &dateFrom,
		DueBefore: &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 100,
			OrderBy: "created_at",
			Order:   "desc",
		},
	}

	result, err := s.loanRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// Build report items
	items := make([]pdf.LoanReportItem, 0, len(result.Data))
	for _, loan := range result.Data {
		customerName := "N/A"
		if loan.Customer != nil {
			customerName = loan.Customer.FirstName + " " + loan.Customer.LastName
		}
		itemName := "N/A"
		if loan.Item != nil {
			itemName = loan.Item.Name
		}

		items = append(items, pdf.LoanReportItem{
			LoanNumber:   loan.LoanNumber,
			CustomerName: customerName,
			ItemName:     itemName,
			Amount:       loan.LoanAmount,
			Interest:     loan.InterestAmount,
			Total:        loan.TotalAmount,
			Status:       string(loan.Status),
			DueDate:      loan.DueDate.Format("02/01/06"),
		})
	}

	data := &pdf.LoanReportData{
		DateFrom:         dateFrom,
		DateTo:           dateTo,
		TotalLoans:       report.TotalLoans,
		TotalAmount:      report.TotalAmount,
		TotalInterest:    report.TotalInterest,
		TotalOutstanding: report.TotalOutstanding,
		ByStatus:         report.ByStatus,
		ByStatusAmount:   report.ByStatusAmount,
		Loans:            items,
	}

	return s.pdfGenerator.GenerateLoanReportPDF(data)
}

// ExportPaymentReportPDF exports the payment report as PDF
func (s *ReportService) ExportPaymentReportPDF(ctx context.Context, branchID int64, dateFrom, dateTo string) ([]byte, error) {
	report, err := s.GetPaymentReport(ctx, branchID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Get payments with details
	params := repository.PaymentListParams{
		BranchID: branchID,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 100,
			OrderBy: "payment_date",
			Order:   "desc",
		},
	}

	result, err := s.paymentRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// Build report items
	items := make([]pdf.PaymentReportItem, 0, len(result.Data))
	for _, payment := range result.Data {
		if payment.Status != domain.PaymentStatusCompleted {
			continue
		}

		customerName := "N/A"
		customer, _ := s.customerRepo.GetByID(ctx, payment.CustomerID)
		if customer != nil {
			customerName = customer.FirstName + " " + customer.LastName
		}

		loanNumber := "N/A"
		loan, _ := s.loanRepo.GetByID(ctx, payment.LoanID)
		if loan != nil {
			loanNumber = loan.LoanNumber
		}

		items = append(items, pdf.PaymentReportItem{
			PaymentNumber: payment.PaymentNumber,
			CustomerName:  customerName,
			LoanNumber:    loanNumber,
			Amount:        payment.Amount,
			Principal:     payment.PrincipalAmount,
			Interest:      payment.InterestAmount,
			LateFee:       payment.LateFeeAmount,
			Method:        string(payment.PaymentMethod),
			Date:          payment.PaymentDate.Format("02/01/06"),
		})
	}

	data := &pdf.PaymentReportData{
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		TotalPayments:  report.TotalPayments,
		TotalAmount:    report.TotalAmount,
		TotalPrincipal: report.TotalPrincipal,
		TotalInterest:  report.TotalInterest,
		TotalLateFees:  report.TotalLateFees,
		ByMethod:       report.ByMethod,
		ByMethodAmount: report.ByMethodAmount,
		Payments:       items,
	}

	return s.pdfGenerator.GeneratePaymentReportPDF(data)
}

// ExportSalesReportPDF exports the sales report as PDF
func (s *ReportService) ExportSalesReportPDF(ctx context.Context, branchID int64, dateFrom, dateTo string) ([]byte, error) {
	report, err := s.GetSalesReport(ctx, branchID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Get sales with details
	params := repository.SaleListParams{
		BranchID: branchID,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		PaginationParams: repository.PaginationParams{
			PerPage: 100,
			OrderBy: "sale_date",
			Order:   "desc",
		},
	}

	result, err := s.saleRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// Build report items
	items := make([]pdf.SaleReportItem, 0, len(result.Data))
	for _, sale := range result.Data {
		if sale.Status != domain.SaleStatusCompleted {
			continue
		}

		customerName := "Sin cliente"
		if sale.CustomerID != nil {
			customer, _ := s.customerRepo.GetByID(ctx, *sale.CustomerID)
			if customer != nil {
				customerName = customer.FirstName + " " + customer.LastName
			}
		}

		itemName := "N/A"
		item, _ := s.itemRepo.GetByID(ctx, sale.ItemID)
		if item != nil {
			itemName = item.Name
		}

		items = append(items, pdf.SaleReportItem{
			SaleNumber:   sale.SaleNumber,
			CustomerName: customerName,
			ItemName:     itemName,
			Price:        sale.SalePrice,
			Discount:     sale.DiscountAmount,
			Total:        sale.FinalPrice,
			Method:       string(sale.PaymentMethod),
			Date:         sale.SaleDate.Format("02/01/06"),
		})
	}

	data := &pdf.SalesReportData{
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		TotalSales:     report.TotalSales,
		TotalAmount:    report.TotalAmount,
		TotalDiscounts: report.TotalDiscounts,
		NetAmount:      report.NetAmount,
		ByMethod:       report.ByMethod,
		ByMethodAmount: report.ByMethodAmount,
		Sales:          items,
	}

	return s.pdfGenerator.GenerateSalesReportPDF(data)
}

// ExportOverdueReportPDF exports the overdue report as PDF
func (s *ReportService) ExportOverdueReportPDF(ctx context.Context, branchID int64) ([]byte, error) {
	report, err := s.GetOverdueReport(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Build overdue items
	overdueItems := make([]pdf.OverdueReportItem, 0, len(report.OverdueLoans))
	for _, loan := range report.OverdueLoans {
		customerName := "N/A"
		customer, _ := s.customerRepo.GetByID(ctx, loan.CustomerID)
		if customer != nil {
			customerName = customer.FirstName + " " + customer.LastName
		}

		itemName := "N/A"
		item, _ := s.itemRepo.GetByID(ctx, loan.ItemID)
		if item != nil {
			itemName = item.Name
		}

		graceEnds := loan.DueDate.AddDate(0, 0, loan.GracePeriodDays).Format("02/01/06")

		overdueItems = append(overdueItems, pdf.OverdueReportItem{
			LoanNumber:   loan.LoanNumber,
			CustomerName: customerName,
			ItemName:     itemName,
			Amount:       loan.RemainingBalance(),
			LateFee:      loan.LateFeeAmount,
			DaysOverdue:  loan.DaysOverdue,
			DueDate:      loan.DueDate.Format("02/01/06"),
			GraceEnds:    graceEnds,
		})
	}

	// Build approaching due items
	approachingItems := make([]pdf.OverdueReportItem, 0, len(report.ApproachingDue))
	for _, loan := range report.ApproachingDue {
		customerName := "N/A"
		customer, _ := s.customerRepo.GetByID(ctx, loan.CustomerID)
		if customer != nil {
			customerName = customer.FirstName + " " + customer.LastName
		}

		approachingItems = append(approachingItems, pdf.OverdueReportItem{
			LoanNumber:   loan.LoanNumber,
			CustomerName: customerName,
			Amount:       loan.RemainingBalance(),
			DueDate:      loan.DueDate.Format("02/01/06"),
		})
	}

	data := &pdf.OverdueReportData{
		GeneratedAt:    time.Now().Format("02/01/2006 15:04"),
		TotalOverdue:   report.TotalOverdue,
		TotalAmount:    report.TotalAmount,
		TotalLateFees:  report.TotalLateFees,
		OverdueLoans:   overdueItems,
		ApproachingDue: approachingItems,
	}

	return s.pdfGenerator.GenerateOverdueReportPDF(data)
}
