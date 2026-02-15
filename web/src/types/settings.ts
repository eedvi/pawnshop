// Settings types - mirrors internal/domain/setting.go

export type SettingDataType = 'string' | 'number' | 'boolean' | 'json'

export interface Setting {
  id: number
  key: string
  value: string
  data_type: SettingDataType
  branch_id?: number
  description?: string
  is_public: boolean
  created_at: string
  updated_at: string
}

export interface SetSettingInput {
  key: string
  value: string
  data_type?: SettingDataType
  branch_id?: number
  description?: string
}

export interface SetMultipleSettingsInput {
  settings: SetSettingInput[]
}

// Common setting keys
export const SETTING_KEYS = {
  // Company info
  COMPANY_NAME: 'company.name',
  COMPANY_ADDRESS: 'company.address',
  COMPANY_PHONE: 'company.phone',
  COMPANY_EMAIL: 'company.email',
  COMPANY_LOGO: 'company.logo_url',
  COMPANY_TAX_ID: 'company.tax_id',

  // Loan defaults
  DEFAULT_INTEREST_RATE: 'loan.default_interest_rate',
  DEFAULT_LOAN_TERM_DAYS: 'loan.default_term_days',
  DEFAULT_GRACE_PERIOD: 'loan.default_grace_period',
  LATE_FEE_RATE: 'loan.late_fee_rate',
  MIN_LOAN_AMOUNT: 'loan.min_amount',
  MAX_LOAN_AMOUNT: 'loan.max_amount',

  // Notifications
  EMAIL_ENABLED: 'notification.email_enabled',
  SMS_ENABLED: 'notification.sms_enabled',
  WHATSAPP_ENABLED: 'notification.whatsapp_enabled',
  REMINDER_DAYS_BEFORE: 'notification.reminder_days_before',

  // System
  TIMEZONE: 'system.timezone',
  CURRENCY: 'system.currency',
  DATE_FORMAT: 'system.date_format',
  ALLOW_NEGATIVE_CASH: 'system.allow_negative_cash',

  // Documents
  CONTRACT_TEMPLATE: 'document.contract_template',
  RECEIPT_TEMPLATE: 'document.receipt_template',
  TICKET_FOOTER: 'document.ticket_footer',

  // Loyalty
  LOYALTY_ENABLED: 'loyalty.enabled',
  LOYALTY_POINTS_PER_CURRENCY: 'loyalty.points_per_currency',
  LOYALTY_REDEMPTION_RATE: 'loyalty.redemption_rate',
} as const

// Helper to parse setting value based on data type
export function parseSettingValue(value: string, dataType: SettingDataType): unknown {
  switch (dataType) {
    case 'number':
      return parseFloat(value)
    case 'boolean':
      return value === 'true'
    case 'json':
      try {
        return JSON.parse(value)
      } catch {
        return null
      }
    default:
      return value
  }
}
