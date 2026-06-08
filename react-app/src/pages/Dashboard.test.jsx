import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, useNavigate } from 'react-router-dom'
import Dashboard from './Dashboard'

const mockNavigate = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

vi.mock('../components/ParabolicBackground', () => {
  return {
    default: () => null,
  }
})

vi.mock('../components/TapeDetail', () => ({
  default: ({ tape, onBack, onRent, error }) => (
    <div>
      <p>{tape.title}</p>
      <button onClick={onBack}>Back</button>
      <button onClick={onRent}>Rent</button>
      {error && <p>{error}</p>}
    </div>
  ),
}))

vi.mock('../components/MyRentals', () => ({
  default: ({ rentals, onRentalClick, onBack }) => (
    <div>
      {rentals.map(r => (
        <p key={r.public_id} onClick={() => onRentalClick(r)}>{r.tape_title}</p>
      ))}
      <button onClick={onBack}>Back</button>
    </div>
  ),
}))

vi.mock('../components/RentalDetail', () => ({
  default: ({ rental, onBack, onReturnRent, error }) => (
    <div>
      <p>{rental.tape_title}</p>
      <button onClick={onBack}>Back</button>
      <button onClick={onReturnRent}>Return</button>
      {error && <p>{error}</p>}
    </div>
  ),
}))

const fakeTapes = [
  { public_id: '1', title: 'Blade Runner', director: 'Ridley Scott' },
  { public_id: '2', title: 'Alien', director: 'Ridley Scott' },
]

const fakeRentals = [
  { public_id: 'r1', tape_title: 'Blade Runner', username: 'admin' },
]

function renderDashboard() {
  return render(
    <MemoryRouter>
      <Dashboard />
    </MemoryRouter>
  )
}

beforeEach(() => {
  vi.resetAllMocks()
  global.fetch = vi.fn()
  localStorage.clear()
})

describe('Dashboard', () => {

  it('redirects to /login if no token in localStorage', () => {
    renderDashboard()
    expect(mockNavigate).toHaveBeenCalledWith('/login')
  })

  it('shows loading state while fetching tapes', () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn(() => new Promise(() => { }))

    renderDashboard()
    expect(screen.getByText('Loading tapes...')).toBeInTheDocument()
  })

  it('displays username from localStorage', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => fakeTapes,
    })

    renderDashboard()
    expect(screen.getByText('admin')).toBeInTheDocument()
  })

  it('fetches and displays tapes on mount', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => fakeTapes,
    })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
      expect(screen.getByText('Alien')).toBeInTheDocument()
    })
  })

  it('displays error when fetch fails', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'Server error' }),
    })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })

  it('logs out, clears localStorage and navigates to /login', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => fakeTapes,
    })

    renderDashboard()
    await userEvent.click(screen.getByText('Logout'))

    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('username')).toBeNull()
    expect(mockNavigate).toHaveBeenCalledWith('/login')
  })

  it('shows tape detail when a tape is clicked', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => fakeTapes,
    })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))

    expect(screen.getByText('Rent')).toBeInTheDocument()
    expect(screen.getByText('Back')).toBeInTheDocument()
  })

  it('navigates back to catalog form tape detail', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => fakeTapes,
    })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Back'))

    await waitFor(() => {
      expect(screen.getByText('VHS Catalog')).toBeInTheDocument()
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })
  })

  it('fetches and displays rentals when My Rentals is clicked', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes, })
      .mockResolvedValueOnce({ ok: true, json: async () => fakeRentals, })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('My Rentals'))

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: 'My Rentals' })).toBeInTheDocument()
    })
  })

  it('rents a tape successfully', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockResolvedValueOnce({ ok: true, json: async () => ({}) })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Rent'))

    expect(global.fetch).toHaveBeenCalledWith('api/rentals/1', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer abc123' },
    })
  })

  it('shows error when renting fails', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockResolvedValueOnce({ ok: false, json: async () => ({ error: 'Already rented' }) })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Rent'))

    await waitFor(() => {
      expect(screen.getByText('Already rented')).toBeInTheDocument()
    })
  })

  it('returns a rental successfully', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockResolvedValueOnce({ ok: true, json: async () => fakeRentals })
      .mockResolvedValueOnce({ ok: true, json: async () => ({}) })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByRole('button', { name: 'My Rentals' }))

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Return'))

    expect(global.fetch).toHaveBeenCalledWith('api/rentals/r1', {
      method: 'PATCH',
      headers: { 'Authorization': 'Bearer abc123' },
    })
  })

  it('shows error when returning a rental fails', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockResolvedValueOnce({ ok: true, json: async () => fakeRentals })
      .mockResolvedValueOnce({ ok: false, json: async () => ({ error: 'Return failed' }) })

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByRole('button', { name: 'My Rentals' }))

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Return'))

    await waitFor(() => {
      expect(screen.getByText('Return failed')).toBeInTheDocument()
    })
  })

  it('shows network error when renting fails due to network', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockRejectedValueOnce(new Error('Network error'))

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Rent'))

    await waitFor(() => {
      expect(screen.getByText(/Unable to connect to the server/)).toBeInTheDocument()
    })
  })

  it('shows network error when returning a rental fails due to network', async () => {
    localStorage.setItem('token', 'abc123')
    localStorage.setItem('username', 'admin')
    global.fetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => fakeTapes })
      .mockResolvedValueOnce({ ok: true, json: async () => fakeRentals })
      .mockRejectedValueOnce(new Error('Network error'))

    renderDashboard()

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByRole('button', { name: 'My Rentals' }))

    await waitFor(() => {
      expect(screen.getByText('Blade Runner')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByText('Blade Runner'))
    await userEvent.click(screen.getByText('Return'))

    await waitFor(() => {
      expect(screen.getByText(/Unable to connect to the server/)).toBeInTheDocument()
    })
  })

})










