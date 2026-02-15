import { useState } from 'react'
import { Plus, ChevronRight, ChevronDown, Pencil, Trash2, FolderTree, Percent } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { StatusBadge } from '@/components/common/status-badge'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useCategoryTree, useCreateCategory, useUpdateCategory, useDeleteCategory } from '@/hooks/use-categories'
import { useConfirm } from '@/hooks'
import { Category } from '@/types'
import { formatPercent } from '@/lib/format'
import { cn } from '@/lib/utils'
import { CategoryForm } from './category-form'
import { CategoryFormValues } from './schemas'

interface CategoryTreeNodeProps {
  category: Category
  level: number
  expandedIds: Set<number>
  onToggle: (id: number) => void
  onEdit: (category: Category) => void
  onDelete: (category: Category) => void
  onAddChild: (parentId: number) => void
}

function CategoryTreeNode({
  category,
  level,
  expandedIds,
  onToggle,
  onEdit,
  onDelete,
  onAddChild,
}: CategoryTreeNodeProps) {
  const hasChildren = category.children && category.children.length > 0
  const isExpanded = expandedIds.has(category.id)

  return (
    <div>
      <div
        className={cn(
          'flex items-center gap-2 rounded-lg border p-3 transition-colors hover:bg-muted/50',
          !category.is_active && 'opacity-60'
        )}
        style={{ marginLeft: `${level * 24}px` }}
      >
        {/* Expand/Collapse Button */}
        <button
          type="button"
          className={cn(
            'flex h-6 w-6 items-center justify-center rounded hover:bg-muted',
            !hasChildren && 'invisible'
          )}
          onClick={() => onToggle(category.id)}
        >
          {isExpanded ? (
            <ChevronDown className="h-4 w-4" />
          ) : (
            <ChevronRight className="h-4 w-4" />
          )}
        </button>

        {/* Icon */}
        <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10">
          <FolderTree className="h-4 w-4 text-primary" />
        </div>

        {/* Name and Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-medium truncate">{category.name}</span>
            <StatusBadge status={category.is_active ? 'active' : 'inactive'} size="sm" />
          </div>
          {category.description && (
            <p className="text-sm text-muted-foreground truncate">{category.description}</p>
          )}
        </div>

        {/* Rate and LTV */}
        <div className="hidden sm:flex items-center gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Percent className="h-3 w-3" />
            <span>{formatPercent(category.default_interest_rate)}</span>
          </div>
          <div>
            <span>LTV: {formatPercent(category.loan_to_value_ratio)}</span>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => onAddChild(category.id)}
            title="Agregar subcategoría"
          >
            <Plus className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => onEdit(category)}
            title="Editar"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 text-destructive hover:text-destructive"
            onClick={() => onDelete(category)}
            title="Eliminar"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Children */}
      {hasChildren && isExpanded && (
        <div className="mt-2 space-y-2">
          {category.children!.map((child) => (
            <CategoryTreeNode
              key={child.id}
              category={child}
              level={level + 1}
              expandedIds={expandedIds}
              onToggle={onToggle}
              onEdit={onEdit}
              onDelete={onDelete}
              onAddChild={onAddChild}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export default function CategoryPage() {
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set())
  const [dialogOpen, setDialogOpen] = useState(false)
  const [selectedCategory, setSelectedCategory] = useState<Category | undefined>()
  const [parentIdForNew, setParentIdForNew] = useState<number | undefined>()

  const { data: categories, isLoading } = useCategoryTree()
  const createMutation = useCreateCategory()
  const updateMutation = useUpdateCategory()
  const deleteMutation = useDeleteCategory()
  const confirmDelete = useConfirm()

  const handleToggle = (id: number) => {
    setExpandedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const handleExpandAll = () => {
    if (!categories) return
    const getAllIds = (cats: Category[]): number[] => {
      const ids: number[] = []
      for (const cat of cats) {
        ids.push(cat.id)
        if (cat.children) {
          ids.push(...getAllIds(cat.children))
        }
      }
      return ids
    }
    setExpandedIds(new Set(getAllIds(categories)))
  }

  const handleCollapseAll = () => {
    setExpandedIds(new Set())
  }

  const handleCreate = () => {
    setSelectedCategory(undefined)
    setParentIdForNew(undefined)
    setDialogOpen(true)
  }

  const handleAddChild = (parentId: number) => {
    setSelectedCategory(undefined)
    setParentIdForNew(parentId)
    setDialogOpen(true)
  }

  const handleEdit = (category: Category) => {
    setSelectedCategory(category)
    setParentIdForNew(undefined)
    setDialogOpen(true)
  }

  const handleDelete = async (category: Category) => {
    const hasChildren = category.children && category.children.length > 0
    const confirmed = await confirmDelete.confirm({
      title: 'Eliminar Categoría',
      description: hasChildren
        ? `"${category.name}" tiene subcategorías. ¿Estás seguro de eliminarla? Las subcategorías también serán eliminadas.`
        : `¿Estás seguro de eliminar "${category.name}"?`,
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      deleteMutation.mutate(category.id)
    }
  }

  const handleSubmit = async (values: CategoryFormValues) => {
    try {
      if (selectedCategory) {
        await updateMutation.mutateAsync({
          id: selectedCategory.id,
          input: {
            name: values.name,
            description: values.description || undefined,
            parent_id: values.parent_id ?? undefined,
            icon: values.icon || undefined,
            default_interest_rate: values.default_interest_rate,
            min_loan_amount: values.min_loan_amount ?? undefined,
            max_loan_amount: values.max_loan_amount ?? undefined,
            loan_to_value_ratio: values.loan_to_value_ratio,
            sort_order: values.sort_order,
            is_active: values.is_active,
          },
        })
      } else {
        await createMutation.mutateAsync({
          name: values.name,
          description: values.description || undefined,
          parent_id: parentIdForNew ?? values.parent_id ?? undefined,
          icon: values.icon || undefined,
          default_interest_rate: values.default_interest_rate,
          min_loan_amount: values.min_loan_amount ?? undefined,
          max_loan_amount: values.max_loan_amount ?? undefined,
          loan_to_value_ratio: values.loan_to_value_ratio,
          sort_order: values.sort_order,
        })
      }
      setDialogOpen(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <div>
      <PageHeader
        title="Categorías"
        description="Gestión de categorías de artículos"
        actions={
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={handleExpandAll}>
              Expandir todo
            </Button>
            <Button variant="outline" size="sm" onClick={handleCollapseAll}>
              Colapsar todo
            </Button>
            <Button onClick={handleCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Nueva Categoría
            </Button>
          </div>
        }
      />

      <div className="rounded-lg border bg-card p-6">
        {isLoading ? (
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex items-center gap-4 rounded-lg border p-3">
                <Skeleton className="h-6 w-6 rounded" />
                <Skeleton className="h-8 w-8 rounded-md" />
                <div className="flex-1">
                  <Skeleton className="h-5 w-32" />
                  <Skeleton className="mt-1 h-4 w-48" />
                </div>
              </div>
            ))}
          </div>
        ) : categories && categories.length > 0 ? (
          <div className="space-y-2">
            {categories.map((category) => (
              <CategoryTreeNode
                key={category.id}
                category={category}
                level={0}
                expandedIds={expandedIds}
                onToggle={handleToggle}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onAddChild={handleAddChild}
              />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <FolderTree className="h-12 w-12 text-muted-foreground" />
            <h3 className="mt-4 text-lg font-semibold">No hay categorías</h3>
            <p className="mt-2 text-muted-foreground">
              Crea tu primera categoría para organizar los artículos.
            </p>
            <Button onClick={handleCreate} className="mt-4">
              <Plus className="mr-2 h-4 w-4" />
              Nueva Categoría
            </Button>
          </div>
        )}
      </div>

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>
              {selectedCategory ? `Editar: ${selectedCategory.name}` : 'Nueva Categoría'}
            </DialogTitle>
          </DialogHeader>
          <CategoryForm
            category={selectedCategory}
            categories={categories ?? []}
            onSubmit={handleSubmit}
            onCancel={() => setDialogOpen(false)}
            isLoading={createMutation.isPending || updateMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      <ConfirmDialog {...confirmDelete} />
    </div>
  )
}
