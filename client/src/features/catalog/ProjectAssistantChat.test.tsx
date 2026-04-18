import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { AppError } from '../../shared/api/errors';
import { ProjectAssistantChat } from './ProjectAssistantChat';
import { sendProjectAssistantMessage } from './api';

vi.mock('./api', () => ({
  sendProjectAssistantMessage: vi.fn(),
}));

const mockedSendProjectAssistantMessage = vi.mocked(sendProjectAssistantMessage);
const scrollIntoViewMock = vi.fn();

interface DeferredResponse {
  promise: Promise<{ answer: string }>;
  resolve: (value: { answer: string }) => void;
}

function createDeferredResponse(): DeferredResponse {
  let resolve!: DeferredResponse['resolve'];
  const promise = new Promise<{ answer: string }>((innerResolve) => {
    resolve = innerResolve;
  });

  return { promise, resolve };
}

describe('ProjectAssistantChat', () => {
  beforeEach(() => {
    mockedSendProjectAssistantMessage.mockReset();
    scrollIntoViewMock.mockReset();
    Object.defineProperty(Element.prototype, 'scrollIntoView', {
      configurable: true,
      value: scrollIntoViewMock,
    });
    window.sessionStorage.clear();
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
    expect(window.sessionStorage.getItem('assistant_history:portfolioforge')).toContain('Grounded answer.');
  });

  it('restores only the newest valid history entries and sends the same bounded subset', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Fresh answer.' });

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'user', content: 'Message 1' },
      { role: 'assistant', content: 'Message 2' },
      { role: 'user', content: 'Message 3' },
      { role: 'assistant', content: 'Message 4' },
      { role: 'user', content: 'Message 5' },
      { role: 'assistant', content: 'Message 6' },
      { role: 'user', content: 'Message 7' },
      { role: 'assistant', content: 'Message 8' },
      { role: 'user', content: 'Message 9' },
      { role: 'assistant', content: 'Message 10' },
      { role: 'system', content: 'Ignore me' },
      { role: 'assistant', content: '   ' },
      { nope: true },
    ]));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));

    expect(screen.queryByText('Message 1')).not.toBeInTheDocument();
    expect(screen.queryByText('Message 2')).not.toBeInTheDocument();
    expect(screen.getByText('Message 3')).toBeInTheDocument();
    expect(screen.getByText('Message 10')).toBeInTheDocument();
    expect(screen.queryByText('Ignore me')).not.toBeInTheDocument();

    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'What changed most recently?' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledWith('portfolioforge', {
        question: 'What changed most recently?',
        history: [
          { role: 'user', content: 'Message 3' },
          { role: 'assistant', content: 'Message 4' },
          { role: 'user', content: 'Message 5' },
          { role: 'assistant', content: 'Message 6' },
          { role: 'user', content: 'Message 7' },
          { role: 'assistant', content: 'Message 8' },
          { role: 'user', content: 'Message 9' },
          { role: 'assistant', content: 'Message 10' },
        ],
        lang: 'en',
      });
    });
  });

  it('submits on Enter and keeps Shift+Enter for editing', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Answer.' });

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    const textbox = screen.getByRole('textbox');

    fireEvent.change(textbox, { target: { value: 'First line' } });
    fireEvent.keyDown(textbox, { key: 'Enter', code: 'Enter' });

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledTimes(1);
    });

    fireEvent.change(textbox, { target: { value: 'Line 1' } });
    fireEvent.keyDown(textbox, { key: 'Enter', code: 'Enter', shiftKey: true });

    expect(mockedSendProjectAssistantMessage).toHaveBeenCalledTimes(1);
    expect(textbox).toHaveValue('Line 1');
  });

  it('keeps prior transcript visible while a reply is pending', async () => {
    const deferred = createDeferredResponse();
    mockedSendProjectAssistantMessage.mockReturnValue(deferred.promise);

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Restored answer.' },
    ]));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Tell me more about the rollout.' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    expect(screen.getByText('Restored answer.')).toBeInTheDocument();
    expect(screen.getByText('Tell me more about the rollout.')).toBeInTheDocument();
    expect(screen.getAllByText('Thinking…')).toHaveLength(2);

    deferred.resolve({ answer: 'Here is more detail.' });

    expect(await screen.findByText('Here is more detail.')).toBeInTheDocument();
  });

  it('auto-scrolls the transcript when the conversation state changes', async () => {
    const deferred = createDeferredResponse();
    mockedSendProjectAssistantMessage.mockReturnValue(deferred.promise);

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    const callsAfterOpen = scrollIntoViewMock.mock.calls.length;

    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Explain the deployment flow' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(scrollIntoViewMock.mock.calls.length).toBeGreaterThan(callsAfterOpen);
    });

    const callsWhilePending = scrollIntoViewMock.mock.calls.length;
    deferred.resolve({ answer: 'Deployment runs from CI after review.' });

    expect(await screen.findByText('Deployment runs from CI after review.')).toBeInTheDocument();

    await waitFor(() => {
      expect(scrollIntoViewMock.mock.calls.length).toBeGreaterThan(callsWhilePending);
    });
  });

  it('restores the draft and preserves previous transcript when the request fails', async () => {
    mockedSendProjectAssistantMessage.mockRejectedValue(new AppError(409, {
      code: 'assistant_unavailable',
      message: 'Assistant unavailable.',
    }));

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Restored answer.' },
    ]));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Hi there' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    expect(await screen.findByText('Assistant unavailable.')).toBeInTheDocument();
    expect(screen.getByText('Restored answer.')).toBeInTheDocument();
    expect(screen.getAllByText('Hi there')).toHaveLength(1);
    expect(screen.getByRole('textbox')).toHaveValue('Hi there');
  });

  it('does not reuse history when the assistant slug changes', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Fresh answer.' });

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Portfolio history.' },
    ]));

    const { rerender } = render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    expect(screen.getByText('Portfolio history.')).toBeInTheDocument();

    rerender(<ProjectAssistantChat slug="other-project" enabled lang="en" />);

    expect(screen.queryByText('Portfolio history.')).not.toBeInTheDocument();
    expect(screen.getByText('Try asking about architecture, results, integrations, or tradeoffs.')).toBeInTheDocument();

    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'What is different here?' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledWith('other-project', {
        question: 'What is different here?',
        history: [],
        lang: 'en',
      });
    });
  });

  it('keeps the current session context after closing and reopening the chat', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Next answer.' });

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Restored answer.' },
    ]));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    expect(screen.getByText('Restored answer.')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Close assistant' }));
    expect(screen.queryByLabelText('Project assistant')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    expect(screen.getByText('Restored answer.')).toBeInTheDocument();

    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Continue the conversation' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledWith('portfolioforge', {
        question: 'Continue the conversation',
        history: [{ role: 'assistant', content: 'Restored answer.' }],
        lang: 'en',
      });
    });
  });

  it('clears the local draft, removes current slug history, and restarts with empty history', async () => {
    mockedSendProjectAssistantMessage.mockResolvedValue({ answer: 'Fresh answer.' });

    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Current project answer.' },
    ]));
    window.sessionStorage.setItem('assistant_history:other-project', JSON.stringify([
      { role: 'assistant', content: 'Other project answer.' },
    ]));

    render(<ProjectAssistantChat slug="portfolioforge" enabled lang="en" />);

    fireEvent.click(screen.getByRole('button', { name: 'Ask project assistant' }));
    expect(screen.getByText('Current project answer.')).toBeInTheDocument();
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Draft to discard' } });

    fireEvent.click(screen.getByRole('button', { name: 'Clear chat' }));

    expect(screen.queryByText('Current project answer.')).not.toBeInTheDocument();
    expect(screen.getByRole('textbox')).toHaveValue('');
    expect(window.sessionStorage.getItem('assistant_history:portfolioforge')).toBeNull();
    expect(window.sessionStorage.getItem('assistant_history:other-project')).toContain('Other project answer.');

    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Start over' } });
    fireEvent.click(screen.getByRole('button', { name: 'Send' }));

    await waitFor(() => {
      expect(mockedSendProjectAssistantMessage).toHaveBeenCalledWith('portfolioforge', {
        question: 'Start over',
        history: [],
        lang: 'en',
      });
    });
  });
});
