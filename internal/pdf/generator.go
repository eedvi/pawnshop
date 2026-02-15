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
