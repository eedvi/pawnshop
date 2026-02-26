package pdf

import (
	"fmt"
	"strings"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"pawnshop/internal/domain"
)

// ThermalTicketGenerator generates 80mm thermal printer tickets
type ThermalTicketGenerator struct {
	companyName    string
	companyAddress string
	companyPhone   string
	companyRFC     string // Tax ID for Mexico/Latin America
}

// NewThermalTicketGenerator creates a new thermal ticket generator
func NewThermalTicketGenerator(companyName, companyAddress, companyPhone, companyRFC string) *ThermalTicketGenerator {
	return &ThermalTicketGenerator{
		companyName:    companyName,
		companyAddress: companyAddress,
		companyPhone:   companyPhone,
		companyRFC:     companyRFC,
	}
}

// createThermalConfig creates the configuration for 80mm thermal paper
func createThermalConfig(height float64) core.Maroto {
	cfg := config.NewBuilder().
		WithDimensions(72, height). // 72mm width (printable area of 80mm paper)
		WithLeftMargin(2).
		WithTopMargin(2).
		WithRightMargin(2).
		WithDefaultFont(&props.Font{
			Size:   8,
			Family: "courier",
		}).
		Build()
	return maroto.New(cfg)
}

// GenerateLoanTicket generates a thermal ticket for a loan
func (g *ThermalTicketGenerator) GenerateLoanTicket(loan *domain.Loan, customer *domain.Customer, item *domain.Item) ([]byte, error) {
	m := createThermalConfig(180)

	// Header
	g.addTicketHeader(m, "BOLETA DE EMPEÑO")

	// Contract number and date
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("No: %s", loan.LoanNumber), props.Text{
		Size:  9,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("Fecha: %s", loan.StartDate.Format("02/01/2006 15:04")), props.Text{
		Size:  7,
		Align: align.Center,
	}))

	g.addSeparator(m)

	// Customer info
	m.AddRow(3, text.NewCol(12, "CLIENTE:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("%s %s", customer.FirstName, customer.LastName), props.Text{Size: 7}))
	if customer.IdentityNumber != "" {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("%s: %s", customer.IdentityType, customer.IdentityNumber), props.Text{Size: 6}))
	}

	g.addSeparator(m)

	// Item info
	m.AddRow(3, text.NewCol(12, "ARTICULO:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3, text.NewCol(12, truncateString(item.Name, 40), props.Text{Size: 7}))
	if item.Brand != nil {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("Marca: %s", *item.Brand), props.Text{Size: 6}))
	}
	if item.SerialNumber != nil {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("Serie: %s", *item.SerialNumber), props.Text{Size: 6}))
	}
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Cond: %s", item.Condition), props.Text{Size: 6}))

	g.addSeparator(m)

	// Loan details
	m.AddRow(3, text.NewCol(12, "DETALLE DEL PRESTAMO:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3,
		text.NewCol(6, "Prestamo:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.LoanAmount), props.Text{Size: 7, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, "Interes:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.InterestAmount), props.Text{Size: 7, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, "Tasa:", props.Text{Size: 6}),
		text.NewCol(6, fmt.Sprintf("%.1f%% mensual", loan.InterestRate), props.Text{Size: 6, Align: align.Right}),
	)

	g.addSeparator(m)

	m.AddRow(4,
		text.NewCol(6, "TOTAL:", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.TotalAmount), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)

	g.addSeparator(m)

	// Due date
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("VENCE: %s", loan.DueDate.Format("02/01/2006")), props.Text{
		Size:  8,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Plazo: %d dias | Gracia: %d dias", loan.LoanTermDays, loan.GracePeriodDays), props.Text{
		Size:  6,
		Align: align.Center,
	}))

	g.addSeparator(m)

	// Footer notice
	m.AddRow(3, text.NewCol(12, "Conserve este ticket para", props.Text{Size: 6, Align: align.Center}))
	m.AddRow(3, text.NewCol(12, "recoger su articulo", props.Text{Size: 6, Align: align.Center}))

	m.AddRow(5)
	m.AddRow(3, text.NewCol(12, strings.Repeat("*", 32), props.Text{Size: 7, Align: align.Center}))

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thermal ticket: %w", err)
	}

	return document.GetBytes(), nil
}

// GeneratePaymentTicket generates a thermal ticket for a payment
func (g *ThermalTicketGenerator) GeneratePaymentTicket(payment *domain.Payment, loan *domain.Loan, customer *domain.Customer) ([]byte, error) {
	m := createThermalConfig(140)

	// Header
	g.addTicketHeader(m, "RECIBO DE PAGO")

	// Receipt number and date
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("Recibo: %s", payment.PaymentNumber), props.Text{
		Size:  8,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("%s", payment.PaymentDate.Format("02/01/2006 15:04")), props.Text{
		Size:  7,
		Align: align.Center,
	}))

	g.addSeparator(m)

	// Customer and loan info
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Cliente: %s %s", customer.FirstName, customer.LastName), props.Text{Size: 7}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Contrato: %s", loan.LoanNumber), props.Text{Size: 7}))

	g.addSeparator(m)

	// Payment breakdown
	m.AddRow(3, text.NewCol(12, "DETALLE:", props.Text{Size: 7, Style: fontstyle.Bold}))

	if payment.PrincipalAmount > 0 {
		m.AddRow(3,
			text.NewCol(6, "Capital:", props.Text{Size: 7}),
			text.NewCol(6, fmt.Sprintf("$%.2f", payment.PrincipalAmount), props.Text{Size: 7, Align: align.Right}),
		)
	}
	if payment.InterestAmount > 0 {
		m.AddRow(3,
			text.NewCol(6, "Interes:", props.Text{Size: 7}),
			text.NewCol(6, fmt.Sprintf("$%.2f", payment.InterestAmount), props.Text{Size: 7, Align: align.Right}),
		)
	}
	if payment.LateFeeAmount > 0 {
		m.AddRow(3,
			text.NewCol(6, "Mora:", props.Text{Size: 7}),
			text.NewCol(6, fmt.Sprintf("$%.2f", payment.LateFeeAmount), props.Text{Size: 7, Align: align.Right}),
		)
	}

	g.addSeparator(m)

	// Total paid
	m.AddRow(4,
		text.NewCol(6, "PAGADO:", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", payment.Amount), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, "Metodo:", props.Text{Size: 6}),
		text.NewCol(6, string(payment.PaymentMethod), props.Text{Size: 6, Align: align.Right}),
	)

	g.addSeparator(m)

	// Remaining balance
	balance := loan.RemainingBalance()
	m.AddRow(4,
		text.NewCol(6, "SALDO:", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", balance), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)

	if loan.Status == domain.LoanStatusPaid {
		m.AddRow(5)
		m.AddRow(4, text.NewCol(12, "*** LIQUIDADO ***", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
		m.AddRow(3, text.NewCol(12, "Puede recoger su articulo", props.Text{Size: 6, Align: align.Center}))
	}

	g.addTicketFooter(m)

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thermal ticket: %w", err)
	}

	return document.GetBytes(), nil
}

// GenerateSaleTicket generates a thermal ticket for a sale
func (g *ThermalTicketGenerator) GenerateSaleTicket(sale *domain.Sale, item *domain.Item, customer *domain.Customer) ([]byte, error) {
	m := createThermalConfig(130)

	// Header
	g.addTicketHeader(m, "TICKET DE VENTA")

	// Sale number and date
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("Venta: %s", sale.SaleNumber), props.Text{
		Size:  8,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(3, text.NewCol(12, sale.SaleDate.Format("02/01/2006 15:04"), props.Text{
		Size:  7,
		Align: align.Center,
	}))

	g.addSeparator(m)

	// Customer (if available)
	if customer != nil {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("Cliente: %s %s", customer.FirstName, customer.LastName), props.Text{Size: 7}))
		g.addSeparator(m)
	}

	// Item details
	m.AddRow(3, text.NewCol(12, "ARTICULO:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3, text.NewCol(12, truncateString(item.Name, 40), props.Text{Size: 7}))
	if item.Brand != nil {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("Marca: %s", *item.Brand), props.Text{Size: 6}))
	}

	g.addSeparator(m)

	// Price breakdown
	m.AddRow(3,
		text.NewCol(6, "Precio:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", sale.SalePrice), props.Text{Size: 7, Align: align.Right}),
	)

	if sale.DiscountAmount > 0 {
		m.AddRow(3,
			text.NewCol(6, "Descuento:", props.Text{Size: 7}),
			text.NewCol(6, fmt.Sprintf("-$%.2f", sale.DiscountAmount), props.Text{Size: 7, Align: align.Right}),
		)
	}

	g.addSeparator(m)

	// Total
	m.AddRow(4,
		text.NewCol(6, "TOTAL:", props.Text{Size: 9, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", sale.FinalPrice), props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, "Pago:", props.Text{Size: 6}),
		text.NewCol(6, string(sale.PaymentMethod), props.Text{Size: 6, Align: align.Right}),
	)

	g.addTicketFooter(m)

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thermal ticket: %w", err)
	}

	return document.GetBytes(), nil
}

// GenerateCashCutTicket generates a thermal ticket for cash register cut (X or Z)
func (g *ThermalTicketGenerator) GenerateCashCutTicket(cutType string, session *domain.CashSession, movements []*domain.CashMovement) ([]byte, error) {
	m := createThermalConfig(200)

	// Header
	title := "CORTE X"
	if cutType == "Z" {
		title = "CORTE Z (CIERRE)"
	}
	g.addTicketHeader(m, title)

	m.AddRow(3, text.NewCol(12, time.Now().Format("02/01/2006 15:04"), props.Text{
		Size:  7,
		Align: align.Center,
	}))

	g.addSeparator(m)

	// Session info
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Sesion: %d", session.ID), props.Text{Size: 7}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Caja: %d", session.CashRegisterID), props.Text{Size: 7}))
	m.AddRow(3, text.NewCol(12, fmt.Sprintf("Apertura: %s", session.OpenedAt.Format("02/01/06 15:04")), props.Text{Size: 6}))

	g.addSeparator(m)

	// Summary
	m.AddRow(3, text.NewCol(12, "RESUMEN:", props.Text{Size: 7, Style: fontstyle.Bold}))

	// Calculate totals from movements
	var totalIncome, totalExpense float64
	var incomeCount, expenseCount int

	for _, mov := range movements {
		if mov.MovementType == domain.CashMovementTypeIncome {
			totalIncome += mov.Amount
			incomeCount++
		} else {
			totalExpense += mov.Amount
			expenseCount++
		}
	}

	m.AddRow(3,
		text.NewCol(6, "Saldo Inicial:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", session.OpeningAmount), props.Text{Size: 7, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, fmt.Sprintf("Ingresos (%d):", incomeCount), props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", totalIncome), props.Text{Size: 7, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, fmt.Sprintf("Egresos (%d):", expenseCount), props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", totalExpense), props.Text{Size: 7, Align: align.Right}),
	)

	g.addSeparator(m)

	// Calculated balance
	calculatedBalance := session.OpeningAmount + totalIncome - totalExpense
	m.AddRow(4,
		text.NewCol(6, "ESPERADO:", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", calculatedBalance), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)

	if session.ClosingAmount != nil {
		m.AddRow(4,
			text.NewCol(6, "CONTADO:", props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("$%.2f", *session.ClosingAmount), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
		)

		diff := *session.ClosingAmount - calculatedBalance
		diffLabel := "DIFERENCIA:"
		if diff != 0 {
			diffLabel = "DIFERENCIA!:"
		}
		m.AddRow(4,
			text.NewCol(6, diffLabel, props.Text{Size: 8, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("$%.2f", diff), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
		)
	}

	g.addSeparator(m)

	// Movement details (last 10)
	m.AddRow(3, text.NewCol(12, "MOVIMIENTOS:", props.Text{Size: 7, Style: fontstyle.Bold}))

	displayMovements := movements
	if len(movements) > 10 {
		displayMovements = movements[len(movements)-10:]
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("(ultimos 10 de %d)", len(movements)), props.Text{Size: 6}))
	}

	for _, mov := range displayMovements {
		sign := "+"
		if mov.MovementType == domain.CashMovementTypeExpense {
			sign = "-"
		}
		desc := truncateString(mov.Description, 20)
		m.AddRow(3,
			text.NewCol(8, desc, props.Text{Size: 6}),
			text.NewCol(4, fmt.Sprintf("%s$%.2f", sign, mov.Amount), props.Text{Size: 6, Align: align.Right}),
		)
	}

	g.addTicketFooter(m)

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thermal ticket: %w", err)
	}

	return document.GetBytes(), nil
}

// addTicketHeader adds the header to a thermal ticket
func (g *ThermalTicketGenerator) addTicketHeader(m core.Maroto, title string) {
	m.AddRow(5, text.NewCol(12, g.companyName, props.Text{
		Size:  10,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	if g.companyAddress != "" {
		m.AddRow(3, text.NewCol(12, truncateString(g.companyAddress, 35), props.Text{
			Size:  6,
			Align: align.Center,
		}))
	}

	if g.companyPhone != "" {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("Tel: %s", g.companyPhone), props.Text{
			Size:  6,
			Align: align.Center,
		}))
	}

	if g.companyRFC != "" {
		m.AddRow(3, text.NewCol(12, fmt.Sprintf("RFC: %s", g.companyRFC), props.Text{
			Size:  6,
			Align: align.Center,
		}))
	}

	m.AddRow(2)
	m.AddRow(4, text.NewCol(12, title, props.Text{
		Size:  9,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(2)
}

// addTicketFooter adds the footer to a thermal ticket
func (g *ThermalTicketGenerator) addTicketFooter(m core.Maroto) {
	m.AddRow(4)
	m.AddRow(3, text.NewCol(12, "Gracias por su preferencia", props.Text{
		Size:  7,
		Align: align.Center,
	}))
	m.AddRow(4)
	m.AddRow(3, text.NewCol(12, strings.Repeat("*", 32), props.Text{
		Size:  7,
		Align: align.Center,
	}))
}

// addSeparator adds a dashed line separator
func (g *ThermalTicketGenerator) addSeparator(m core.Maroto) {
	m.AddRow(3, text.NewCol(12, strings.Repeat("-", 36), props.Text{
		Size:  6,
		Align: align.Center,
	}))
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ThermalPaperSize represents thermal paper width
type ThermalPaperSize int

const (
	ThermalPaper58mm ThermalPaperSize = 58
	ThermalPaper80mm ThermalPaperSize = 80
)

// GenerateLoanTicketWithSize generates a loan ticket for a specific paper size
func (g *ThermalTicketGenerator) GenerateLoanTicketWithSize(loan *domain.Loan, customer *domain.Customer, item *domain.Item, paperSize ThermalPaperSize) ([]byte, error) {
	var width, height float64
	var charLimit int

	switch paperSize {
	case ThermalPaper58mm:
		width = 48  // 48mm printable width for 58mm paper
		height = 160
		charLimit = 28
	default:
		width = 72  // 72mm printable width for 80mm paper
		height = 180
		charLimit = 40
	}

	cfg := config.NewBuilder().
		WithDimensions(width, height).
		WithLeftMargin(1).
		WithTopMargin(2).
		WithRightMargin(1).
		WithDefaultFont(&props.Font{
			Size:   7,
			Family: "courier",
		}).
		Build()

	m := maroto.New(cfg)

	g.addTicketHeader(m, "BOLETA DE EMPEÑO")

	m.AddRow(4, text.NewCol(12, fmt.Sprintf("No: %s", loan.LoanNumber), props.Text{
		Size:  9,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(4, text.NewCol(12, fmt.Sprintf("Fecha: %s", loan.StartDate.Format("02/01/2006 15:04")), props.Text{
		Size:  7,
		Align: align.Center,
	}))

	g.addSeparator(m)

	m.AddRow(3, text.NewCol(12, "CLIENTE:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3, text.NewCol(12, truncateString(fmt.Sprintf("%s %s", customer.FirstName, customer.LastName), charLimit), props.Text{Size: 7}))

	g.addSeparator(m)

	m.AddRow(3, text.NewCol(12, "ARTICULO:", props.Text{Size: 7, Style: fontstyle.Bold}))
	m.AddRow(3, text.NewCol(12, truncateString(item.Name, charLimit), props.Text{Size: 7}))

	g.addSeparator(m)

	m.AddRow(3,
		text.NewCol(6, "Prestamo:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.LoanAmount), props.Text{Size: 7, Align: align.Right}),
	)
	m.AddRow(3,
		text.NewCol(6, "Interes:", props.Text{Size: 7}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.InterestAmount), props.Text{Size: 7, Align: align.Right}),
	)

	g.addSeparator(m)

	m.AddRow(4,
		text.NewCol(6, "TOTAL:", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, fmt.Sprintf("$%.2f", loan.TotalAmount), props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)

	m.AddRow(4, text.NewCol(12, fmt.Sprintf("VENCE: %s", loan.DueDate.Format("02/01/2006")), props.Text{
		Size:  8,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	g.addTicketFooter(m)

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thermal ticket: %w", err)
	}

	return document.GetBytes(), nil
}
