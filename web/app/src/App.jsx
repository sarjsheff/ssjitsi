import { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import BotList from './components/BotList';
import LoadingSpinner from './components/LoadingSpinner';
import { useBots } from './hooks/useBots';
import { checkServerStatus } from './services/api';

/**
 * Основной компонент приложения для мониторинга ботов
 */
function App() {
  const {
    bots,
    loading,
    error,
    lastUpdate,
    refreshAll,
    loadScreenshot,
    loadAllScreenshots
  } = useBots();

  const [serverStatus, setServerStatus] = useState(null);
  const [refreshing, setRefreshing] = useState(false);

  // Проверка статуса сервера
  useEffect(() => {
    const checkStatus = async () => {
      const status = await checkServerStatus();
      setServerStatus(status);
    };

    checkStatus();
    const interval = setInterval(checkStatus, 5000); // Проверка каждые 5 секунд
    return () => clearInterval(interval);
  }, []);

  // Функция для обновления всех данных
  const handleRefreshAll = async () => {
    setRefreshing(true);
    try {
      await refreshAll();
    } finally {
      setRefreshing(false);
    }
  };

  // Функция для обновления скриншота конкретного бота
  const handleRefreshScreenshot = async (botId) => {
    await loadScreenshot(botId);
  };

  // Форматирование времени последнего обновления
  const formatLastUpdate = () => {
    if (!lastUpdate) return 'Никогда';
    return new Date(lastUpdate).toLocaleTimeString('ru-RU');
  };

  return (
    <div className="container-fluid py-3">
      {/* Шапка приложения */}
      <div className="row mb-4">
        <div className="col-12">
          <div className="d-flex justify-content-between align-items-center mb-3">
            <div>
              <h1 className="h3 mb-1">
                <i className="bi bi-robot me-2"></i>
                Мониторинг ботов Jitsi
              </h1>
              <p className="text-muted mb-0">
                Отслеживание состояния и активности ботов в реальном времени
              </p>
            </div>
            
            <div className="d-flex align-items-center gap-3">
              {/* Статус сервера */}
              <div className="d-flex align-items-center">
                <div 
                  className={`badge ${serverStatus === true ? 'bg-success' : serverStatus === false ? 'bg-danger' : 'bg-secondary'} me-2`}
                  style={{ width: '10px', height: '10px', borderRadius: '50%' }}
                ></div>
                <small className="text-muted">
                  {serverStatus === true ? 'Сервер доступен' : 
                   serverStatus === false ? 'Сервер недоступен' : 'Проверка...'}
                </small>
              </div>

              {/* Кнопка обновления */}
              <button
                className="btn btn-primary"
                onClick={handleRefreshAll}
                disabled={refreshing || loading}
              >
                {refreshing ? (
                  <>
                    <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                    Обновление...
                  </>
                ) : (
                  <>
                    <i className="bi bi-arrow-clockwise me-2"></i>
                    Обновить все
                  </>
                )}
              </button>
            </div>
          </div>

          {/* Статистика и информация */}
          <div className="row">
            <div className="col-12">
              <div className="d-flex flex-wrap gap-4 text-muted small">
                <div>
                  <i className="bi bi-collection me-1"></i>
                  <strong>Ботов:</strong> {bots.length}
                </div>
                <div>
                  <i className="bi bi-clock me-1"></i>
                  <strong>Последнее обновление:</strong> {formatLastUpdate()}
                </div>
                <div>
                  <i className="bi bi-arrow-repeat me-1"></i>
                  <strong>Автообновление:</strong> каждые 10 сек
                </div>
                <div>
                  <i className="bi bi-image me-1"></i>
                  <strong>Скриншоты:</strong> каждые 30 сек
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Основной контент - список ботов */}
      <div className="row">
        <div className="col-12">
          <BotList 
            bots={bots}
            loading={loading}
            error={error}
            onRefreshScreenshot={handleRefreshScreenshot}
          />
        </div>
      </div>

      {/* Футер с информацией */}
      <div className="row mt-5">
        <div className="col-12">
          <div className="text-center text-muted small">
            <hr />
            <p className="mb-1">
              Система мониторинга ботов Jitsi &copy; {new Date().getFullYear()}
            </p>
            <p className="mb-0">
              Автоматическое обновление данных и скриншотов в реальном времени
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
