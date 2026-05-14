import { useState, useEffect, useContext, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft, BookOpen, ChevronRight, Trophy, Code, Users,
  CheckCircle, XCircle, Clock, BarChart3, Hash, User, Layers, Calendar
} from 'lucide-react';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import { getExamDetails } from '../api/exams';
import {
  getUsersScoresByExam,
  getExamScoringForUser,
  getExamScoreItemDetails,
} from '../api/scoring';
import PageLoader from '../components/PageLoader';
import './ExamResults.css';

/* ========================================================
   ExamResults — Dedicated per-exam results page
   ======================================================== */

const pick = (obj, ...keys) => {
  for (const k of keys) {
    if (obj && obj[k] !== undefined && obj[k] !== null) return obj[k];
  }
  return undefined;
};

const fmtDate = (v) => {
  if (!v) return '—';
  const d = new Date(v);
  return Number.isNaN(d.getTime())
    ? '—'
    : d.toLocaleString([], { dateStyle: 'short', timeStyle: 'short' });
};

const scoreClass = (pct) =>
  pct === 100 ? 'perfect' : pct >= 50 ? 'partial' : pct > 0 ? 'low' : 'none';

const ExamResults = () => {
  const { id: examId } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const isProfessor =
    user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';
  const userId = user?.id || user?.ID || '';

  // Exam metadata
  const [exam, setExam] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Exam items (challenges)
  const [examItems, setExamItems] = useState([]);

  // Student: own scoring
  const [myScoring, setMyScoring] = useState(null);
  const [myLoading, setMyLoading] = useState(false);
  const [myEmpty, setMyEmpty] = useState(false);

  // Professor: class scores
  const [classScores, setClassScores] = useState([]);
  const [classLoading, setClassLoading] = useState(false);

  // Item details cache { examScoreId: { loaded, detail } }
  const [itemDetails, setItemDetails] = useState({});

  // Expanded states
  const [expandedStudentId, setExpandedStudentId] = useState(null);
  const [expandedScoreId, setExpandedScoreId] = useState(null);

  // Code modal
  const [selectedSubmission, setSelectedSubmission] = useState(null);

  // Active tab for professor
  const [tab, setTab] = useState(isProfessor ? 'class' : 'mine');

  // ──────────────────────────────────────────
  // 1. Load exam + items
  // ──────────────────────────────────────────
  useEffect(() => {
    const load = async () => {
      try {
        const data = await getExamDetails(examId);
        setExam(data);

        const itemsRes = await client.get(`/exams/${examId}/items`);
        const items = Array.isArray(itemsRes.data)
          ? itemsRes.data
          : itemsRes.data?.items || [];
        setExamItems(items);
      } catch (err) {
        setError(err?.response?.data?.error || err?.message || 'Error cargando examen');
      } finally {
        setLoading(false);
      }
    };
    if (user) load();
  }, [examId, user]);

  // ──────────────────────────────────────────
  // 2. Load scores once exam is ready
  // ──────────────────────────────────────────
  useEffect(() => {
    if (!exam || loading) return;

    if (isProfessor) {
      loadClassScores();
    }
    loadMyScoring();
  }, [exam, loading]);

  const loadClassScores = useCallback(async () => {
    if (classLoading) return;
    setClassLoading(true);
    try {
      const data = await getUsersScoresByExam(examId, exam?.course_id || exam?.CourseID || null);
      const users = Array.isArray(data?.users) ? data.users : [];
      setClassScores(users);
    } catch (err) {
      console.warn('Could not load class scores:', err);
    } finally {
      setClassLoading(false);
    }
  }, [examId, exam, classLoading]);

  const loadMyScoring = useCallback(async () => {
    if (!userId || myLoading) return;
    setMyLoading(true);
    try {
      const data = await getExamScoringForUser(examId, userId);
      setMyScoring(data);
      setMyEmpty(false);
    } catch (err) {
      const raw = err?.response?.data?.error || err?.message || '';
      if (/no exam scores/i.test(raw)) {
        setMyEmpty(true);
      } else {
        console.warn('Could not load my scoring:', raw);
      }
    } finally {
      setMyLoading(false);
    }
  }, [examId, userId, myLoading]);

  // ──────────────────────────────────────────
  // 3. Load item-level detail (submissions)
  // ──────────────────────────────────────────
  const loadItemDetail = async (examScoreId) => {
    if (itemDetails[examScoreId]?.loaded) return;
    setItemDetails((prev) => ({
      ...prev,
      [examScoreId]: { loading: true },
    }));
    try {
      const detail = await getExamScoreItemDetails(examScoreId);
      setItemDetails((prev) => ({
        ...prev,
        [examScoreId]: { loaded: true, detail },
      }));
    } catch (err) {
      setItemDetails((prev) => ({
        ...prev,
        [examScoreId]: { loaded: true, error: err?.response?.data?.error || err?.message },
      }));
    }
  };

  const toggleScoreExpansion = (scoreId) => {
    const next = expandedScoreId === scoreId ? null : scoreId;
    setExpandedScoreId(next);
    if (next) loadItemDetail(next);
  };

  // ──────────────────────────────────────────
  // Professor: per-student drill-down
  // ──────────────────────────────────────────
  const [studentScoring, setStudentScoring] = useState({});

  const toggleStudentExpansion = async (uid) => {
    const next = expandedStudentId === uid ? null : uid;
    setExpandedStudentId(next);
    setExpandedScoreId(null);

    if (next && !studentScoring[uid]) {
      try {
        const data = await getExamScoringForUser(examId, uid);
        setStudentScoring((prev) => ({ ...prev, [uid]: data }));
      } catch (err) {
        const raw = err?.response?.data?.error || err?.message || '';
        setStudentScoring((prev) => ({ ...prev, [uid]: { error: raw } }));
      }
    }
  };

  // ──────────────────────────────────────────
  // Renders
  // ──────────────────────────────────────────
  if (loading) {
    return (
      <div className="exam-results-page">
        <PageLoader message="Cargando resultados del examen…" minHeight="300px" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="exam-results-page">
        <button className="er-back-btn" onClick={() => navigate(-1)}>
          <ArrowLeft size={16} /> Volver
        </button>
        <div className="er-error">{error}</div>
      </div>
    );
  }

  const examTitle = exam?.title || exam?.Title || 'Examen';
  const examDesc = exam?.description || exam?.Description || '';
  const courseId = exam?.course_id || exam?.CourseID || null;

  // ──────────────────────────────────────────
  // Render: score hero for a scoring block
  // ──────────────────────────────────────────
  const renderScoreHero = (scoringBlock) => {
    if (!scoringBlock) return null;
    const best = pick(scoringBlock, 'best_score', 'BestScore') ?? 0;
    const scores = Array.isArray(scoringBlock?.exam_scores)
      ? scoringBlock.exam_scores
      : Array.isArray(scoringBlock?.ExamScores)
        ? scoringBlock.ExamScores
        : [];

    const cls = scoreClass(best);
    return (
      <div className="er-score-hero">
        <div className={`er-score-circle ${cls}`}>{best}</div>
        <div className="er-score-info">
          <h3>Puntuación del Examen</h3>
          <p>Mejor calificación obtenida</p>
          <div className="er-score-meta">
            <span><Hash size={13} /> {scores.length} intento{scores.length !== 1 ? 's' : ''}</span>
            <span><Calendar size={13} /> {scores.length > 0 ? fmtDate(pick(scores[scores.length - 1], 'created_at', 'CreatedAt')) : '—'}</span>
          </div>
        </div>
      </div>
    );
  };

  // ──────────────────────────────────────────
  // Render: exam score attempts + item drill-down
  // ──────────────────────────────────────────
  const renderScoreAttempts = (scoringBlock) => {
    const scores = Array.isArray(scoringBlock?.exam_scores)
      ? scoringBlock.exam_scores
      : Array.isArray(scoringBlock?.ExamScores)
        ? scoringBlock.ExamScores
        : [];

    if (scores.length === 0) {
      return <p style={{ color: '#64748b', fontSize: '0.85rem' }}>Sin intentos registrados.</p>;
    }

    return (
      <div className="er-items-grid">
        {scores.map((row, idx) => {
          const sid = pick(row, 'id', 'ID');
          const sc = pick(row, 'score', 'Score') ?? 0;
          const created = pick(row, 'created_at', 'CreatedAt');
          const isOpen = expandedScoreId === sid;
          const detail = itemDetails[sid];

          return (
            <div key={sid || idx} className={`er-item-card ${sc === 100 ? 'solved' : sc > 0 ? 'attempted' : 'unattempted'}`}>
              <div className="er-item-header" onClick={() => toggleScoreExpansion(sid)}>
                <div className="er-item-title">
                  <Trophy size={16} />
                  <span>Intento #{idx + 1}</span>
                </div>
                <div className="er-item-badges">
                  <span className={`er-badge score-${scoreClass(sc)}`}>{sc} pts</span>
                  <span style={{ fontSize: '0.72rem', color: '#94a3b8' }}>{fmtDate(created)}</span>
                  <ChevronRight size={14} style={{ color: '#94a3b8', transform: isOpen ? 'rotate(90deg)' : 'none', transition: 'transform 0.2s' }} />
                </div>
              </div>

              {isOpen && (
                <div className="er-item-submissions">
                  {detail?.loading && (
                    <PageLoader message="Cargando detalle…" compact minHeight="0" size={14} />
                  )}
                  {detail?.error && <div className="er-error">{detail.error}</div>}
                  {detail?.detail && (
                    <>
                      {(detail.detail.exam_items || detail.detail.ExamItems || []).map((it, i) => {
                        const ei = it.exam_item || it.ExamItem || {};
                        const eis = it.exam_item_score || it.ExamItemScore || {};
                        const subs = Array.isArray(it.submissions) ? it.submissions : [];
                        const chTitle = pick(ei, 'title', 'Title') || `Reto #${(pick(ei, 'order', 'Order') ?? i) + 1}`;
                        const maxPts = pick(ei, 'points', 'Points') ?? 0;
                        const itemScore = pick(eis, 'score', 'Score') ?? 0;
                        const tries = pick(eis, 'tries', 'Tries') ?? 0;

                        return (
                          <div key={pick(eis, 'id', 'ID') || i} style={{ marginBottom: '0.75rem' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.35rem' }}>
                              <span style={{ display: 'flex', alignItems: 'center', gap: '6px', fontWeight: 700, fontSize: '0.85rem', color: '#1e293b' }}>
                                <Code size={14} style={{ color: '#c8102e' }} />
                                {chTitle}
                              </span>
                              <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                                <span className={`er-badge score-${scoreClass(itemScore)}`}>{itemScore}/{maxPts} pts</span>
                                <span style={{ fontSize: '0.72rem', color: '#94a3b8' }}>{tries} intento{tries !== 1 ? 's' : ''}</span>
                              </div>
                            </div>

                            {subs.length > 0 && subs.map((sub, si) => {
                              const subScore = pick(sub, 'score', 'Score') ?? 0;
                              const lang = pick(sub, 'language', 'Language') ?? '—';
                              const subDate = pick(sub, 'created_at', 'CreatedAt');

                              return (
                                <div
                                  key={pick(sub, 'id', 'ID') || si}
                                  className={`er-sub-row ${subScore === 100 ? 'passed' : ''}`}
                                  onClick={() => setSelectedSubmission({ ...sub, attemptNumber: si + 1, challengeTitle: chTitle })}
                                >
                                  <span className="er-sub-left">
                                    {subScore === 100
                                      ? <CheckCircle size={12} style={{ color: '#16a34a' }} />
                                      : <XCircle size={12} style={{ color: '#ef4444' }} />}
                                    Envío #{si + 1} — {lang}
                                  </span>
                                  <span className="er-sub-right">
                                    <span
                                      className="er-sub-score"
                                      style={{
                                        background: subScore === 100 ? '#dcfce7' : subScore >= 50 ? '#fef3c7' : '#fee2e2',
                                        color: subScore === 100 ? '#16a34a' : subScore >= 50 ? '#d97706' : '#ef4444',
                                      }}
                                    >
                                      {subScore}%
                                    </span>
                                    <span className="er-sub-date">{fmtDate(subDate)}</span>
                                  </span>
                                </div>
                              );
                            })}
                          </div>
                        );
                      })}
                    </>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    );
  };

  // ──────────────────────────────────────────
  // Render: Student view (My Results)
  // ──────────────────────────────────────────
  const renderStudentView = () => {
    if (myLoading) {
      return <PageLoader message="Cargando tus resultados…" compact minHeight="200px" />;
    }

    if (myEmpty || !myScoring) {
      return (
        <div className="er-empty">
          <div className="er-empty-icon"><Layers size={28} /></div>
          <h3>Sin resultados aún</h3>
          <p>Aún no has completado este examen o no se han registrado calificaciones.</p>
        </div>
      );
    }

    return (
      <>
        {renderScoreHero(myScoring)}
        <h3 style={{ fontSize: '1rem', fontWeight: 800, color: '#1e293b', marginBottom: '1rem' }}>
          Historial de Intentos
        </h3>
        {renderScoreAttempts(myScoring)}
      </>
    );
  };

  // ──────────────────────────────────────────
  // Render: Professor class view
  // ──────────────────────────────────────────
  const renderClassView = () => {
    if (classLoading) {
      return <PageLoader message="Cargando resultados de estudiantes…" compact minHeight="200px" />;
    }

    if (classScores.length === 0) {
      return (
        <div className="er-empty">
          <div className="er-empty-icon"><Users size={28} /></div>
          <h3>Sin resultados</h3>
          <p>Ningún estudiante ha realizado este examen todavía.</p>
        </div>
      );
    }

    return (
      <>
        <table className="er-students-table">
          <thead>
            <tr>
              <th>Estudiante</th>
              <th>Mejor Puntaje</th>
              <th>Intentos</th>
              <th>Último Intento</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {classScores.map((urow) => {
              const uid = pick(urow, 'user_id', 'UserID', 'userId') || 'unknown';
              const block = urow.exam_scores || urow.ExamScores || {};
              const best = pick(block, 'best_score', 'BestScore') ?? 0;
              const scores = Array.isArray(block?.exam_scores)
                ? block.exam_scores
                : Array.isArray(block?.ExamScores)
                  ? block.ExamScores
                  : [];
              const lastDate = scores.length > 0
                ? fmtDate(pick(scores[scores.length - 1], 'created_at', 'CreatedAt'))
                : '—';
              const isOpen = expandedStudentId === uid;
              const stuScoring = studentScoring[uid];

              return (
                <tr key={uid} style={{ verticalAlign: 'top' }}>
                  <td colSpan={5} style={{ padding: 0 }}>
                    <div
                      onClick={() => toggleStudentExpansion(uid)}
                      style={{
                        display: 'grid',
                        gridTemplateColumns: '2fr 1fr 1fr 1.5fr 40px',
                        alignItems: 'center',
                        padding: '0.75rem 1rem',
                        cursor: 'pointer',
                        transition: 'background 0.15s',
                      }}
                      onMouseEnter={(e) => (e.currentTarget.style.background = '#f8fafc')}
                      onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                    >
                      <div className="er-student-id">
                        <div className="er-student-avatar"><User size={14} /></div>
                        <span title={uid}>{uid.length > 24 ? uid.slice(0, 12) + '…' : uid}</span>
                      </div>
                      <div>
                        <span className={`er-score-pill ${scoreClass(best)}`}>{best} pts</span>
                      </div>
                      <div style={{ fontSize: '0.85rem', color: '#475569', fontWeight: 600 }}>
                        {scores.length}
                      </div>
                      <div style={{ fontSize: '0.8rem', color: '#94a3b8' }}>{lastDate}</div>
                      <div>
                        <button className={`er-expand-btn ${isOpen ? 'open' : ''}`}>
                          <ChevronRight size={16} />
                        </button>
                      </div>
                    </div>

                    {isOpen && (
                      <div className="er-student-detail">
                        {!stuScoring ? (
                          <PageLoader message="Cargando detalle…" compact minHeight="0" size={14} />
                        ) : stuScoring.error ? (
                          <div className="er-error">{stuScoring.error}</div>
                        ) : (
                          <>
                            {renderScoreHero(stuScoring)}
                            {renderScoreAttempts(stuScoring)}
                          </>
                        )}
                      </div>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </>
    );
  };

  // ──────────────────────────────────────────
  // Render: Code modal
  // ──────────────────────────────────────────
  const renderCodeModal = () => {
    if (!selectedSubmission) return null;
    const sc = pick(selectedSubmission, 'score', 'Score') ?? 0;
    const lang = pick(selectedSubmission, 'language', 'Language') ?? '—';
    const code = pick(selectedSubmission, 'code', 'Code') ?? 'Sin código';
    const results = Array.isArray(selectedSubmission.results) ? selectedSubmission.results : [];

    return (
      <div className="er-modal-backdrop" onClick={() => setSelectedSubmission(null)}>
        <div className="er-modal" onClick={(e) => e.stopPropagation()}>
          <div className="er-modal-header">
            <div>
              <h3><Code size={18} /> Detalles del Envío</h3>
              <div className="er-modal-meta">
                <span>{selectedSubmission.challengeTitle || 'Reto'}</span>
                <span>•</span>
                <span>Envío #{selectedSubmission.attemptNumber}</span>
                <span>•</span>
                <span>{lang}</span>
                <span>•</span>
                <span style={{ color: sc === 100 ? '#16a34a' : '#ef4444', fontWeight: 'bold' }}>{sc}%</span>
              </div>
            </div>
            <button className="er-modal-close" onClick={() => setSelectedSubmission(null)}>
              <XCircle size={24} />
            </button>
          </div>
          <div className="er-modal-body">
            <div className="er-modal-code">
              <div className="er-modal-section-title">Código Enviado</div>
              <pre>{code}</pre>
            </div>
            <div className="er-modal-results">
              <div className="er-modal-section-title">Casos de Prueba</div>
              <div className="er-modal-results-body">
                {results.length === 0 ? (
                  <div style={{ color: '#64748b', fontSize: '0.85rem', textAlign: 'center', marginTop: '2rem' }}>
                    No hay detalles de resultados.
                  </div>
                ) : (
                  results.map((r, i) => {
                    const st = (pick(r, 'status', 'Status') || 'unknown').toLowerCase();
                    const isAcc = st === 'accepted';
                    return (
                      <div key={i} className={`er-result-card ${isAcc ? 'accepted' : 'failed'}`}>
                        <div className="er-result-header">
                          <span>Caso #{i + 1}</span>
                          <span className={`er-result-status ${isAcc ? 'accepted' : 'failed'}`}>
                            {isAcc ? 'ACEPTADO' : 'FALLIDO'}
                          </span>
                        </div>
                        {!isAcc && (pick(r, 'error_message', 'ErrorMessage', 'errorMessage')) && (
                          <div className="er-result-error">
                            {pick(r, 'error_message', 'ErrorMessage', 'errorMessage')}
                          </div>
                        )}
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  };

  // ──────────────────────────────────────────
  // Main Render
  // ──────────────────────────────────────────
  return (
    <div className="exam-results-page">
      <button className="er-back-btn" onClick={() => navigate(-1)}>
        <ArrowLeft size={16} /> Volver
      </button>

      <div className="er-header">
        <div className="er-breadcrumb">
          <BookOpen size={13} />
          <Link to={courseId ? `/courses/${courseId}` : '/public-exams'}>
            {courseId ? 'Curso' : 'Exámenes Públicos'}
          </Link>
          <ChevronRight size={11} />
          <span>{examTitle}</span>
          <ChevronRight size={11} />
          <span>Resultados</span>
        </div>
        <h1 className="er-title">Resultados: {examTitle}</h1>
        {examDesc && <p className="er-subtitle">{examDesc}</p>}
      </div>

      {/* Tabs: professor gets class + mine; students see their own only */}
      {isProfessor && (
        <div className="er-tabs">
          <button
            className={`er-tab ${tab === 'class' ? 'active' : ''}`}
            onClick={() => setTab('class')}
          >
            <Users size={15} /> Resultados de Clase
          </button>
          <button
            className={`er-tab ${tab === 'mine' ? 'active' : ''}`}
            onClick={() => setTab('mine')}
          >
            <BarChart3 size={15} /> Mis Resultados
          </button>
        </div>
      )}

      {/* Content */}
      {tab === 'class' && isProfessor ? renderClassView() : renderStudentView()}

      {/* Code modal */}
      {renderCodeModal()}
    </div>
  );
};

export default ExamResults;
