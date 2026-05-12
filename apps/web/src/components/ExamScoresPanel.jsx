import { useState, useCallback, useEffect } from 'react';
import { BarChart3, ChevronDown, ChevronRight, Trophy, Layers } from 'lucide-react';
import PageLoader from './PageLoader';
import {
  getUsersScoresByExam,
  getExamScoringForUser,
  getExamScoreItemDetails,
} from '../api/scoring';
import './ExamScoresPanel.css';

const pick = (obj, ...keys) => {
  for (const k of keys) {
    if (obj && obj[k] !== undefined && obj[k] !== null) return obj[k];
  }
  return undefined;
};

const formatShortDate = (value) => {
  if (!value) return '—';
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? '—' : d.toLocaleString([], { dateStyle: 'short', timeStyle: 'short' });
};

/**
 * Panel de calificaciones oficiales (API scoring).
 * - Modo clase: profesor dueño del examen (showClassScores).
 * - Modo propio: estudiante o cualquier usuario viendo solo su userId.
 */
const ExamScoresPanel = ({
  examId,
  courseId,
  userId,
  isProfessor,
  showClassScores,
  displayNames = {},
  triggerLabelClass = 'Calificaciones',
  triggerLabelMine = 'Mis calificaciones',
  compact = false,
}) => {
  const [expanded, setExpanded] = useState(false);
  const [view, setView] = useState(showClassScores ? 'class' : 'mine'); // 'class' | 'mine'

  const [classState, setClassState] = useState({ loading: false, loaded: false, error: '', data: null });
  const [mineState, setMineState] = useState({
    loading: false,
    loaded: false,
    error: '',
    data: null,
    empty: false,
  });
  const [itemDetailByScoreId, setItemDetailByScoreId] = useState({});
  const [openScoreIds, setOpenScoreIds] = useState({});

  const uid = String(userId || '');

  useEffect(() => {
    setExpanded(false);
    setView(showClassScores ? 'class' : 'mine');
    setClassState({ loading: false, loaded: false, error: '', data: null });
    setMineState({ loading: false, loaded: false, error: '', data: null, empty: false });
    setItemDetailByScoreId({});
    setOpenScoreIds({});
  }, [examId, showClassScores]);

  const loadClass = useCallback(async () => {
    if (classState.loaded || classState.loading) return;
    setClassState((s) => ({ ...s, loading: true, error: '' }));
    try {
      const data = await getUsersScoresByExam(examId, courseId);
      setClassState({ loading: false, loaded: true, error: '', data });
    } catch (err) {
      const msg =
        err?.response?.data?.error ||
        err?.response?.data?.message ||
        err?.message ||
        'No se pudieron cargar las calificaciones.';
      setClassState({ loading: false, loaded: true, error: msg, data: null });
    }
  }, [examId, courseId, classState.loaded, classState.loading]);

  const loadMine = useCallback(async () => {
    if (!uid) return;
    if (mineState.loaded || mineState.loading) return;
    setMineState((s) => ({ ...s, loading: true, error: '' }));
    try {
      const data = await getExamScoringForUser(examId, uid);
      setMineState({ loading: false, loaded: true, error: '', data, empty: false });
    } catch (err) {
      const raw =
        err?.response?.data?.error ||
        err?.response?.data?.message ||
        err?.message ||
        '';
      const noScores = /no exam scores/i.test(String(raw));
      if (noScores) {
        setMineState({
          loading: false,
          loaded: true,
          error: '',
          data: null,
          empty: true,
        });
      } else {
        setMineState({
          loading: false,
          loaded: true,
          error: raw || 'No se pudieron cargar tus calificaciones.',
          data: null,
          empty: false,
        });
      }
    }
  }, [examId, uid, mineState.loaded, mineState.loading]);

  const toggle = () => {
    const next = !expanded;
    setExpanded(next);
    if (next) {
      if (view === 'class' && showClassScores) loadClass();
      else loadMine();
    }
  };

  const switchView = (v) => {
    setView(v);
    if (v === 'class' && showClassScores) loadClass();
    if (v === 'mine') loadMine();
  };

  const toggleItemDetail = async (examScoreId) => {
    setOpenScoreIds((prev) => ({ ...prev, [examScoreId]: !prev[examScoreId] }));
    if (itemDetailByScoreId[examScoreId]?.loaded || itemDetailByScoreId[examScoreId]?.loading) return;
    setItemDetailByScoreId((prev) => ({
      ...prev,
      [examScoreId]: { ...prev[examScoreId], loading: true, error: '' },
    }));
    try {
      const detail = await getExamScoreItemDetails(examScoreId);
      setItemDetailByScoreId((prev) => ({
        ...prev,
        [examScoreId]: { loading: false, loaded: true, error: '', detail },
      }));
    } catch (err) {
      const msg =
        err?.response?.data?.error ||
        err?.response?.data?.message ||
        err?.message ||
        'Error al cargar el detalle por ítem.';
      setItemDetailByScoreId((prev) => ({
        ...prev,
        [examScoreId]: { loading: false, loaded: false, error: msg },
      }));
    }
  };

  const usersList = Array.isArray(classState.data?.users) ? classState.data.users : [];

  const renderExamScoresRows = (examScoresBlock) => {
    const best = pick(examScoresBlock, 'best_score', 'BestScore') ?? '—';
    const scores = Array.isArray(examScoresBlock?.exam_scores)
      ? examScoresBlock.exam_scores
      : Array.isArray(examScoresBlock?.ExamScores)
        ? examScoresBlock.ExamScores
        : [];

    return (
      <div className="exam-scores-panel-block">
        <div className="exam-scores-best">
          <Trophy size={14} />
          <span>Mejor puntuación: <strong>{best}</strong></span>
        </div>
        {scores.length === 0 ? (
          <p className="exam-scores-muted">Sin intentos registrados.</p>
        ) : (
          <ul className="exam-scores-attempts">
            {scores.map((row) => {
              const sid = pick(row, 'id', 'ID');
              const sc = pick(row, 'score', 'Score');
              const created = pick(row, 'created_at', 'CreatedAt', 'createdAt');
              const open = openScoreIds[sid];
              const detail = itemDetailByScoreId[sid];
              return (
                <li key={sid} className="exam-scores-attempt">
                  <button
                    type="button"
                    className="exam-scores-attempt-head"
                    onClick={() => toggleItemDetail(sid)}
                  >
                    {open ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                    <span>Intento: <strong>{sc}</strong> pts</span>
                    <span className="exam-scores-muted">{formatShortDate(created)}</span>
                  </button>
                  {open && (
                    <div className="exam-scores-attempt-body">
                      {detail?.loading && (
                        <PageLoader message="Cargando ítems…" compact minHeight="0" size={14} />
                      )}
                      {detail?.error && <p className="exam-scores-err">{detail.error}</p>}
                      {detail?.detail && (
                        <ul className="exam-scores-items">
                          {(detail.detail.exam_items || detail.detail.ExamItems || []).map((it, idx) => {
                            const ei = it.exam_item || it.ExamItem || {};
                            const eis = it.exam_item_score || it.ExamItemScore || {};
                            const ord = pick(ei, 'order', 'Order');
                            const chall = pick(ei, 'challenge_id', 'challengeId', 'ChallengeID');
                            const maxPts = pick(ei, 'points', 'Points');
                            const title =
                              pick(ei, 'title', 'Title') ||
                              (ord !== undefined && ord !== null
                                ? `Ítem #${ord}${typeof maxPts === 'number' ? ` · máx ${maxPts} pts` : ''}`
                                : chall
                                  ? `Reto ${String(chall).slice(0, 12)}`
                                  : `Ítem ${idx + 1}`);
                            const pts = pick(eis, 'score', 'Score');
                            return (
                              <li key={pick(eis, 'id', 'ID') || idx}>
                                <span className="exam-scores-item-title">{title}</span>
                                <span className="exam-scores-item-pts">{pts ?? '—'} pts</span>
                              </li>
                            );
                          })}
                        </ul>
                      )}
                    </div>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </div>
    );
  };

  if (!uid && !showClassScores) return null;

  return (
    <div className={`exam-scores-panel ${compact ? 'compact' : ''}`}>
      <button type="button" className="exam-scores-trigger" onClick={toggle}>
        <BarChart3 size={16} />
        <span>{showClassScores && isProfessor ? triggerLabelClass : triggerLabelMine}</span>
        {expanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
      </button>

      {expanded && (
        <div className="exam-scores-body">
          {showClassScores && isProfessor && (
            <div className="exam-scores-tabs">
              <button
                type="button"
                className={view === 'class' ? 'active' : ''}
                onClick={() => switchView('class')}
              >
                Clase
              </button>
              <button
                type="button"
                className={view === 'mine' ? 'active' : ''}
                onClick={() => switchView('mine')}
              >
                Mis intentos
              </button>
            </div>
          )}

          {view === 'class' && showClassScores && (
            <>
              {classState.loading && (
                <PageLoader message="Cargando calificaciones…" compact minHeight="60px" size={16} />
              )}
              {classState.error && <p className="exam-scores-err">{classState.error}</p>}
              {!classState.loading && classState.loaded && !classState.error && usersList.length === 0 && (
                <p className="exam-scores-muted">Ningún estudiante tiene calificación registrada en este examen.</p>
              )}
              {!classState.loading && usersList.length > 0 && (
                <ul className="exam-scores-user-list">
                  {usersList.map((urow) => {
                    const suid = pick(urow, 'user_id', 'UserID', 'userId');
                    const block = urow.exam_scores || urow.ExamScores;
                    const label = displayNames[String(suid)] || String(suid);
                    return (
                      <li key={String(suid)} className="exam-scores-user-card">
                        <div className="exam-scores-user-head">
                          <Layers size={14} />
                          <span title={String(suid)}>{label}</span>
                        </div>
                        {renderExamScoresRows(block)}
                      </li>
                    );
                  })}
                </ul>
              )}
            </>
          )}

          {(view === 'mine' || !showClassScores || !isProfessor) && (
            <>
              {!uid ? (
                <p className="exam-scores-muted">Inicia sesión para ver tus calificaciones.</p>
              ) : (
                <>
                  {mineState.loading && (
                    <PageLoader message="Cargando tus calificaciones…" compact minHeight="60px" size={16} />
                  )}
                  {mineState.error && <p className="exam-scores-err">{mineState.error}</p>}
                  {mineState.empty && !mineState.error && (
                    <p className="exam-scores-muted">Aún no hay calificación registrada para ti en este examen.</p>
                  )}
                  {mineState.data && renderExamScoresRows(mineState.data)}
                </>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default ExamScoresPanel;
