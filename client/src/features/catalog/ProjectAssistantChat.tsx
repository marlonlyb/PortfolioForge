import { useEffect, useState, type FormEvent, type KeyboardEvent } from 'react';

import { sendProjectAssistantMessage } from './api';
import { AppError } from '../../shared/api/errors';
import type { ProjectAssistantMessage } from '../../shared/types/project';
import {
  clearAssistantHistory,
  getAssistantHistory,
  normalizeAssistantHistory,
  setAssistantHistory,
} from '../../shared/storage';

interface ProjectAssistantChatProps {
  slug: string;
  enabled: boolean;
  lang: string;
}

export function ProjectAssistantChat({ slug, enabled, lang }: ProjectAssistantChatProps) {
  const [open, setOpen] = useState(false);
  const [question, setQuestion] = useState('');
  const [history, setHistory] = useState<ProjectAssistantMessage[]>(() => getAssistantHistory(slug));
  const [loading, setLoading] = useState(false);
  const [pendingQuestion, setPendingQuestion] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setHistory(getAssistantHistory(slug));
    setQuestion('');
    setPendingQuestion(null);
    setError(null);
  }, [slug]);

  useEffect(() => {
    setAssistantHistory(slug, history);
  }, [history, slug]);

  if (!enabled) {
    return null;
  }

  async function submitQuestion() {
    const nextQuestion = question.trim();
    if (nextQuestion.length < 2 || loading) {
      return;
    }

    const requestHistory = normalizeAssistantHistory(history);
    setLoading(true);
    setError(null);
    setPendingQuestion(nextQuestion);
    setQuestion('');

    try {
      const response = await sendProjectAssistantMessage(slug, {
        question: nextQuestion,
        history: requestHistory,
        lang,
      });

      setHistory(normalizeAssistantHistory([
        ...requestHistory,
        { role: 'user', content: nextQuestion },
        { role: 'assistant', content: response.answer },
      ]));
      setPendingQuestion(null);
    } catch (err: unknown) {
      setPendingQuestion(null);
      setQuestion(nextQuestion);
      setError(err instanceof AppError ? err.message : 'Assistant unavailable.');
    } finally {
      setLoading(false);
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await submitQuestion();
  }

  function handleInputKeyDown(event: KeyboardEvent<HTMLTextAreaElement>) {
    if (event.key !== 'Enter' || event.shiftKey) {
      return;
    }

    if (loading || question.trim().length < 2) {
      return;
    }

    event.preventDefault();
    void submitQuestion();
  }

  function handleClearChat() {
    clearAssistantHistory(slug);
    setHistory([]);
    setQuestion('');
    setPendingQuestion(null);
    setError(null);
  }

  const hasConversation = history.length > 0 || pendingQuestion !== null;

  return (
    <div className={`assistant-chat ${open ? 'assistant-chat--open' : ''}`}>
      <button type="button" className="assistant-chat__toggle" onClick={() => setOpen((current) => !current)}>
        {open ? 'Close assistant' : 'Ask project assistant'}
      </button>

      {open ? (
        <section className="assistant-chat__panel card" aria-label="Project assistant">
          <div className="assistant-chat__header">
            <p className="eyebrow">Project assistant</p>
            <p className="assistant-chat__copy">
              {hasConversation
                ? 'Continue the conversation with context from this browser session only.'
                : 'Ask detailed questions grounded in the project documentation.'}
            </p>
          </div>

          <div className="assistant-chat__history">
            {history.length === 0 ? (
              <p className="assistant-chat__empty">Try asking about architecture, results, integrations, or tradeoffs.</p>
            ) : null}

            {history.map((message, index) => (
              <article key={`${message.role}-${index}`} className={`assistant-chat__message assistant-chat__message--${message.role}`}>
                <strong>{message.role === 'assistant' ? 'Assistant' : 'You'}</strong>
                <p>{message.content}</p>
              </article>
            ))}

            {pendingQuestion ? (
              <article className="assistant-chat__message assistant-chat__message--user assistant-chat__message--pending">
                <strong>You</strong>
                <p>{pendingQuestion}</p>
              </article>
            ) : null}

            {loading ? (
              <article className="assistant-chat__message assistant-chat__message--assistant assistant-chat__message--typing" aria-live="polite">
                <strong>Assistant</strong>
                <p>Thinking…</p>
              </article>
            ) : null}
          </div>

          <form className="assistant-chat__form" onSubmit={handleSubmit}>
            <textarea
              className="admin__textarea assistant-chat__input"
              rows={4}
              value={question}
              onChange={(event) => setQuestion(event.target.value)}
              onKeyDown={handleInputKeyDown}
              placeholder="Ask a detailed question about this project"
            />

            {error ? <p className="admin__error assistant-chat__error">{error}</p> : null}

            <div className="assistant-chat__actions">
              <button className="assistant-chat__clear" type="button" onClick={handleClearChat} disabled={loading || (!hasConversation && question.length === 0 && !error)}>
                Clear chat
              </button>
              <button className="btn" type="submit" disabled={loading || question.trim().length < 2}>
                {loading ? 'Thinking…' : 'Send'}
              </button>
            </div>
          </form>
        </section>
      ) : null}
    </div>
  );
}
