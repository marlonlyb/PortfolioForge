import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { AdminUserFormPage } from './AdminUserFormPage';
import { AdminUserListPage } from './AdminUserListPage';
import type { AdminUserDetail, AdminUserSummary } from '../../shared/types/admin-user';
import { deleteAdminUser, fetchAdminUserById, fetchAdminUsers, updateAdminUser } from './api';
import { useSession, type SessionUser } from '../../app/providers/SessionProvider';

vi.mock('./api', async () => {
	const actual = await vi.importActual<typeof import('./api')>('./api');
	return {
		...actual,
		deleteAdminUser: vi.fn(),
		fetchAdminUserById: vi.fn(),
		fetchAdminUsers: vi.fn(),
		updateAdminUser: vi.fn(),
	};
});

vi.mock('../../app/providers/SessionProvider', async () => {
	const actual = await vi.importActual<typeof import('../../app/providers/SessionProvider')>('../../app/providers/SessionProvider');
	return {
		...actual,
		useSession: vi.fn(),
	};
});

const mockedDeleteAdminUser = vi.mocked(deleteAdminUser);
const mockedFetchAdminUserById = vi.mocked(fetchAdminUserById);
const mockedFetchAdminUsers = vi.mocked(fetchAdminUsers);
const mockedUpdateAdminUser = vi.mocked(updateAdminUser);
const mockedUseSession = vi.mocked(useSession);

function buildSessionUser(overrides: Partial<SessionUser> = {}): SessionUser {
	return {
		id: 'admin-1',
		email: 'admin@example.com',
		is_admin: true,
		auth_provider: 'local',
		email_verified: true,
		full_name: 'Admin User',
		company: 'PortfolioForge',
		profile_completed: true,
		assistant_eligible: true,
		can_use_project_assistant: true,
		created_at: '2026-04-15T00:00:00Z',
		...overrides,
	};
}

function buildAdminUserSummary(overrides: Partial<AdminUserSummary> = {}): AdminUserSummary {
	return {
		id: 'user-1',
		email: 'ada@example.com',
		is_admin: false,
		auth_provider: 'local',
		email_verified: true,
		full_name: 'Ada Lovelace',
		company: 'Analytical Engines',
		created_at: '2026-04-15T00:00:00Z',
		updated_at: '2026-04-15T01:00:00Z',
		last_login_at: '2026-04-15T02:00:00Z',
		...overrides,
	};
}

function buildAdminUserDetail(overrides: Partial<AdminUserDetail> = {}): AdminUserDetail {
	return {
		...buildAdminUserSummary(),
		...overrides,
	};
}

function mockSession(overrides: Partial<SessionUser> = {}) {
	mockedUseSession.mockReturnValue({
		user: buildSessionUser(overrides),
		token: 'token',
		loading: false,
		login: vi.fn(),
		refreshSession: vi.fn(async () => null),
		setUser: vi.fn(),
		logout: vi.fn(),
	});
}

function renderListPage() {
	return render(
		<MemoryRouter>
			<AdminUserListPage />
		</MemoryRouter>,
	);
}

function renderFormPage() {
	return render(
		<MemoryRouter initialEntries={['/admin/users/user-2']}>
			<Routes>
				<Route path="/admin/users" element={<p>users destination</p>} />
				<Route path="/admin/users/:id" element={<AdminUserFormPage />} />
			</Routes>
		</MemoryRouter>,
	);
}

describe('Admin user pages', () => {
	beforeEach(() => {
		mockSession();
		mockedDeleteAdminUser.mockReset();
		mockedFetchAdminUserById.mockReset();
		mockedFetchAdminUsers.mockReset();
		mockedUpdateAdminUser.mockReset();
		vi.spyOn(window, 'confirm').mockReturnValue(true);
	});

	afterEach(() => {
		cleanup();
		vi.restoreAllMocks();
	});

	it('lists active users and removes a soft-deleted standard user from the table', async () => {
		mockedFetchAdminUsers.mockResolvedValue({
			items: [
				buildAdminUserSummary({ id: 'user-1', email: 'ada@example.com' }),
				buildAdminUserSummary({ id: 'admin-1', email: 'admin@example.com', is_admin: true, full_name: 'Admin User' }),
			],
		});
		mockedDeleteAdminUser.mockResolvedValue(undefined);

		renderListPage();

		expect(await screen.findByRole('heading', { name: 'Users' })).toBeInTheDocument();
		expect(screen.getByText('Only active users are listed. Deleted identities remain reserved and lose access immediately.')).toBeInTheDocument();
		expect(screen.getByRole('link', { name: 'ada@example.com' })).toHaveAttribute('href', '/admin/users/user-1');

		const firstDeleteButton = screen.getAllByRole('button', { name: 'Delete' }).at(0);
		expect(firstDeleteButton).toBeDefined();
		if (!firstDeleteButton) {
			throw new Error('Expected at least one delete button');
		}

		fireEvent.click(firstDeleteButton);

		await waitFor(() => {
			expect(mockedDeleteAdminUser).toHaveBeenCalledWith('user-1');
		});
		await waitFor(() => {
			expect(screen.queryByText('ada@example.com')).not.toBeInTheDocument();
		});
		expect(screen.getByText('admin@example.com')).toBeInTheDocument();
	});

	it('loads immutable user detail and submits is_admin-only updates', async () => {
		const targetUser = buildAdminUserDetail({
			id: 'user-2',
			email: 'grace@example.com',
			is_admin: false,
			full_name: 'Grace Hopper',
			company: 'Compilers Inc.',
		});
		mockedFetchAdminUserById.mockResolvedValue(targetUser);
		mockedUpdateAdminUser.mockResolvedValue({ ...targetUser, is_admin: true });

		renderFormPage();

		expect(await screen.findByRole('heading', { name: 'User detail' })).toBeInTheDocument();
		expect(screen.getByDisplayValue('grace@example.com')).toHaveAttribute('readonly');
		expect(screen.getByDisplayValue('Grace Hopper')).toHaveAttribute('readonly');
		expect(screen.getByText(/This flow only mutates/i)).toBeInTheDocument();

		fireEvent.click(screen.getByRole('checkbox', { name: 'Grant admin access' }));
		fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));

		await waitFor(() => {
			expect(mockedUpdateAdminUser).toHaveBeenCalledWith('user-2', { is_admin: true });
		});
		expect(screen.getByRole('checkbox', { name: 'Grant admin access' })).toBeChecked();
	});
});
