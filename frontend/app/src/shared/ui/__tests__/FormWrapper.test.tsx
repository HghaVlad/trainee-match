import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { z } from 'zod'
import { FormWrapper } from '../Form/FormWrapper'
import { FormField, FormItem, FormLabel, FormControl, FormMessage } from '../form'
import { Input } from '../input'
import { Button } from '../button'

const testSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
})

describe('FormWrapper', () => {
  it('renders form children correctly', () => {
    const onSubmit = vi.fn()
    render(
      <FormWrapper schema={testSchema} onSubmit={onSubmit}>
        <div data-testid="form-child">Child</div>
      </FormWrapper>
    )
    expect(screen.getByTestId('form-child')).toBeInTheDocument()
  })

  it('validates with zod and shows error', async () => {
    const onSubmit = vi.fn()
    render(
      <FormWrapper schema={testSchema} onSubmit={onSubmit} defaultValues={{ name: '' }}>
        {(form) => (
          <>
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Submit</Button>
          </>
        )}
      </FormWrapper>
    )

    const submitBtn = screen.getByRole('button', { name: 'Submit' })
    fireEvent.click(submitBtn)

    await waitFor(() => {
      expect(screen.getByText('Name must be at least 2 characters')).toBeInTheDocument()
    })
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('submits valid data', async () => {
    const onSubmit = vi.fn()
    render(
      <FormWrapper schema={testSchema} onSubmit={onSubmit} defaultValues={{ name: '' }}>
        {(form) => (
          <>
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Submit</Button>
          </>
        )}
      </FormWrapper>
    )

    const input = screen.getByPlaceholderText('Enter name')
    fireEvent.change(input, { target: { value: 'John' } })

    const submitBtn = screen.getByRole('button', { name: 'Submit' })
    fireEvent.click(submitBtn)

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({ name: 'John' }, expect.anything())
    })
  })
})
