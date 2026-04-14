import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { AppError } from '../../shared/api/errors';
import { ProjectAssistantChat } from './ProjectAssistantChat';
import { sendProjectAssistantMessage } from './api';

vi.mock('./api', () => ({
  sendProjectAssistantMessage: vi.fn(),
}));

const mockedSendProjectAssistantMessage = vi.mocked(sendProjectAssistantMessage);

describe('ProjectAssistantChat', () => {
  beforeEach(() => {
    mockedSendProjectAssistantMessage.mockReset();
  });

  afterEach(() => {
    cleanup();
  });

  it('does not render when disabled', () => {
    const { container } = render(<ProjectAssistantChat slug="portfolioforge" enabled={false} lang="es" />);

    expect(container).toBeEmptyDOMElement();
  });

  it('sends the approved contract and renders the answer', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Grounded answer.' });

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'How does the architecture work?' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledWith('portfolioforge', {
        question: 'How does the architecture work?',
        history: [],
        lang: 'en',
      });
    });

    expect(await screen.findByText('Grounded answer.')).toBeInTheDocument();
    expect(screen.getByText('How does the architecture work?')).toBeInTheDocument();
  });

  it('restores history and shows backend errors', async () => {
    mockedSendProjectAssistantMessage.mockRejectedValue(new AppError(409, {
      code: 'assistant_unavailable',
      message: 'Assistant unavailable.',
    }));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Hi there' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    expect(await screen.findByText('Assistant unavailable.')).toBeInTheDocument();
    expect(screen.queryByText('Hi there')).not.toBeInTheDocument();
  });
});
