import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { settingService } from '@/services/setting-service'
import type { SetSettingInput, SetMultipleSettingsInput } from '@/types'

export const settingKeys = {
  all: ['settings'] as const,
  lists: () => [...settingKeys.all, 'list'] as const,
  list: (branchId?: number) => [...settingKeys.lists(), branchId] as const,
  byKey: (key: string, branchId?: number) => [...settingKeys.all, 'key', key, branchId] as const,
  byKeys: (keys: string[], branchId?: number) => [...settingKeys.all, 'keys', keys, branchId] as const,
  public: () => [...settingKeys.all, 'public'] as const,
}

export function useSettings(branchId?: number) {
  return useQuery({
    queryKey: settingKeys.list(branchId),
    queryFn: () => settingService.list(branchId),
  })
}

export function useSetting(key: string, branchId?: number) {
  return useQuery({
    queryKey: settingKeys.byKey(key, branchId),
    queryFn: () => settingService.getByKey(key, branchId),
    enabled: !!key,
  })
}

export function useSettingsByKeys(keys: string[], branchId?: number) {
  return useQuery({
    queryKey: settingKeys.byKeys(keys, branchId),
    queryFn: () => settingService.getByKeys(keys, branchId),
    enabled: keys.length > 0,
  })
}

export function usePublicSettings() {
  return useQuery({
    queryKey: settingKeys.public(),
    queryFn: () => settingService.getPublic(),
  })
}

export function useSetSetting() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: SetSettingInput) => settingService.set(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: settingKeys.lists() })
      queryClient.invalidateQueries({
        queryKey: settingKeys.byKey(variables.key, variables.branch_id),
      })
    },
  })
}

export function useSetMultipleSettings() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: SetMultipleSettingsInput) => settingService.setMultiple(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingKeys.all })
    },
  })
}

export function useResetSetting() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ key, branchId }: { key: string; branchId?: number }) =>
      settingService.reset(key, branchId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingKeys.all })
    },
  })
}
