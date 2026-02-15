// Report types - matches report service outputs

export interface DashboardStats {
  // Loan stats
  active_loans: number
  active_loans_amount: number
  overdue_loans: number
  overdue_loans_amount: number
  loans_due_today: number
  loans_due_this_week: number

  // Payment stats
  payments_today: number
  payments_today_amount: number
  payments_this_month: number
  payments_this_month_amount: number

  // Sales stats
  sales_today: number
  sales_today_amount: number
  sales_this_month: number
  sales_this_month_amount: number

  // Customer stats
  total_customers: number
  new_customers_this_month: number
  blocked_customers: number

  // Item stats
  items_available: number
  items_for_sale: number
  items_pawned: number

  // Cash stats
  cash_sessions_open: number
  cash_balance_today: number

  // Trends (percentage change from previous period)
  loans_trend?: number
  payments_trend?: number
  sales_trend?: number
  customers_trend?: number
}

export interface LoanReportItem {
  loan_id: number
  loan_number: string
  customer_name: string
  item_name: string
  loan_amount: number
  interest_rate: number
  principal_remaining: number
  interest_remaining: number
  total_remaining: number
  status: string
  start_date: string
  due_date: string
  days_overdue: number
  branch_name: string
}

export interface LoanReport {
  items: LoanReportItem[]
  summary: {
    total_loans: number
    total_principal: number
    total_interest: number
    total_remaining: number
    by_status: {
      status: string
      count: number
      amount: number
    }[]
  }
}

export interface PaymentReportItem {
  payment_id: number
  payment_number: string
  loan_number: string
  customer_name: string
  amount: number
  principal_amount: number
  interest_amount: number
  late_fee_amount: number
  payment_method: string
  payment_date: string
  status: string
  branch_name: string
}

export interface PaymentReport {
  items: PaymentReportItem[]
  summary: {
    total_payments: number
    total_amount: number
    total_principal: number
    total_interest: number
    total_late_fees: number
    by_method: {
      method: string
      count: number
      amount: number
    }[]
    by_status: {
      status: string
      count: number
      amount: number
    }[]
  }
}

export interface SalesReportItem {
  sale_id: number
  sale_number: string
  item_name: string
  customer_name?: string
  sale_price: number
  discount_amount: number
  final_price: number
  payment_method: string
  sale_date: string
  status: string
  branch_name: string
}

export interface SalesReport {
  items: SalesReportItem[]
  summary: {
    total_sales: number
    gross_amount: number
    total_discounts: number
    net_amount: number
    total_refunds: number
    refund_amount: number
    by_method: {
      method: string
      count: number
      amount: number
    }[]
  }
}

export interface OverdueReportItem {
  loan_id: number
  loan_number: string
  customer_id: number
  customer_name: string
  customer_phone: string
  item_name: string
  loan_amount: number
  total_remaining: number
  due_date: string
  days_overdue: number
  late_fee_amount: number
  grace_period_ends: string
  branch_name: string
}

export interface OverdueReport {
  items: OverdueReportItem[]
  summary: {
    total_overdue_loans: number
    total_overdue_amount: number
    total_late_fees: number
    by_days_overdue: {
      range: string
      count: number
      amount: number
    }[]
  }
}

export interface ReportFilters {
  branch_id?: number
  date_from?: string
  date_to?: string
  status?: string
  customer_id?: number
}
