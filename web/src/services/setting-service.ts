import { apiGet, apiPost, apiPut } from '@/lib/api-client'
import type {
  Setting,
  SetSettingInput,
  SetMultipleSettingsInput,
} from '@/types'

export const settingService = {
  // Get all settings
  list: (branchId?: number) =>
    apiGet<Setting[]>('/settings', branchId ? { branch_id: branchId } : undefined),

  // Get a single setting by key
  getByKey: (key: string, branchId?: number) =>
    apiGet<Setting>(`/settings/${key}`, branchId ? { branch_id: branchId } : undefined),

  // Get multiple settings by keys
  getByKeys: (keys: string[], branchId?: number) =>
    apiPost<Setting[]>('/settings/batch', {
      keys,
      branch_id: branchId,
    }),

  // Set a single setting
  set: (input: SetSettingInput) =>
    apiPost<Setting>('/settings', input),

  // Set multiple settings
  setMultiple: (input: SetMultipleSettingsInput) =>
    apiPut<Setting[]>('/settings/batch', input),

  // Get public settings (no auth required)
  getPublic: () =>
    apiGet<Setting[]>('/settings/public'),

  // Reset setting to default
  reset: (key: string, branchId?: number) =>
    apiPost<void>(`/settings/${key}/reset`, {
      branch_id: branchId,
    }),
}
