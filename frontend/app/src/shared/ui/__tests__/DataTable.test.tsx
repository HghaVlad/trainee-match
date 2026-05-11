import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { createColumnHelper } from '@tanstack/react-table'
import { DataTable } from '../DataTable/DataTable'

interface Person {
  id: number
  name: string
  age: number
}

const data: Person[] = [
  { id: 1, name: 'Alice', age: 24 },
  { id: 2, name: 'Bob', age: 30 },
]

const columnHelper = createColumnHelper<Person>()
const columns = [
  columnHelper.accessor('id', { header: 'ID' }),
  columnHelper.accessor('name', { header: 'Name' }),
  columnHelper.accessor('age', { header: 'Age' }),
]

describe('DataTable', () => {
  it('renders correct number of columns and rows', () => {
    render(<DataTable columns={columns} data={data} />)
    
    const headers = screen.getAllByRole('columnheader')
    expect(headers).toHaveLength(3)
    expect(headers[0]).toHaveTextContent('ID')
    expect(headers[1]).toHaveTextContent('Name')
    expect(headers[2]).toHaveTextContent('Age')

    const rows = screen.getAllByRole('row')
    expect(rows).toHaveLength(3) 

    const cells = screen.getAllByRole('cell')
    expect(cells).toHaveLength(6)
    expect(cells[0]).toHaveTextContent('1')
    expect(cells[1]).toHaveTextContent('Alice')
    expect(cells[2]).toHaveTextContent('24')
  })

  it('renders empty state when no data provided', () => {
    render(<DataTable columns={columns} data={[]} />)
    
    const headers = screen.getAllByRole('columnheader')
    expect(headers).toHaveLength(3)

    const cell = screen.getByRole('cell')
    expect(cell).toHaveTextContent('No results.')
    expect(cell).toHaveAttribute('colSpan', '3')
  })
})
