import { useEffect } from 'react';
import Editor from '@monaco-editor/react';
import {
  AlertTriangle,
  ArrowLeft,
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  Clock,
  Code,
  LoaderCircle,
  Send,
  Target,
  Timer,
  Eye,
} from 'lucide-react';
import './ExamPreview.css';

const MONACO_LANG_MAP = {
  python: 'python',
  javascript: 'javascript',
  java: 'java',
  cpp: 'cpp',
  go: 'go',
};

const VALID_LANGUAGES = Object.keys(MONACO_LANG_MAP);

function formatPreviewTimer(totalMinutes) {
  if (!Number.isFinite(totalMinutes) || totalMinutes <= 0) {
    return null;
  }

  const totalSeconds = totalMinutes * 60;
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (hours > 0) {
    return `${hours}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
  }

  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
}

function normalizeCodeTemplates(rawTemplates) {
  let parsedTemplates = rawTemplates;

  if (parsedTemplates == null) {
    parsedTemplates = {};
  }

  if (typeof parsedTemplates === 'string') {
    for (let attempt = 0; attempt < 3 && typeof parsedTemplates === 'string'; attempt += 1) {
      try {
        parsedTemplates = JSON.parse(parsedTemplates);
      } catch {
        parsedTemplates = {};
        break;
      }
    }
  }

  const cleanTemplates = {};

  if (Array.isArray(parsedTemplates)) {
    parsedTemplates.forEach((templateEntry) => {
      if (
        templateEntry &&
        typeof templateEntry === 'object' &&
        VALID_LANGUAGES.includes(templateEntry.language) &&
        typeof templateEntry.template === 'string'
      ) {
        cleanTemplates[templateEntry.language] = templateEntry.template;
      }
    });
  } else if (parsedTemplates && typeof parsedTemplates === 'object') {
    Object.entries(parsedTemplates).forEach(([language, template]) => {
      if (VALID_LANGUAGES.includes(language) && typeof template === 'string') {
        cleanTemplates[language] = template;
      }
    });
  }

  return cleanTemplates;
}

export function buildPreviewChallenges(examItems) {
  return examItems
    .map((item, index) => {
      const challenge = item.challenge || item.Challenge || {};
      const templates = normalizeCodeTemplates(
        challenge.code_templates || challenge.CodeTemplates || challenge.codeTemplates,
      );

      return {
        ...challenge,
        id: challenge.id || challenge.ID || item.challenge_id || item.challengeID,
        title: challenge.title || challenge.Title || `Reto ${index + 1}`,
        description: challenge.description || challenge.Description || '',
        difficulty: (challenge.difficulty || challenge.Difficulty || 'medium').toLowerCase(),
        constraints: challenge.constraints || challenge.Constraints || '',
        points: item.points || item.Points || 0,
        order: item.order || item.Order || index + 1,
        code_templates: templates,
      };
    })
    .filter((challenge) => challenge.id)
    .sort((left, right) => left.order - right.order);
}

function getDifficultyLabel(difficulty) {
  if (difficulty === 'easy') return 'Facil';
  if (difficulty === 'hard') return 'Dificil';
  return 'Medio';
}

function formatAttemptLimit(tryLimit) {
  if (tryLimit === -1) {
    return '0/∞ intentos';
  }

  return `0/${tryLimit} intentos`;
}

function ExamPreview({
  examTitle,
  examDescription,
  timeLimitMinutes,
  tryLimit,
  challenges,
  publicTestCasesMap,
  previewCodeMap,
  setPreviewCodeMap,
  previewLanguage,
  setPreviewLanguage,
  previewCurrentIndex,
  setPreviewCurrentIndex,
  onExit,
  isLoading = false,
}) {
  const currentChallenge = challenges[previewCurrentIndex] || null;
  const currentCode = currentChallenge ? previewCodeMap[currentChallenge.id] || '' : '';
  const previewTimer = formatPreviewTimer(timeLimitMinutes);

  useEffect(() => {
    if (!currentChallenge) {
      return;
    }

    const templates = currentChallenge.code_templates || {};
    const languages = Object.keys(templates).filter((language) => VALID_LANGUAGES.includes(language));

    if (languages.length === 0) {
      if (previewLanguage !== 'python') {
        setPreviewLanguage('python');
      }
      return;
    }

    if (!languages.includes(previewLanguage)) {
      setPreviewLanguage(languages[0]);
    }

    if (!previewCodeMap[currentChallenge.id]) {
      setPreviewCodeMap((previous) => ({
        ...previous,
        [currentChallenge.id]: templates[languages[0]] || '',
      }));
    }
  }, [
    currentChallenge,
    previewCodeMap,
    previewLanguage,
    setPreviewCodeMap,
    setPreviewLanguage,
  ]);

  const handleSelectChallenge = (index) => {
    setPreviewCurrentIndex(index);
  };

  const handleCodeChange = (value) => {
    if (!currentChallenge) {
      return;
    }

    setPreviewCodeMap((previous) => ({
      ...previous,
      [currentChallenge.id]: value || '',
    }));
  };

  const handleLanguageChange = (nextLanguage) => {
    setPreviewLanguage(nextLanguage);

    if (!currentChallenge) {
      return;
    }

    const template = currentChallenge.code_templates?.[nextLanguage] || '';
    setPreviewCodeMap((previous) => ({
      ...previous,
      [currentChallenge.id]: previous[currentChallenge.id] || template,
    }));
  };

  return (
    <div className="exam-preview-page">
      <div className="exam-preview-notice">
        <div className="exam-preview-notice-copy">
          <span className="exam-preview-badge">
            <Eye size={14} /> Modo Preview
          </span>
          <p>
            Estás viendo la experiencia del estudiante. Puedes navegar y escribir codigo, pero nada se guarda ni se envía.
          </p>
        </div>
        <button type="button" className="exam-preview-exit" onClick={onExit}>
          <ArrowLeft size={16} /> Volver a configuracion
        </button>
      </div>

      <div className="exam-preview-frame">
        <div className="exam-preview-topbar">
          <div className="exam-preview-topbar-title">
            <Target size={20} />
            <div>
              <h3>{examTitle || 'Vista previa del examen'}</h3>
              <p>{examDescription || 'Sin descripcion disponible.'}</p>
            </div>
          </div>

          <div className="exam-preview-topbar-meta">
            <div className="exam-preview-meta-pill muted">
              <CheckCircle2 size={15} /> 0/{challenges.length} resueltos
            </div>
            {previewTimer ? (
              <div className="exam-preview-meta-pill accent">
                <Timer size={15} /> {previewTimer}
              </div>
            ) : (
              <div className="exam-preview-meta-pill unlimited">
                <Clock size={15} /> Sin límite
              </div>
            )}
            <button type="button" className="exam-preview-finish" disabled>
              <ArrowLeft size={15} /> Terminar Actividad
            </button>
          </div>
        </div>

        <div className="exam-preview-body">
          <aside className="exam-preview-sidebar">
            <div className="exam-preview-sidebar-title">Retos del examen</div>
            {challenges.map((challenge, index) => {
              const isActive = index === previewCurrentIndex;
              return (
                <button
                  key={challenge.id}
                  type="button"
                  onClick={() => handleSelectChallenge(index)}
                  className={`exam-preview-sidebar-item ${isActive ? 'active' : ''}`}
                >
                  <div className={`exam-preview-sidebar-index ${isActive ? 'active' : ''}`}>
                    {index + 1}
                  </div>
                  <div className="exam-preview-sidebar-copy">
                    <strong>{challenge.title}</strong>
                    <span>
                      {challenge.points} pts • {getDifficultyLabel(challenge.difficulty)}
                    </span>
                  </div>
                </button>
              );
            })}
          </aside>

          <div className="exam-preview-content">
            {currentChallenge ? (
              <>
                <section className="exam-preview-description">
                  <div className="exam-preview-challenge-header">
                    <div className="exam-preview-challenge-title-block">
                      <div className="exam-preview-challenge-title">
                        <Code size={18} />
                        <h4>{currentChallenge.title}</h4>
                      </div>
                      <div className="exam-preview-challenge-tags">
                        <span className={`difficulty ${currentChallenge.difficulty}`}>
                          {getDifficultyLabel(currentChallenge.difficulty)}
                        </span>
                        <span className="points">{currentChallenge.points} pts</span>
                      </div>
                    </div>
                  </div>

                  <p className="exam-preview-text">{currentChallenge.description || 'Sin descripcion.'}</p>

                  {currentChallenge.constraints ? (
                    <div className="exam-preview-constraints">
                      <h5>Restricciones / explicacion</h5>
                      <p>{currentChallenge.constraints}</p>
                    </div>
                  ) : null}

                  {publicTestCasesMap[currentChallenge.id]?.length > 0 ? (
                    <div className="exam-preview-cases">
                      <h5>Casos de prueba visibles</h5>
                      {publicTestCasesMap[currentChallenge.id].map((testCase, index) => {
                        const inputValues = Array.isArray(testCase.input) ? testCase.input : [];
                        const firstInput = inputValues[0] || {};
                        const expectedOutput =
                          testCase.expected_output ||
                          testCase.ExpectedOutput ||
                          testCase.expectedOutput ||
                          {};

                        return (
                          <div key={`${currentChallenge.id}-case-${index}`} className="exam-preview-case-card">
                            <div>
                              <strong>Entrada ({firstInput.name || firstInput.Name || 'entrada'})</strong>
                              <pre>{firstInput.value || ''}</pre>
                            </div>
                            <div>
                              <strong>Salida esperada</strong>
                              <textarea readOnly value={expectedOutput.value || ''} rows="5" />
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  ) : (
                    <div className="exam-preview-empty-note">
                      <AlertTriangle size={16} /> Este reto no tiene casos de prueba publicos para mostrar en preview.
                    </div>
                  )}

                  <div className="exam-preview-nav">
                    <button
                      type="button"
                      onClick={() => handleSelectChallenge(previewCurrentIndex - 1)}
                      disabled={previewCurrentIndex === 0}
                    >
                      <ChevronLeft size={16} /> Anterior
                    </button>
                    <button
                      type="button"
                      onClick={() => handleSelectChallenge(previewCurrentIndex + 1)}
                      disabled={previewCurrentIndex === challenges.length - 1}
                    >
                      Siguiente <ChevronRight size={16} />
                    </button>
                  </div>
                </section>

                <section className="exam-preview-editor-panel">
                  <div className="exam-preview-editor-toolbar">
                    <div className="exam-preview-editor-toolbar-start">
                      <select
                        value={previewLanguage}
                        onChange={(event) => handleLanguageChange(event.target.value)}
                      >
                        {(() => {
                          const templates = currentChallenge.code_templates || {};
                          const languages = Object.keys(templates).filter((language) =>
                            VALID_LANGUAGES.includes(language),
                          );

                          if (languages.length === 0) {
                            return <option value="python">Python</option>;
                          }

                          return languages.map((language) => (
                            <option key={language} value={language}>
                              {language === 'cpp'
                                ? 'C++'
                                : language.charAt(0).toUpperCase() + language.slice(1)}
                            </option>
                          ));
                        })()}
                      </select>
                    </div>

                    <div className="exam-preview-editor-toolbar-end">
                      <span className="exam-preview-attempts">{formatAttemptLimit(tryLimit)}</span>
                      <div className="exam-preview-editor-actions">
                        <button type="button" disabled className="guide-button secondary" title="Vista guía del botón que verá el estudiante">
                          <Timer size={14} /> Probar Código
                        </button>
                        <button type="button" disabled className="guide-button primary" title="Vista guía del botón que verá el estudiante">
                          <Send size={14} /> Enviar Solución
                        </button>
                      </div>
                    </div>
                  </div>

                  <div className="exam-preview-toolbar-guide">
                    Así se muestran el temporizador y las acciones principales durante la actividad real del estudiante.
                  </div>

                  <div className="exam-preview-editor-shell">
                    {isLoading ? (
                      <div className="exam-preview-loader-overlay" role="status" aria-live="polite">
                        <div className="exam-preview-loader-card">
                          <LoaderCircle size={24} className="exam-preview-loader-icon" />
                          <strong>Preparando preview</strong>
                          <p>Cargando retos, plantillas y casos visibles para mostrar la vista del estudiante.</p>
                        </div>
                      </div>
                    ) : null}
                    <Editor
                      key={`${currentChallenge.id}-${previewLanguage}-preview`}
                      height="100%"
                      theme="vs-dark"
                      language={MONACO_LANG_MAP[previewLanguage] || 'python'}
                      value={currentCode}
                      onChange={handleCodeChange}
                      options={{
                        fontSize: 14,
                        minimap: { enabled: false },
                        scrollBeyondLastLine: false,
                        automaticLayout: true,
                      }}
                    />
                  </div>

                  <div className="exam-preview-console">
                    <div className="exam-preview-console-header">Salida / Consola</div>
                    <div className="exam-preview-console-body">
                      <p>
                        El preview replica la interfaz del estudiante, pero la ejecucion real esta deshabilitada en esta vista.
                      </p>
                      <p>
                        Puedes usar este espacio para revisar distribucion, lectura del reto, casos visibles y sensacion general del examen.
                      </p>
                    </div>
                  </div>
                </section>
              </>
            ) : (
              <div className="exam-preview-no-challenges">
                <AlertTriangle size={18} /> Este examen aun no tiene retos asignados para previsualizar.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default ExamPreview;