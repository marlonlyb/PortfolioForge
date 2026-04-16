import { useEffect, useState } from 'react';

import { sendProjectAssistantMessage } from './api';
import { AppError } from '../../shared/api/errors';
import type { ProjectAssistantMessage } from '../../shared/types/project';
import { getAssistantHistory, setAssistantHistory } from '../../shared/storage';

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
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setHistory(getAssistantHistory(slug));
  }, [slug]);

  useEffect(() => {
    setAssistantHistory(slug, history);
  }, [history, slug]);

  if (!enabled) {
    return null;
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextQuestion = question.trim();
    if (nextQuestion.length < 2 || loading) {
      return;
    }

    const nextHistory = [...history, { role: 'user' as const, content: nextQuestion }];
    setLoading(true);
    setError(null);
    setHistory(nextHistory);
    setQuestion('');

    try {
      const response = await sendProjectAssistantMessage(slug, {
        question: nextQuestion,
        history,
        lang,
      });

      setHistory([...nextHistory, { role: 'assistant', content: response.answer }]);
    } catch (err: unknown) {
      setHistory(history);
      setError(err instanceof AppError ? err.message : 'Assistant unavailable.');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className={`assistant-chat ${open ? 'assistant-chat--open' : ''}`}>
      <button type="button" className="assistant-chat__toggle" onClick={() => setOpen((current) => !current)}>
        {open ? 'Close assistant' : 'Ask project assistant'}
      </button>

      {open ? (
        <section className="assistant-chat__panel card" aria-label="Project assistant">
          <div className="assistant-chat__header">
            <p className="eyebrow">Project assistant</p>
            <p className="assistant-chat__copy">Ask detailed questions grounded in the project documentation.</p>
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
          </div>

          <form className="assistant-chat__form" onSubmit={handleSubmit}>
            <textarea
              className="admin__textarea assistant-chat__input"
              rows={4}
              value={question}
              onChange={(event) => setQuestion(event.target.value)}
              placeholder="Ask a detailed question about this project"
            />

            {error ? <p className="admin__error">{error}</p> : null}

            <div className="assistant-chat__actions">
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
