package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"pawnshop/internal/domain"
)

// Generator handles PDF document generation
type Generator struct {
	companyName    string
	companyAddress string
	companyPhone   string
}

// NewGenerator creates a new PDF generator
func NewGenerator(companyName, companyAddress, companyPhone string) *Generator {
	return &Generator{
		companyName:    companyName,
		companyAddress: companyAddress,
		companyPhone:   companyPhone,
	}
}

// GenerateLoanContract generates a loan contract PDF
func (g *Generator) GenerateLoanContract(loan *domain.Loan, customer *domain.Customer, item *domain.Item) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)
	doc := m.GetStructure()

	// Header
	g.addHeader(m, "CONTRATO DE EMPEÑO")

	// Contract Info
	m.AddRow(8, text.NewCol(12, fmt.Sprintf("Contrato No: %s", loan.LoanNumber), props.Text{
		Size:  12,
		Style: fontstyle.Bold,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Fecha: %s", loan.StartDate.Format("02/01/2006")), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Fecha de Vencimiento: %s", loan.DueDate.Format("02/01/2006")), props.Text{Size: 10}))

	// Customer Section
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "DATOS DEL CLIENTE", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Nombre: %s %s", customer.FirstName, customer.LastName), props.Text{Size: 10}))
	if customer.IdentityNumber != "" {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Identificación: %s - %s", customer.IdentityType, customer.IdentityNumber), props.Text{Size: 10}))
	}
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Teléfono: %s", customer.Phone), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Dirección: %s", customer.Address), props.Text{Size: 10}))

	// Item Section
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "ARTÍCULO EN PRENDA", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Descripción: %s", item.Name), props.Text{Size: 10}))
	if item.Brand != nil {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Marca: %s", *item.Brand), props.Text{Size: 10}))
	}
	if item.Model != nil {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Modelo: %s", *item.Model), props.Text{Size: 10}))
	}
	if item.SerialNumber != nil {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("No. Serie: %s", *item.SerialNumber), props.Text{Size: 10}))
	}
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Condición: %s", item.Condition), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Valor Avalúo: $%.2f", item.AppraisedValue), props.Text{Size: 10}))

	// Loan Details
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "CONDICIONES DEL PRÉSTAMO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Monto del Préstamo: $%.2f", loan.LoanAmount), props.Text{Size: 10, Style: fontstyle.Bold}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Tasa de Interés: %.2f%% mensual", loan.InterestRate), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Interés: $%.2f", loan.InterestAmount), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Total a Pagar: $%.2f", loan.TotalAmount), props.Text{Size: 10, Style: fontstyle.Bold}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Plazo: %d días", loan.LoanTermDays), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Período de Gracia: %d días", loan.GracePeriodDays), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Mora por día vencido: %.2f%%", loan.LateFeeRate), props.Text{Size: 10}))

	// Terms and Conditions
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "TÉRMINOS Y CONDICIONES", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	terms := []string{
		"1. El cliente acepta las condiciones del préstamo establecidas en este contrato.",
		"2. El artículo quedará en custodia hasta el pago total del préstamo.",
		"3. Si el préstamo no es pagado en la fecha de vencimiento, se aplicarán cargos por mora.",
		fmt.Sprintf("4. Después de %d días de vencido el período de gracia, el artículo pasará a propiedad de la casa de empeño.", loan.GracePeriodDays),
		"5. El cliente puede renovar el préstamo pagando los intereses acumulados.",
	}

	for _, term := range terms {
		m.AddRow(5, text.NewCol(12, term, props.Text{Size: 9}))
	}

	// Signatures
	m.AddRow(25)
	m.AddRows(
		row.New(15).Add(
			col.New(6).Add(
				text.New("_______________________", props.Text{Size: 10, Align: align.Center}),
			),
			col.New(6).Add(
				text.New("_______________________", props.Text{Size: 10, Align: align.Center}),
			),
		),
		row.New(6).Add(
			col.New(6).Add(
				text.New("Firma del Cliente", props.Text{Size: 9, Align: align.Center}),
			),
			col.New(6).Add(
				text.New("Firma Autorizada", props.Text{Size: 9, Align: align.Center}),
			),
		),
	)

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	_ = doc // Avoid unused variable warning
	return document.GetBytes(), nil
}

// GeneratePaymentReceipt generates a payment receipt PDF
func (g *Generator) GeneratePaymentReceipt(payment *domain.Payment, loan *domain.Loan, customer *domain.Customer) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "RECIBO DE PAGO")

	// Receipt Info
	m.AddRow(8, text.NewCol(12, fmt.Sprintf("Recibo No: %s", payment.PaymentNumber), props.Text{
		Size:  12,
		Style: fontstyle.Bold,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Fecha: %s", payment.PaymentDate.Format("02/01/2006 15:04")), props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Contrato: %s", loan.LoanNumber), props.Text{Size: 10}))

	// Customer Info
	m.AddRow(10)
	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Cliente: %s %s", customer.FirstName, customer.LastName), props.Text{Size: 10}))

	// Payment Details
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "DETALLES DEL PAGO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Monto del Pago: $%.2f", payment.Amount), props.Text{Size: 10, Style: fontstyle.Bold}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Método de Pago: %s", payment.PaymentMethod), props.Text{Size: 10}))

	if payment.PrincipalAmount > 0 {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Aplicado a Capital: $%.2f", payment.PrincipalAmount), props.Text{Size: 10}))
	}
	if payment.InterestAmount > 0 {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Aplicado a Intereses: $%.2f", payment.InterestAmount), props.Text{Size: 10}))
	}
	if payment.LateFeeAmount > 0 {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Aplicado a Mora: $%.2f", payment.LateFeeAmount), props.Text{Size: 10}))
	}

	// Balance After Payment
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "SALDO DESPUÉS DEL PAGO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	balanceRemaining := loan.RemainingBalance()
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Saldo Pendiente: $%.2f", balanceRemaining), props.Text{Size: 10, Style: fontstyle.Bold}))

	if loan.Status == domain.LoanStatusPaid {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "*** PRÉSTAMO LIQUIDADO ***", props.Text{
			Size:  14,
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
	}

	// Footer
	m.AddRow(20)
	m.AddRow(6, text.NewCol(12, "Gracias por su pago.", props.Text{Size: 10, Align: align.Center}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// GenerateSaleReceipt generates a sale receipt PDF
func (g *Generator) GenerateSaleReceipt(sale *domain.Sale, item *domain.Item, customer *domain.Customer) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "RECIBO DE VENTA")

	// Sale Info
	m.AddRow(8, text.NewCol(12, fmt.Sprintf("Venta No: %s", sale.SaleNumber), props.Text{
		Size:  12,
		Style: fontstyle.Bold,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Fecha: %s", sale.SaleDate.Format("02/01/2006 15:04")), props.Text{Size: 10}))

	// Customer Info (if available)
	if customer != nil {
		m.AddRow(6, text.NewCol(12, fmt.Sprintf("Cliente: %s %s", customer.FirstName, customer.LastName), props.Text{Size: 10}))
	}

	// Item Details
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "ARTÍCULO VENDIDO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Descripción: %s", item.Name), props.Text{Size: 10}))
	if item.Brand != nil {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Marca: %s", *item.Brand), props.Text{Size: 10}))
	}
	if item.Model != nil {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Modelo: %s", *item.Model), props.Text{Size: 10}))
	}
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Condición: %s", item.Condition), props.Text{Size: 10}))

	// Price Details
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "DETALLES DE LA VENTA", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
		Top:   2,
	}))

	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Precio: $%.2f", sale.SalePrice), props.Text{Size: 10}))
	if sale.DiscountAmount > 0 {
		m.AddRow(6, text.NewCol(6, fmt.Sprintf("Descuento: $%.2f", sale.DiscountAmount), props.Text{Size: 10}))
	}
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Total: $%.2f", sale.FinalPrice), props.Text{Size: 10, Style: fontstyle.Bold}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("Método de Pago: %s", sale.PaymentMethod), props.Text{Size: 10}))

	// Footer
	m.AddRow(20)
	m.AddRow(6, text.NewCol(12, "Gracias por su compra.", props.Text{Size: 10, Align: align.Center}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// addHeader adds a common header to the document
func (g *Generator) addHeader(m core.Maroto, title string) {
	m.AddRow(10, text.NewCol(12, g.companyName, props.Text{
		Size:  16,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRow(5, text.NewCol(12, g.companyAddress, props.Text{
		Size:  9,
		Align: align.Center,
	}))

	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Tel: %s", g.companyPhone), props.Text{
		Size:  9,
		Align: align.Center,
	}))

	m.AddRow(10)

	m.AddRow(8, text.NewCol(12, title, props.Text{
		Size:  14,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRow(5)
}

// GenerateDailyReport generates a daily summary report PDF
func (g *Generator) GenerateDailyReport(report *DailyReport) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "REPORTE DIARIO")

	// Report Date
	m.AddRow(8, text.NewCol(12, fmt.Sprintf("Fecha: %s", report.Date.Format("02/01/2006")), props.Text{
		Size:  12,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRow(10)

	// Summary Section
	m.AddRow(8, text.NewCol(12, "RESUMEN DE OPERACIONES", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRow(6, text.NewCol(6, "Préstamos nuevos:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("%d ($%.2f)", report.NewLoansCount, report.NewLoansAmount), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Pagos recibidos:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("%d ($%.2f)", report.PaymentsCount, report.PaymentsAmount), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Ventas realizadas:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("%d ($%.2f)", report.SalesCount, report.SalesAmount), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Renovaciones:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("%d", report.RenewalsCount), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Préstamos vencidos:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("%d ($%.2f)", report.OverdueCount, report.OverdueAmount), props.Text{Size: 10}))

	// Cash Summary
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "RESUMEN DE CAJA", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRow(6, text.NewCol(6, "Efectivo inicial:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("$%.2f", report.OpeningCash), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Ingresos:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("$%.2f", report.TotalIncome), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Egresos:", props.Text{Size: 10}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("$%.2f", report.TotalExpenses), props.Text{Size: 10}))

	m.AddRow(6, text.NewCol(6, "Efectivo final:", props.Text{Size: 10, Style: fontstyle.Bold}))
	m.AddRow(6, text.NewCol(6, fmt.Sprintf("$%.2f", report.ClosingCash), props.Text{Size: 10, Style: fontstyle.Bold}))

	// Generated timestamp
	m.AddRow(20)
	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
		Size:  8,
		Align: align.Right,
	}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// DailyReport contains daily report data
type DailyReport struct {
	Date           time.Time
	BranchID       int64
	NewLoansCount  int
	NewLoansAmount float64
	PaymentsCount  int
	PaymentsAmount float64
	SalesCount     int
	SalesAmount    float64
	RenewalsCount  int
	OverdueCount   int
	OverdueAmount  float64
	OpeningCash    float64
	ClosingCash    float64
	TotalIncome    float64
	TotalExpenses  float64
}

// SaveToBuffer saves the PDF to a buffer
func SaveToBuffer(data []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(data)
	return buf, nil
}

// LoanReportData contains loan report data for PDF generation
type LoanReportData struct {
	DateFrom         string
	DateTo           string
	TotalLoans       int
	TotalAmount      float64
	TotalInterest    float64
	TotalOutstanding float64
	ByStatus         map[string]int
	ByStatusAmount   map[string]float64
	Loans            []LoanReportItem
}

// LoanReportItem represents a loan in the report
type LoanReportItem struct {
	LoanNumber   string
	CustomerName string
	ItemName     string
	Amount       float64
	Interest     float64
	Total        float64
	Status       string
	DueDate      string
}

// GenerateLoanReportPDF generates a loan report PDF
func (g *Generator) GenerateLoanReportPDF(data *LoanReportData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "REPORTE DE PRÉSTAMOS")

	// Date range
	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Período: %s - %s", data.DateFrom, data.DateTo), props.Text{
		Size:  10,
		Align: align.Center,
	}))

	m.AddRow(10)

	// Summary
	m.AddRow(8, text.NewCol(12, "RESUMEN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(text.New("Total de Préstamos:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", data.TotalLoans), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Monto Total Prestado:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalAmount), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Intereses Totales:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalInterest), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Saldo Pendiente:", props.Text{Size: 10, Style: fontstyle.Bold})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalOutstanding), props.Text{Size: 10, Style: fontstyle.Bold, Align: align.Right})),
		),
	)

	// By status
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "POR ESTADO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	statusLabels := map[string]string{
		"active":      "Activos",
		"paid":        "Pagados",
		"overdue":     "Vencidos",
		"defaulted":   "En Mora",
		"renewed":     "Renovados",
		"confiscated": "Confiscados",
	}

	for status, count := range data.ByStatus {
		label := statusLabels[status]
		if label == "" {
			label = status
		}
		amount := data.ByStatusAmount[status]
		m.AddRow(5,
			text.NewCol(4, label, props.Text{Size: 9}),
			text.NewCol(4, fmt.Sprintf("%d", count), props.Text{Size: 9, Align: align.Center}),
			text.NewCol(4, fmt.Sprintf("Q%.2f", amount), props.Text{Size: 9, Align: align.Right}),
		)
	}

	// Loan list
	if len(data.Loans) > 0 {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "DETALLE DE PRÉSTAMOS", props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}))

		// Table header
		m.AddRow(6,
			text.NewCol(2, "No.", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(3, "Cliente", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(2, "Monto", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Total", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Estado", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(1, "Vence", props.Text{Size: 8, Style: fontstyle.Bold}),
		)

		for _, loan := range data.Loans {
			statusLabel := statusLabels[loan.Status]
			if statusLabel == "" {
				statusLabel = loan.Status
			}
			m.AddRow(5,
				text.NewCol(2, loan.LoanNumber, props.Text{Size: 7}),
				text.NewCol(3, truncateName(loan.CustomerName, 20), props.Text{Size: 7}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", loan.Amount), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", loan.Total), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, statusLabel, props.Text{Size: 7, Align: align.Center}),
				text.NewCol(1, loan.DueDate, props.Text{Size: 7}),
			)
		}
	}

	// Footer
	m.AddRow(15)
	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
		Size:  8,
		Align: align.Right,
	}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// PaymentReportData contains payment report data for PDF generation
type PaymentReportData struct {
	DateFrom       string
	DateTo         string
	TotalPayments  int
	TotalAmount    float64
	TotalPrincipal float64
	TotalInterest  float64
	TotalLateFees  float64
	ByMethod       map[string]int
	ByMethodAmount map[string]float64
	Payments       []PaymentReportItem
}

// PaymentReportItem represents a payment in the report
type PaymentReportItem struct {
	PaymentNumber string
	CustomerName  string
	LoanNumber    string
	Amount        float64
	Principal     float64
	Interest      float64
	LateFee       float64
	Method        string
	Date          string
}

// GeneratePaymentReportPDF generates a payment report PDF
func (g *Generator) GeneratePaymentReportPDF(data *PaymentReportData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "REPORTE DE PAGOS")

	// Date range
	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Período: %s - %s", data.DateFrom, data.DateTo), props.Text{
		Size:  10,
		Align: align.Center,
	}))

	m.AddRow(10)

	// Summary
	m.AddRow(8, text.NewCol(12, "RESUMEN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(text.New("Total de Pagos:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", data.TotalPayments), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Monto Total:", props.Text{Size: 10, Style: fontstyle.Bold})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalAmount), props.Text{Size: 10, Style: fontstyle.Bold, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Capital Recuperado:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalPrincipal), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Intereses Cobrados:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalInterest), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Moras Cobradas:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalLateFees), props.Text{Size: 10, Align: align.Right})),
		),
	)

	// By method
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "POR MÉTODO DE PAGO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	methodLabels := map[string]string{
		"cash":     "Efectivo",
		"card":     "Tarjeta",
		"transfer": "Transferencia",
	}

	for method, count := range data.ByMethod {
		label := methodLabels[method]
		if label == "" {
			label = method
		}
		amount := data.ByMethodAmount[method]
		m.AddRow(5,
			text.NewCol(4, label, props.Text{Size: 9}),
			text.NewCol(4, fmt.Sprintf("%d", count), props.Text{Size: 9, Align: align.Center}),
			text.NewCol(4, fmt.Sprintf("Q%.2f", amount), props.Text{Size: 9, Align: align.Right}),
		)
	}

	// Payment list
	if len(data.Payments) > 0 {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "DETALLE DE PAGOS", props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}))

		// Table header
		m.AddRow(6,
			text.NewCol(2, "Recibo", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(3, "Cliente", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(2, "Préstamo", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(2, "Monto", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Método", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(1, "Fecha", props.Text{Size: 8, Style: fontstyle.Bold}),
		)

		for _, payment := range data.Payments {
			methodLabel := methodLabels[payment.Method]
			if methodLabel == "" {
				methodLabel = payment.Method
			}
			m.AddRow(5,
				text.NewCol(2, payment.PaymentNumber, props.Text{Size: 7}),
				text.NewCol(3, truncateName(payment.CustomerName, 20), props.Text{Size: 7}),
				text.NewCol(2, payment.LoanNumber, props.Text{Size: 7}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", payment.Amount), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, methodLabel, props.Text{Size: 7, Align: align.Center}),
				text.NewCol(1, payment.Date, props.Text{Size: 7}),
			)
		}
	}

	// Footer
	m.AddRow(15)
	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
		Size:  8,
		Align: align.Right,
	}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// SalesReportData contains sales report data for PDF generation
type SalesReportData struct {
	DateFrom       string
	DateTo         string
	TotalSales     int
	TotalAmount    float64
	TotalDiscounts float64
	NetAmount      float64
	ByMethod       map[string]int
	ByMethodAmount map[string]float64
	Sales          []SaleReportItem
}

// SaleReportItem represents a sale in the report
type SaleReportItem struct {
	SaleNumber   string
	CustomerName string
	ItemName     string
	Price        float64
	Discount     float64
	Total        float64
	Method       string
	Date         string
}

// GenerateSalesReportPDF generates a sales report PDF
func (g *Generator) GenerateSalesReportPDF(data *SalesReportData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "REPORTE DE VENTAS")

	// Date range
	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Período: %s - %s", data.DateFrom, data.DateTo), props.Text{
		Size:  10,
		Align: align.Center,
	}))

	m.AddRow(10)

	// Summary
	m.AddRow(8, text.NewCol(12, "RESUMEN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(text.New("Total de Ventas:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", data.TotalSales), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Monto Bruto:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalAmount), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Descuentos:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("-Q%.2f", data.TotalDiscounts), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Monto Neto:", props.Text{Size: 10, Style: fontstyle.Bold})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.NetAmount), props.Text{Size: 10, Style: fontstyle.Bold, Align: align.Right})),
		),
	)

	// By method
	m.AddRow(10)
	m.AddRow(8, text.NewCol(12, "POR MÉTODO DE PAGO", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	methodLabels := map[string]string{
		"cash":     "Efectivo",
		"card":     "Tarjeta",
		"transfer": "Transferencia",
	}

	for method, count := range data.ByMethod {
		label := methodLabels[method]
		if label == "" {
			label = method
		}
		amount := data.ByMethodAmount[method]
		m.AddRow(5,
			text.NewCol(4, label, props.Text{Size: 9}),
			text.NewCol(4, fmt.Sprintf("%d", count), props.Text{Size: 9, Align: align.Center}),
			text.NewCol(4, fmt.Sprintf("Q%.2f", amount), props.Text{Size: 9, Align: align.Right}),
		)
	}

	// Sales list
	if len(data.Sales) > 0 {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "DETALLE DE VENTAS", props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}))

		// Table header
		m.AddRow(6,
			text.NewCol(2, "No.", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(3, "Artículo", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(2, "Precio", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Total", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Método", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(1, "Fecha", props.Text{Size: 8, Style: fontstyle.Bold}),
		)

		for _, sale := range data.Sales {
			methodLabel := methodLabels[sale.Method]
			if methodLabel == "" {
				methodLabel = sale.Method
			}
			m.AddRow(5,
				text.NewCol(2, sale.SaleNumber, props.Text{Size: 7}),
				text.NewCol(3, truncateName(sale.ItemName, 20), props.Text{Size: 7}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", sale.Price), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", sale.Total), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, methodLabel, props.Text{Size: 7, Align: align.Center}),
				text.NewCol(1, sale.Date, props.Text{Size: 7}),
			)
		}
	}

	// Footer
	m.AddRow(15)
	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
		Size:  8,
		Align: align.Right,
	}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// OverdueReportData contains overdue report data for PDF generation
type OverdueReportData struct {
	GeneratedAt    string
	TotalOverdue   int
	TotalAmount    float64
	TotalLateFees  float64
	OverdueLoans   []OverdueReportItem
	ApproachingDue []OverdueReportItem
}

// OverdueReportItem represents an overdue loan in the report
type OverdueReportItem struct {
	LoanNumber   string
	CustomerName string
	ItemName     string
	Amount       float64
	LateFee      float64
	DaysOverdue  int
	DueDate      string
	GraceEnds    string
}

// GenerateOverdueReportPDF generates an overdue loans report PDF
func (g *Generator) GenerateOverdueReportPDF(data *OverdueReportData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	m := maroto.New(cfg)

	// Header
	g.addHeader(m, "REPORTE DE PRÉSTAMOS VENCIDOS")

	m.AddRow(6, text.NewCol(12, fmt.Sprintf("Generado: %s", data.GeneratedAt), props.Text{
		Size:  10,
		Align: align.Center,
	}))

	m.AddRow(10)

	// Summary
	m.AddRow(8, text.NewCol(12, "RESUMEN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	}))

	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(text.New("Total Préstamos Vencidos:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", data.TotalOverdue), props.Text{Size: 10, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Monto Pendiente:", props.Text{Size: 10, Style: fontstyle.Bold})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalAmount), props.Text{Size: 10, Style: fontstyle.Bold, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(6).Add(text.New("Moras Acumuladas:", props.Text{Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("Q%.2f", data.TotalLateFees), props.Text{Size: 10, Align: align.Right})),
		),
	)

	// Overdue loans list
	if len(data.OverdueLoans) > 0 {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "PRÉSTAMOS VENCIDOS", props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}))

		// Table header
		m.AddRow(6,
			text.NewCol(2, "No.", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(3, "Cliente", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(2, "Saldo", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(2, "Mora", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(1, "Días", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(2, "Gracia", props.Text{Size: 8, Style: fontstyle.Bold}),
		)

		for _, loan := range data.OverdueLoans {
			m.AddRow(5,
				text.NewCol(2, loan.LoanNumber, props.Text{Size: 7}),
				text.NewCol(3, truncateName(loan.CustomerName, 20), props.Text{Size: 7}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", loan.Amount), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(2, fmt.Sprintf("Q%.2f", loan.LateFee), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(1, fmt.Sprintf("%d", loan.DaysOverdue), props.Text{Size: 7, Align: align.Center}),
				text.NewCol(2, loan.GraceEnds, props.Text{Size: 7}),
			)
		}
	}

	// Approaching due
	if len(data.ApproachingDue) > 0 {
		m.AddRow(10)
		m.AddRow(8, text.NewCol(12, "PRÓXIMOS A VENCER (7 días)", props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}))

		// Table header
		m.AddRow(6,
			text.NewCol(2, "No.", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(4, "Cliente", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(3, "Saldo", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
			text.NewCol(3, "Vence", props.Text{Size: 8, Style: fontstyle.Bold}),
		)

		for _, loan := range data.ApproachingDue {
			m.AddRow(5,
				text.NewCol(2, loan.LoanNumber, props.Text{Size: 7}),
				text.NewCol(4, truncateName(loan.CustomerName, 25), props.Text{Size: 7}),
				text.NewCol(3, fmt.Sprintf("Q%.2f", loan.Amount), props.Text{Size: 7, Align: align.Right}),
				text.NewCol(3, loan.DueDate, props.Text{Size: 7}),
			)
		}
	}

	// Footer
	m.AddRow(15)
	m.AddRow(5, text.NewCol(12, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
		Size:  8,
		Align: align.Right,
	}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// Helper function to truncate names
func truncateName(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
