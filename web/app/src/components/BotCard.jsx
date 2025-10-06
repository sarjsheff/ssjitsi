import { useState } from 'react';
import LoadingSpinner from './LoadingSpinner';

/**
 * Компонент карточки бота
 */
const BotCard = ({
  bot,
  onRefreshScreenshot,
  onStopBot,
  onRestartBot
}) => {
  const [actionInProgress, setActionInProgress] = useState(false);

  const handleStop = async () => {
    setActionInProgress(true);
    try {
      await onStopBot(bot.id);
    } catch (error) {
      console.error('Ошибка остановки бота:', error);
    } finally {
      setActionInProgress(false);
    }
  };

  const handleRestart = async () => {
    setActionInProgress(true);
    try {
      await onRestartBot(bot.id);
    } catch (error) {
      console.error('Ошибка перезапуска бота:', error);
    } finally {
      setActionInProgress(false);
    }
  };

  const getStatusBadgeClass = (status) => {
    switch (status) {
      case 'running':
        return 'bg-success';
      case 'stopped':
        return 'bg-secondary';
      case 'starting':
        return 'bg-warning';
      case 'stopping':
        return 'bg-warning';
      default:
        return 'bg-secondary';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'running':
        return 'Работает';
      case 'stopped':
        return 'Остановлен';
      case 'starting':
        return 'Запускается';
      case 'stopping':
        return 'Останавливается';
      default:
        return 'Неизвестно';
    }
  };
  const formatTime = (date) => {
    if (!date) return 'Никогда';
    return new Date(date).toLocaleTimeString('ru-RU');
  };

  const formatDate = (date) => {
    if (!date) return '';
    return new Date(date).toLocaleDateString('ru-RU');
  };

  return (
    <div className="card h-100 shadow-sm">
      {/* Заголовок карточки */}
      <div className="card-header bg-light p-2">
        <div className="d-flex justify-content-between align-items-center mb-1">
          <h6 className="card-title mb-0 text-truncate" title={`${bot.room} | ${bot.botName}`}>
            {bot.room} | {bot.botName}
          </h6>
        </div>
        <div className="d-flex justify-content-between align-items-center">
          <span className={`badge ${getStatusBadgeClass(bot.status)}`}>
            {getStatusText(bot.status)}
          </span>
          <div className="btn-group btn-group-sm" role="group">
            <button
              className="btn btn-outline-success"
              onClick={handleRestart}
              disabled={actionInProgress || bot.status === 'starting' || bot.status === 'stopping'}
              title="Перезапустить"
            >
              <i className="bi bi-arrow-repeat"></i>
            </button>
            <button
              className="btn btn-outline-danger"
              onClick={handleStop}
              disabled={actionInProgress || bot.status === 'stopped' || bot.status === 'stopping'}
              title="Остановить"
            >
              <i className="bi bi-stop-circle"></i>
            </button>
            <button
              className="btn btn-outline-primary"
              onClick={() => onRefreshScreenshot(bot.id)}
              disabled={bot.loadingScreenshot}
              title="Обновить скриншот"
            >
              <i className="bi bi-arrow-clockwise"></i>
            </button>
          </div>
        </div>
      </div>

      <div className="card-body d-flex flex-column p-2">
        {/* Скриншот */}
        <div className="mb-2 text-center">
          {bot.loadingScreenshot ? (
            <div className="d-flex justify-content-center align-items-center" style={{ height: '120px' }}>
              <LoadingSpinner size="sm" text="Загрузка..." />
            </div>
          ) : bot.screenshotError ? (
            <div className="alert alert-warning text-center py-2 mb-0" style={{ height: '120px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <div>
                <i className="bi bi-exclamation-triangle fs-5"></i>
                <p className="mb-0 small">{bot.screenshotError}</p>
              </div>
            </div>
          ) : bot.screenshot ? (
            <img
              src={bot.screenshot}
              alt={`Скриншот бота ${bot.id}`}
              className="img-fluid rounded border"
              style={{
                maxHeight: '120px',
                width: 'auto',
                objectFit: 'contain'
              }}
              onError={(e) => {
                e.target.style.display = 'none';
                e.target.nextSibling.style.display = 'block';
              }}
            />
          ) : (
            <div className="text-muted d-flex justify-content-center align-items-center" style={{ height: '120px' }}>
              <div>
                <i className="bi bi-image fs-3"></i>
                <p className="mb-0 small">Нет скриншота</p>
              </div>
            </div>
          )}
        </div>

        {/* Информация о боте */}
        <div className="mt-auto">
          <div className="row small text-muted">
            <div className="col-12 mb-1">
              <strong>Сервер:</strong>
              <div className="text-truncate" title={bot.server}>
                {bot.server}
              </div>
            </div>

            <div className="col-12">
              <div className="d-flex justify-content-between align-items-center">
                <div>
                  <strong>Авторизация:</strong>
                  <span className={`badge ms-1 ${bot.authMethod === 'JWT' ? 'bg-success' : 'bg-primary'}`}>
                    {bot.authMethod}
                  </span>
                </div>
                <div className="d-flex align-items-center">
                  <div
                    className={`badge ${bot.screenshot ? 'bg-success' : 'bg-secondary'} me-1`}
                    style={{ width: '8px', height: '8px', borderRadius: '50%' }}
                  ></div>
                  <small>
                    {bot.screenshot ? 'Активен' : 'Неактивен'}
                  </small>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Футер карточки */}
      <div className="card-footer bg-transparent border-top-0 py-1 px-2">
        <small className="text-muted">
          {bot.lastScreenshotUpdate
            ? `${formatTime(bot.lastScreenshotUpdate)}`
            : 'Ожидание...'
          }
        </small>
      </div>
    </div>
  );
};

export default BotCard;